package groups

import (
	"context"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*clients.Client)
}

// NewGroupsDataSource is a helper function to simplify the provider implementation.
func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

// groupsDataSource is the data source implementation.
type groupsDataSource struct {
	client *clients.Client
}

type groupFlags struct {
	Flags []string `tfsdk:"flags"`
}

// groupsDataSourceModel maps the data source schema data.
type groupsDataSourceModel struct {
	ID     types.String          `tfsdk:"id"`
	Groups map[string]groupFlags `tfsdk:"groups"`
}

// Metadata returns the data source type name.
func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

// Schema defines the schema for the data source.
func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of list of groups. Required by terraform-plugin-testing",
				Computed:    true,
			},
			"groups": schema.MapNestedAttribute{
				Description: "List of groups",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flags": schema.ListAttribute{
							Description: "List of flags of group.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupsDataSourceModel
	state.Groups = make(map[string]groupFlags)

	responseBastionGroupList, err := d.client.GetListGroup(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TheBastion Groups",
			err.Error(),
		)
		return
	}

	for name, group := range responseBastionGroupList.Value {
		groupState := groupFlags{
			Flags: group.Flags,
		}
		state.Groups[name] = groupState
	}

	stateIDString, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating id",
			"Could not create id for testing: "+err.Error(),
		)
		return
	}
	state.ID = types.StringValue(stateIDString)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
