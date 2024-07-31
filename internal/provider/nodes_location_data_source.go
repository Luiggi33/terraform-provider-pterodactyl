package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/Luiggi33/pterodactyl-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &nodesLocationDataSource{}
	_ datasource.DataSourceWithConfigure = &nodesLocationDataSource{}
)

// NewNodesLocationDataSource is a helper function to simplify the provider implementation.
func NewNodesLocationDataSource() datasource.DataSource {
	return &nodesLocationDataSource{}
}

// nodesLocationDataSource is the data source implementation.
type nodesLocationDataSource struct {
	client *pterodactyl.Client
}

// nodesLocationDataSourceModel maps the data source schema data.
type nodesLocationDataSourceModel struct {
	LocationID types.Int32 `tfsdk:"location_id"`
	Nodes      []Node      `tfsdk:"nodes"`
}

// Metadata returns the data source type name.
func (d *nodesLocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nodes_location"
}

// Schema defines the schema for the data source.
func (d *nodesLocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Pterodactyl nodes data source allows Terraform to read nodes from the Pterodactyl API.",
		Attributes: map[string]schema.Attribute{
			"location_id": schema.Int32Attribute{
				Description: "The ID of the location.",
				Required:    true,
			},
			"nodes": schema.ListNestedAttribute{
				Description: "The list of nodes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Description: "The ID of the node.",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Description: "The UUID of the node.",
							Computed:    true,
						},
						"public": schema.BoolAttribute{
							Description: "The public status of the node.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the node.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the node.",
							Computed:    true,
						},
						"location_id": schema.Int32Attribute{
							Description: "The location ID of the node.",
							Computed:    true,
						},
						"fqdn": schema.StringAttribute{
							Description: "The FQDN of the node.",
							Computed:    true,
						},
						"scheme": schema.StringAttribute{
							Description: "The scheme of the node.",
							Computed:    true,
						},
						"behind_proxy": schema.BoolAttribute{
							Description: "The behind proxy status of the node.",
							Computed:    true,
						},
						"maintenance_mode": schema.BoolAttribute{
							Description: "The maintenance mode status of the node.",
							Computed:    true,
						},
						"memory": schema.Int32Attribute{
							Description: "The memory of the node.",
							Computed:    true,
						},
						"memory_overallocate": schema.Int32Attribute{
							Description: "The memory overallocate of the node.",
							Computed:    true,
						},
						"disk": schema.Int32Attribute{
							Description: "The disk of the node.",
							Computed:    true,
						},
						"disk_overallocate": schema.Int32Attribute{
							Description: "The disk overallocate of the node.",
							Computed:    true,
						},
						"upload_size": schema.Int32Attribute{
							Description: "The upload size of the node.",
							Computed:    true,
						},
						"daemon_listen": schema.Int32Attribute{
							Description: "The daemon listen of the node.",
							Computed:    true,
						},
						"daemon_sftp": schema.Int32Attribute{
							Description: "The daemon SFTP of the node.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The creation date of the node.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "The last update date of the node.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *nodesLocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state nodesLocationDataSourceModel

	nodes, err := d.client.GetNodes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Pterodactyl Nodes",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, node := range nodes {
		if node.LocationID != state.LocationID.ValueInt32() {
			continue
		}
		state.Nodes = append(state.Nodes, Node{
			ID:                 types.Int32Value(node.ID),
			UUID:               types.StringValue(node.UUID),
			Public:             types.BoolValue(node.Public),
			Name:               types.StringValue(node.Name),
			Description:        types.StringValue(node.Description),
			LocationID:         types.Int32Value(node.LocationID),
			FQDN:               types.StringValue(node.FQDN),
			Scheme:             types.StringValue(node.Scheme),
			BehindProxy:        types.BoolValue(node.BehindProxy),
			MaintenanceMode:    types.BoolValue(node.MaintenanceMode),
			Memory:             types.Int32Value(node.Memory),
			MemoryOverallocate: types.Int32Value(node.MemoryOverallocate),
			Disk:               types.Int32Value(node.Disk),
			DiskOverallocate:   types.Int32Value(node.DiskOverallocate),
			UploadSize:         types.Int32Value(node.UploadSize),
			DaemonListen:       types.Int32Value(node.DaemonListen),
			DaemonSFTP:         types.Int32Value(node.DaemonSFTP),
			CreatedAt:          types.StringValue(node.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:          types.StringValue(node.UpdatedAt.Format(time.RFC3339)),
		})
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *nodesLocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pterodactyl.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pterodactyl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
