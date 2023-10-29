package groups

import (
	"context"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// groupResource is the resource implementation.
type groupResource struct {
	client *clients.Client
}

type groupModel struct {
	ID      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Owner   types.String  `tfsdk:"owner"`
	Algo    types.String  `tfsdk:"algo"`
	Size    types.Int64   `tfsdk:"size"`
	Servers []serverModel `tfsdk:"servers"`
}

type Server struct {
	Host    string
	User    string
	Port    int64
	Comment string
}

type serverModel struct {
	Host    types.String `tfsdk:"host"`
	User    types.String `tfsdk:"user"`
	Port    types.Int64  `tfsdk:"port"`
	Comment types.String `tfsdk:"comment"`
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
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the group",
				Required:    true,
			},
			"algo": schema.StringAttribute{
				Description: "Algorithm used to generate the key",
				Required:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Size of the key",
				Required:    true,
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
						"comment": schema.StringAttribute{
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

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.String()
	owner := plan.Owner.String()
	algo := plan.Algo.String()
	size := plan.Size.ValueInt64()
	servers := plan.Servers

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the group
	_, err := r.client.CreateGroup(ctx, name, owner, algo, size)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating group",
			err.Error(),
		)
		return
	}

	// Add servers to the group
	for _, server := range servers {
		_, err := r.client.AddServerToGroup(ctx, name, server.Host.String(), server.User.String(), server.Port.ValueInt64(), server.Comment.String())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while adding server to group",
				err.Error(),
			)
			return
		}
	}

	planIDString, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating id",
			"Could not create id for testing: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(planIDString)

	// Set state to the resource
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
			"Error while destroying group",
			err.Error(),
		)
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
