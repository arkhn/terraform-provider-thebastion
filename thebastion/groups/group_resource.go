package groups

import (
	"context"
	"terraform-provider-thebastion/thebastion/clients"

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
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	Owner            types.String  `tfsdk:"owner"`
	Algo             types.String  `tfsdk:"algo"`
	Size             types.Int64   `tfsdk:"size"`
	Key_is_encrypted types.Bool    `tfsdk:"key_is_encrypted"`
	Servers          []serverModel `tfsdk:"servers"`
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
			"name": schema.StringAttribute{
				Description: "Name of the group",
				Required:    true,
			},
			"owner": schema.StringAttribute{
				Required: true,
			},
			"algo": schema.StringAttribute{
				Required: true,
			},
			"size": schema.NumberAttribute{
				Required: true,
			},
			"servers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Required: true,
						},
						"user": schema.StringAttribute{
							Required: true,
						},
						"port": schema.NumberAttribute{
							Required: true,
						},
						"comment": schema.StringAttribute{
							Required: true,
						},
					},
				},
				Required: false,
			},
		},
	}
}

// Create a new resource
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Throw exception not implemented

}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
