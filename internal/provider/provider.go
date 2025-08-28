package provider

import (
	"context"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure KindProvider satisfies various provider interfaces.
var _ provider.Provider = &KindProvider{}

// KindProvider defines the provider implementation.
type KindProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KindProviderModel describes the provider data model.
type KindProviderModel struct {
	Runtime types.String `tfsdk:"runtime"`
}

// ProviderData contains provider configuration that gets passed to resources
type ProviderData struct {
	Runtime string
}

func (p *KindProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kind"
	resp.Version = p.version
}

func (p *KindProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Kind provider allows Terraform to manage Kubernetes Kind clusters.",
		Attributes: map[string]schema.Attribute{
			"runtime": schema.StringAttribute{
				Description: "Container runtime to use for Kind clusters. Valid values are 'docker' and 'podman'. If not specified, the provider will auto-detect the available runtime.",
				Optional:    true,
			},
		},
	}
}

func (p *KindProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KindProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Determine runtime
	runtime := "docker" // default
	if !data.Runtime.IsNull() && !data.Runtime.IsUnknown() {
		specifiedRuntime := data.Runtime.ValueString()
		if specifiedRuntime != "docker" && specifiedRuntime != "podman" {
			resp.Diagnostics.AddError(
				"Invalid Runtime",
				"Runtime must be either 'docker' or 'podman'",
			)
			return
		}
		runtime = specifiedRuntime
	} else {
		// Auto-detect runtime
		if detected := detectRuntime(); detected != "" {
			runtime = detected
		}
	}

	// Create provider data
	providerData := &ProviderData{
		Runtime: runtime,
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

// detectRuntime attempts to auto-detect the available container runtime
func detectRuntime() string {
	// Check for podman first since it's less common and users who have it likely want to use it
	if _, err := exec.LookPath("podman"); err == nil {
		return "podman"
	}

	// Check for docker
	if _, err := exec.LookPath("docker"); err == nil {
		return "docker"
	}

	// Default to docker if nothing is found
	return "docker"
}

func (p *KindProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

func (p *KindProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Data sources can be added here if needed
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KindProvider{
			version: version,
		}
	}
}
