package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Luiggi33/pterodactyl-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &nodeResource{}
	_ resource.ResourceWithConfigure = &nodeResource{}
	// _ resource.ResourceWithImportState = &nodeResource{}
)

// NewNodeResource is a helper function to simplify the provider implementation.
func NewNodeResource() resource.Resource {
	return &nodeResource{}
}

// nodeResource is the resource implementation.
type nodeResource struct {
	client *pterodactyl.Client
}

// nodeResourceModel maps the resource schema data.
type nodeResourceModel struct {
	ID                 types.Int32  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Public             types.Bool   `tfsdk:"public"`
	BehindProxy        types.Bool   `tfsdk:"behind_proxy"`
	MaintenanceMode    types.Bool   `tfsdk:"maintenance_mode"`
	LocationID         types.Int32  `tfsdk:"location_id"`
	FQDN               types.String `tfsdk:"fqdn"`
	Scheme             types.String `tfsdk:"scheme"`
	Memory             types.Int32  `tfsdk:"memory"`
	MemoryOverallocate types.Int32  `tfsdk:"memory_overallocate"`
	Disk               types.Int32  `tfsdk:"disk"`
	DiskOverallocate   types.Int32  `tfsdk:"disk_overallocate"`
	UploadSize         types.Int32  `tfsdk:"upload_size"`
	DaemonSFTP         types.Int32  `tfsdk:"daemon_sftp"`
	DaemonListen       types.Int32  `tfsdk:"daemon_listen"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *nodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

// Schema defines the schema for the resource.
func (r *nodeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Pterodactyl node resource allows Terraform to manage nodes in the Pterodactyl Panel API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Description: "The ID of the node.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the node.",
				Optional:    true,
			},
			"public": schema.BoolAttribute{
				Description: "The public status of the node.",
				Required:    true,
			},
			"behind_proxy": schema.BoolAttribute{
				Description: "The behind proxy status of the node.",
				Required:    true,
			},
			"maintenance_mode": schema.BoolAttribute{
				Description: "The maintenance mode status of the node.",
				Required:    true,
			},
			"location_id": schema.Int32Attribute{
				Description: "The location ID of the node.",
				Required:    true,
			},
			"fqdn": schema.StringAttribute{
				Description: "The FQDN of the node.",
				Required:    true,
			},
			"scheme": schema.StringAttribute{
				Description: "The scheme of the node.",
				Required:    true,
			},
			"memory": schema.Int32Attribute{
				Description: "The memory of the node.",
				Required:    true,
			},
			"memory_overallocate": schema.Int32Attribute{
				Description: "The memory overallocate of the node.",
				Required:    true,
			},
			"disk": schema.Int32Attribute{
				Description: "The disk of the node.",
				Required:    true,
			},
			"disk_overallocate": schema.Int32Attribute{
				Description: "The disk overallocate of the node.",
				Required:    true,
			},
			"upload_size": schema.Int32Attribute{
				Description: "The upload size of the node.",
				Required:    true,
			},
			"daemon_sftp": schema.Int32Attribute{
				Description: "The daemon SFTP of the node.",
				Required:    true,
			},
			"daemon_listen": schema.Int32Attribute{
				Description: "The daemon listen of the node.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The creation date of the node.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "The last update date of the node.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create a new resource.
func (r *nodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan nodeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create partial node
	partialNode := pterodactyl.PartialNode{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Public:             plan.Public.ValueBool(),
		BehindProxy:        plan.BehindProxy.ValueBool(),
		MaintenanceMode:    plan.MaintenanceMode.ValueBool(),
		LocationID:         plan.LocationID.ValueInt32(),
		FQDN:               plan.FQDN.ValueString(),
		Scheme:             plan.Scheme.ValueString(),
		Memory:             plan.Memory.ValueInt32(),
		MemoryOverallocate: plan.MemoryOverallocate.ValueInt32(),
		Disk:               plan.Disk.ValueInt32(),
		DiskOverallocate:   plan.DiskOverallocate.ValueInt32(),
		UploadSize:         plan.UploadSize.ValueInt32(),
		DaemonSFTP:         plan.DaemonSFTP.ValueInt32(),
		DaemonListen:       plan.DaemonListen.ValueInt32(),
	}

	// Create new node
	node, err := r.client.CreateNode(partialNode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating node",
			"Could not create node, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int32Value(node.ID)
	plan.CreatedAt = types.StringValue(node.CreatedAt.Format(time.RFC3339))

	if !plan.Description.IsNull() {
		updatedNode, err := r.client.UpdateNode(int32(plan.ID.ValueInt32()), partialNode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Pterodactyl Node",
				"Could not update node, unexpected error: "+err.Error(),
			)
			return
		}

		// Update resource state with updated values
		plan.Description = types.StringValue(updatedNode.Description)
		plan.UpdatedAt = types.StringValue(updatedNode.UpdatedAt.Format(time.RFC3339))
	}

	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *nodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state nodeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed node value from Pterodactyl
	node, err := r.client.GetNode(int32(state.ID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pterodactyl Node",
			"Could not read Pterodactyl node ID "+strconv.FormatInt(int64(state.ID.ValueInt32()), 10)+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(node.Name)
	state.Description = types.StringValue(node.Description)
	state.Public = types.BoolValue(node.Public)
	state.BehindProxy = types.BoolValue(node.BehindProxy)
	state.MaintenanceMode = types.BoolValue(node.MaintenanceMode)
	state.LocationID = types.Int32Value(node.LocationID)
	state.FQDN = types.StringValue(node.FQDN)
	state.Scheme = types.StringValue(node.Scheme)
	state.Memory = types.Int32Value(node.Memory)
	state.MemoryOverallocate = types.Int32Value(node.MemoryOverallocate)
	state.Disk = types.Int32Value(node.Disk)
	state.DiskOverallocate = types.Int32Value(node.DiskOverallocate)
	state.UploadSize = types.Int32Value(node.UploadSize)
	state.DaemonSFTP = types.Int32Value(node.DaemonSFTP)
	state.DaemonListen = types.Int32Value(node.DaemonListen)
	state.CreatedAt = types.StringValue(node.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(node.UpdatedAt.Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *nodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan nodeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create partial node
	partialNode := pterodactyl.PartialNode{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Public:             plan.Public.ValueBool(),
		BehindProxy:        plan.BehindProxy.ValueBool(),
		MaintenanceMode:    plan.MaintenanceMode.ValueBool(),
		LocationID:         plan.LocationID.ValueInt32(),
		FQDN:               plan.FQDN.ValueString(),
		Scheme:             plan.Scheme.ValueString(),
		Memory:             plan.Memory.ValueInt32(),
		MemoryOverallocate: plan.MemoryOverallocate.ValueInt32(),
		Disk:               plan.Disk.ValueInt32(),
		DiskOverallocate:   plan.DiskOverallocate.ValueInt32(),
		UploadSize:         plan.UploadSize.ValueInt32(),
		DaemonSFTP:         plan.DaemonSFTP.ValueInt32(),
		DaemonListen:       plan.DaemonListen.ValueInt32(),
	}

	// Update existing node
	node, err := r.client.UpdateNode(int32(plan.ID.ValueInt32()), partialNode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Pterodactyl Node",
			"Could not update node, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource plan with updated values
	plan.Name = types.StringValue(node.Name)
	plan.Description = types.StringValue(node.Description)
	plan.Public = types.BoolValue(node.Public)
	plan.BehindProxy = types.BoolValue(node.BehindProxy)
	plan.MaintenanceMode = types.BoolValue(node.MaintenanceMode)
	plan.LocationID = types.Int32Value(node.LocationID)
	plan.FQDN = types.StringValue(node.FQDN)
	plan.Scheme = types.StringValue(node.Scheme)
	plan.Memory = types.Int32Value(node.Memory)
	plan.MemoryOverallocate = types.Int32Value(node.MemoryOverallocate)
	plan.Disk = types.Int32Value(node.Disk)
	plan.DiskOverallocate = types.Int32Value(node.DiskOverallocate)
	plan.UploadSize = types.Int32Value(node.UploadSize)
	plan.DaemonSFTP = types.Int32Value(node.DaemonSFTP)
	plan.DaemonListen = types.Int32Value(node.DaemonListen)
	plan.CreatedAt = types.StringValue(node.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(node.UpdatedAt.Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *nodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state nodeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing node
	err := r.client.DeleteNode(state.ID.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Pterodactyl Node",
			"Could not delete node, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *nodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *nodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)

	node, err := r.client.GetNode(int32(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Pterodactyl User",
			"Could not import node: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state := nodeResourceModel{
		ID:                 types.Int32Value(node.ID),
		Name:               types.StringValue(node.Name),
		Description:        types.StringValue(node.Description),
		Public:             types.BoolValue(node.Public),
		BehindProxy:        types.BoolValue(node.BehindProxy),
		MaintenanceMode:    types.BoolValue(node.MaintenanceMode),
		LocationID:         types.Int32Value(node.LocationID),
		FQDN:               types.StringValue(node.FQDN),
		Scheme:             types.StringValue(node.Scheme),
		Memory:             types.Int32Value(node.Memory),
		MemoryOverallocate: types.Int32Value(node.MemoryOverallocate),
		Disk:               types.Int32Value(node.Disk),
		DiskOverallocate:   types.Int32Value(node.DiskOverallocate),
		UploadSize:         types.Int32Value(node.UploadSize),
		DaemonSFTP:         types.Int32Value(node.DaemonSFTP),
		DaemonListen:       types.Int32Value(node.DaemonListen),
		CreatedAt:          types.StringValue(node.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:          types.StringValue(node.UpdatedAt.Format(time.RFC3339)),
	}

	// Set state to fully populated data
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
