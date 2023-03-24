package thebastion

import (
	"context"
	"os"
	"terraform-provider-thebastion/thebastion/clients"
	"terraform-provider-thebastion/thebastion/users"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &thebastionProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &thebastionProvider{}
}

// thebastionProvider is the provider implementation.
type thebastionProvider struct{}

// hashicupsProviderModel maps provider schema data to a Go type.
type thebastionProviderModel struct {
	Host           types.String `tfsdk:"host"`
	Username       types.String `tfsdk:"username"`
	PathPrivateKey types.String `tfsdk:"path_private_key"`
	PathKnownHost  types.String `tfsdk:"path_known_host"`
}

// Metadata returns the provider type name.
func (p *thebastionProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "thebastion"
}

// Schema defines the provider-level schema for configuration data.
func (p *thebastionProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"path_private_key": schema.StringAttribute{
				Optional: true,
			},
			"path_known_host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *thebastionProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config thebastionProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown TheBastion API Host",
			"The provider cannot create the TheBastion API client as there is an unknown configuration value for the TheBastion API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the THEBASTION_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown TheBastion API Username",
			"The provider cannot create the TheBastion API client as there is an unknown configuration value for the TheBastion API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the THEBASTION_USERNAME environment variable.",
		)
	}

	if config.PathPrivateKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path_known_host"),
			"Unknown TheBastion API PathPrivateKey",
			"The provider cannot create the TheBastion API client as there is an unknown configuration value for the TheBastion API PathPrivateKey. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the THEBASTION_PATH_PRIVATE_KEY environment variable.",
		)
	}

	if config.PathKnownHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path_known_host"),
			"Unknown TheBastion API PathKnownHost",
			"The provider cannot create the TheBastion API client as there is an unknown configuration value for the TheBastion API PathKnownHost. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the THEBASTION_PATH_KNOWN_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("THEBASTION_HOST")
	username := os.Getenv("THEBASTION_USERNAME")
	path_known_host := os.Getenv("THEBASTION_PATH_KNOWN_HOST")
	path_private_key := os.Getenv("THEBASTION_PATH_PRIVATE_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.PathKnownHost.IsNull() {
		path_known_host = config.PathKnownHost.ValueString()
	}

	if !config.PathPrivateKey.IsNull() {
		path_private_key = config.PathPrivateKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing TheBastion API Host",
			"The provider cannot create the TheBastion API client as there is a missing or empty value for the TheBastion API host. "+
				"Set the host value in the configuration or use the THEBASTION_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing TheBastion API Username",
			"The provider cannot create the TheBastion API client as there is a missing or empty value for the TheBastion API username. "+
				"Set the username value in the configuration or use the THEBASTION_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if path_known_host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path_known_host"),
			"Missing TheBastion API PathKnownHost",
			"The provider cannot create the TheBastion API client as there is a missing or empty value for the TheBastion API PathKnownHost. "+
				"Set the host value in the configuration or use the THEBASTION_PATH_KNOWN_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if path_private_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path_private_key"),
			"Missing TheBastion API PathPrivateKey",
			"The provider cannot create the TheBastion API client as there is a missing or empty value for the TheBastion API PathPrivateKey. "+
				"Set the username value in the configuration or use the THEBASTION_PATH_PRIVATE_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// log about the configuration of bastion
	ctx = tflog.SetField(ctx, "thebastion_host", host)
	ctx = tflog.SetField(ctx, "thebastion_username", username)
	ctx = tflog.SetField(ctx, "thebastion_path_known_host", path_known_host)
	ctx = tflog.SetField(ctx, "thebastion_path_private_key", path_private_key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "thebastion_host")
	tflog.Info(ctx, "Setting TheBastion parameters")

	// Create a new TheBastion client using the configuration values
	client, err := clients.NewClient(host, username, path_private_key, path_known_host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create TheBastion API Client",
			"An unexpected error occurred when creating the TheBastion API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"TheBastion Client Error: "+err.Error(),
		)
		return
	}

	// Make the TheBastion client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *thebastionProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		users.NewUsersDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *thebastionProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
