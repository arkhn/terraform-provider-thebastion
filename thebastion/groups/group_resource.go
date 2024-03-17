package groups

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// groupResource is the resource implementation.
type groupResource struct {
	client *clients.Client
}

type groupModel struct {
	ID      types.String          `tfsdk:"id"`
	Name    types.String          `tfsdk:"name"`
	Owners  []types.String        `tfsdk:"owners"`
	Algo    types.String          `tfsdk:"algo"`
	Size    types.Int64           `tfsdk:"size"`
	Servers []clients.ServerModel `tfsdk:"servers"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

// Configure adds the provider configured client to the data source.
func (d *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*clients.Client)
}

// NewUserResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

// Metadata returns the resource type name.
func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (group groupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of resource. Required by terraform-plugin-testing",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the group",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owners": schema.ListAttribute{
				Description: "List of owners",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"algo": schema.StringAttribute{
				Description: "Algorithm used to generate the key",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Size of the key",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"servers": schema.ListNestedAttribute{
				Description: "List of servers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Description: "Host of the server",
							Required:    true,
						},
						"user": schema.StringAttribute{
							Description: "User of the server",
							Required:    true,
						},
						"port": schema.Int64Attribute{
							Description: "Port of the server",
							Required:    true,
						},
						"user_comment": schema.StringAttribute{
							Description: "Comment of the server",
							Required:    true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

// Create a new resource
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := groupModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, algo, size := plan.Name.String(), plan.Algo.String(), plan.Size.ValueInt64()
	owners := []string{}
	for _, owner := range plan.Owners {
		owners = append(owners, owner.String())
	}

	servers := plan.Servers

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the group
	_, err := r.client.CreateGroup(ctx, name, owners, algo, size)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while creating group: %s", err.Error()),
		)
		return
	}

	// Add servers to the group
	for _, server := range servers {
		_, err := r.client.AddServerToGroup(ctx, name, server.Host.String(), server.User.String(), server.Port.ValueInt64(), server.UserComment.String())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Error while adding server to group: %s", err.Error()),
			)
			return
		}
	}

	uuid, _ := uuid.GenerateUUID()
	plan.ID = types.StringValue(uuid)

	// Set state to the resource
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	state := groupModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := r.client.GetListGroup(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while reading groups: %s", err.Error()),
		)
		return
	}

	// Check if group in groups keys
	if _, ok := groups.Value[state.Name.ValueString()]; !ok {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Group %s not found", state.Name.ValueString()),
		)
		return
	}

	groupInfo, err := r.client.GetGroupInfo(ctx, state.Name.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while reading group info: %s", err.Error()),
		)
		return
	}

	// Get servers of the group
	servers, err := r.client.GetListServer(ctx, state.Name.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while reading servers of group: %s", err.Error()),
		)
		return
	}

	// Overwrite state with the latest data
	state.Name = types.StringValue(groupInfo.Value.Group)
	state.Owners = []types.String{}
	for _, owner := range groupInfo.Value.Owners {
		state.Owners = append(state.Owners, types.StringValue(owner))
	}
	for _, key := range groupInfo.Value.Keys {
		state.Algo = types.StringValue(key.Typecode)
		state.Size = types.Int64Value(key.Size)
		break
	}

	state.Servers = []clients.ServerModel{}
	for _, server := range servers.Value {
		port, err := strconv.ParseInt(server.Port, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Error while parsing port: %s", err.Error()),
			)
			return
		}
		state.Servers = append(state.Servers, clients.ServerModel{
			Host:        types.StringValue(server.IP),
			User:        types.StringValue(server.User),
			Port:        types.Int64Value(port),
			UserComment: types.StringValue(server.UserComment),
		})
	}

	uuid, _ := uuid.GenerateUUID()
	state.ID = types.StringValue(uuid)

	// Set state to the resource
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan, state := groupModel{}, groupModel{}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update servers of the group
	servers, err := r.client.UpdateServerFromGroup(ctx, state.Name.String(), plan.Servers, state.Servers)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while updating servers of group: %s", err.Error()),
		)
		return
	}

	planOwners := []string{}
	for _, owner := range plan.Owners {
		planOwners = append(planOwners, owner.String())
	}

	stateOwners := []string{}
	for _, owner := range state.Owners {
		stateOwners = append(stateOwners, owner.String())
	}
	// Update owners of the group
	_, err = r.client.UpdateOwnerFromGroup(ctx, state.Name.String(), planOwners, stateOwners)

	// Overwrite state with the latest data
	state.Servers = []clients.ServerModel{}
	for _, server := range servers.Value {
		port, err := strconv.ParseInt(server.Port, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Error while parsing port: %s", err.Error()),
			)
			return
		}
		state.Servers = append(state.Servers, clients.ServerModel{
			Host:        types.StringValue(server.IP),
			User:        types.StringValue(server.User),
			Port:        types.Int64Value(port),
			UserComment: types.StringValue(server.UserComment),
		})
	}

	uuid, _ := uuid.GenerateUUID()
	plan.ID = types.StringValue(uuid)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the current state
	state := groupModel{}
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the group
	_, err := r.client.DeleteGroup(ctx, state.Name.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error while deleting group: %s", err.Error()),
		)
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
