package users

import (
	"context"
	"fmt"
	"sort"
	"terraform-provider-thebastion/thebastion/clients"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *clients.Client
}

// usersDataSourceModel maps the data source schema data.
type usersDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Users []userModel  `tfsdk:"users"`
}

// userModel maps users schema data.
type userModel struct {
	ID           types.String `tfsdk:"id"`
	Uid          types.Int64  `tfsdk:"uid"`
	Name         types.String `tfsdk:"name"`
	Is_active    types.Int64  `tfsdk:"is_active"`
	Ingress_keys types.List   `tfsdk:"ingress_keys"`
}

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

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of list of users. Required by terraform-plugin-testing",
				Computed:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "List of users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of user. Required by terraform-plugin-testing",
							Computed:    true,
						},
						"uid": schema.Int64Attribute{
							Description: "UID of user.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of user. Used as an unique identifier by TheBastion.",
							Computed:    true,
						},
						"is_active": schema.Int64Attribute{
							Description: "Is the user active.",
							Computed:    true,
						},
						"ingress_keys": schema.ListAttribute{
							Description: "List of ingress keys of users.",
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
	state := usersDataSourceModel{}

	responseBastionAccountList, err := d.client.GetListAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TheBastion Users",
			err.Error(),
		)
		return
	}

	for name, bastionUser := range responseBastionAccountList.Value {
		responseBastionIngressKeys, err := d.client.GetListIngressKeys(ctx, name)
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

		stateIDString, err := uuid.GenerateUUID()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error generating id",
				"Could not create id for testing: "+err.Error(),
			)
			return
		}

		userState := userModel{
			ID:           types.StringValue(stateIDString),
			Uid:          types.Int64Value(bastionUser.UID),
			Name:         types.StringValue(bastionUser.Name),
			Is_active:    types.Int64Value(bastionUser.IsActive),
			Ingress_keys: list,
		}

		state.Users = append(state.Users, userState)
	}

	// Sort users by their uid
	sort.Slice(state.Users, func(i, j int) bool {
		return state.Users[i].Uid.ValueInt64() < state.Users[j].Uid.ValueInt64()
	})

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
