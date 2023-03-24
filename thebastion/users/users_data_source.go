package users

import (
	"context"
	"fmt"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*clients.Client)
}

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *clients.Client
}

// usersDataSourceModel maps the data source schema data.
type usersDataSourceModel struct {
	Users []usersModel `tfsdk:"users"`
}

// usersModel maps users schema data.
type usersModel struct {
	Uid          types.Int64  `tfsdk:"uid"`
	Name         types.String `tfsdk:"name"`
	Is_active    types.Int64  `tfsdk:"is_active"`
	Ingress_keys types.List   `tfsdk:"ingress_keys"`
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"is_active": schema.Int64Attribute{
							Computed: true,
						},
						"ingress_keys": schema.ListAttribute{
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
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	responseBastionAccountList, err := d.client.GetListAccount()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TheBastion Users",
			err.Error(),
		)
		return
	}

	for name, bastionUser := range responseBastionAccountList.Value {
		responseBastionIngressKeys, err := d.client.GetListIngressKeys(name)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Read TheBastion User Ingress Keys: %s", name),
				err.Error(),
			)
			return
		}
		listIngressKeys := []string{}
		for _, key := range responseBastionIngressKeys.Value.Keys {
			listIngressKeys = append(listIngressKeys, key.Line)
		}

		list, _ := types.ListValueFrom(ctx, types.StringType, listIngressKeys)

		userState := usersModel{
			Uid:          types.Int64Value(bastionUser.UID),
			Name:         types.StringValue(bastionUser.Name),
			Is_active:    types.Int64Value(bastionUser.IsActive),
			Ingress_keys: list,
		}

		state.Users = append(state.Users, userState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
