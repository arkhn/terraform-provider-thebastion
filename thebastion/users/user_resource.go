package users

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// userResource is the resource implementation.
type userResource struct {
	client *clients.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// Configure adds the provider configured client to the data source.
func (d *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*clients.Client)
}

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of resource. Required by terraform-plugin-testing",
				Computed:    true,
			},
			"uid": schema.Int64Attribute{
				Description: "UID of user.",
				Required:    true,
				// Need to replace this resource if this attribute is planned for update
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of user. Used as an unique identifier by TheBastion.",
				Required:    true,
				// Need to replace this resource if this attribute is planned for update
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					// Check regex to validate name
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([a-zA-Z0-9._-]+)$`),
						"must contain only UNIX valid characters",
					),
				},
			},
			"ingress_keys": schema.ListAttribute{
				Description: "List of ingress keys of users.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"is_active": schema.Int64Attribute{
				Description: "Is the user active.",
				Computed:    true,
			},
		},
	}
}

// Create a new resource
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan := userModel{}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get parameters from plan
	name := plan.Name.String()
	uid := plan.Uid.ValueInt64()
	ingress_keys := []string{}
	diags = plan.Ingress_keys.ElementsAs(ctx, &ingress_keys, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new user
	_, err := r.client.CreateAccount(ctx, name, uid, ingress_keys)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	// New user are always active when created
	plan.Is_active = types.Int64Value(1)
	planIDString, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating id",
			"Could not create id for testing: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(planIDString)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	state := userModel{}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed user value from TheBastion
	name := state.Name.ValueString()
	account, err := r.client.GetAccount(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			"Could not get user, unexpected error: "+err.Error(),
		)
		return
	}

	if account != nil {
		account_ingress_keys, err := r.client.GetListIngressKeys(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Error while reading ingress keys of user: %s", err.Error()),
			)
			return
		}
		ingress_keys_string := []string{}
		for _, key := range account_ingress_keys.Value.Keys {
			ingress_keys_string = append(ingress_keys_string, key.Line)
		}
		ingress_keys, diags := types.ListValueFrom(ctx, types.StringType, ingress_keys_string)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Overwrite items with refreshed state
		state.Name = types.StringValue(account.Name)
		state.Uid = types.Int64Value(account.UID)
		state.Is_active = types.Int64Value(account.IsActive)
		state.Ingress_keys = ingress_keys

		uuid, _ := uuid.GenerateUUID()
		state.ID = types.StringValue(uuid)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	plan, state := userModel{}, userModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	name := plan.Name.ValueString()
	planIngressKeys, stateIngressKeys := []string{}, []string{}
	for _, key := range plan.Ingress_keys.Elements() {
		planIngressKeys = append(planIngressKeys, key.String())
	}

	resp.Diagnostics.Append(plan.Ingress_keys.ElementsAs(ctx, &planIngressKeys, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.Ingress_keys.ElementsAs(ctx, &stateIngressKeys, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing user
	err := r.client.UpdateListIngressKeys(ctx, name, stateIngressKeys, planIngressKeys)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while updating ingress keys of user: %s", err.Error()),
		)
		return
	}

	// Fetch updated items from GetAccount as UpdateUser items are not
	// populated.
	account, err := r.client.GetAccount(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while reading user: %s", err.Error()),
		)
		return
	}

	ingress_keys := basetypes.ListValue{}

	// If an account with this name is found
	if account != nil {
		account_ingress_keys, err := r.client.GetListIngressKeys(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Error while reading ingress keys of user: %s", err.Error()),
			)
			return
		}
		ingress_keys_string := []string{}
		for _, key := range account_ingress_keys.Value.Keys {
			ingress_keys_string = append(ingress_keys_string, key.Line)
		}
		_, diags := types.ListValueFrom(ctx, types.StringType, ingress_keys_string)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update resource state with updated items and timestamp
	plan.Name = types.StringValue(account.Name)
	plan.Is_active = types.Int64Value(account.IsActive)
	plan.Uid = types.Int64Value(account.UID)
	plan.Ingress_keys = ingress_keys

	uuid, _ := uuid.GenerateUUID()
	plan.ID = types.StringValue(uuid)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state := userModel{}
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing user
	_, err := r.client.DeleteAccount(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting TheBastion User",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
