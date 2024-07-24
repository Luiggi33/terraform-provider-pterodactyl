package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Luiggi33/pterodactyl-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &userResource{}
	_ resource.ResourceWithConfigure = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *pterodactyl.Client
}

// userResourceModel maps the resource schema data.
type userResourceModel struct {
	ID        int    `tfsdk:"id"`
	Username  string `tfsdk:"username"`
	Email     string `tfsdk:"email"`
	FirstName string `tfsdk:"first_name"`
	LastName  string `tfsdk:"last_name"`
	CreatedAt string `tfsdk:"created_at"`
	UpdatedAt string `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"first_name": schema.StringAttribute{
				Required: true,
			},
			"last_name": schema.StringAttribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new resource.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create partial user
	partialUser := pterodactyl.PartialUser{
		Username:  plan.Username,
		Email:     plan.Email,
		FirstName: plan.FirstName,
		LastName:  plan.LastName,
	}

	// Create new order
	user, err := r.client.CreateUser(partialUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = user.ID
	plan.CreatedAt = user.CreatedAt.Format(time.RFC3339)
	plan.CreatedAt = time.Now().Format(time.RFC3339)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed user value from Pterodactyl
	user, err := r.client.GetUser(state.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pterodactyl User",
			"Could not read Pterodactyl user ID "+strconv.Itoa(state.ID)+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Email = user.Email
	state.FirstName = user.FirstName
	state.LastName = user.LastName
	state.UpdatedAt = user.UpdatedAt.Format(time.RFC3339)
	state.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
