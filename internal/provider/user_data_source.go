package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/Luiggi33/pterodactyl-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// userDataSourceModel maps the data source schema data.
type userDataSourceModel struct {
	ID         int    `tfsdk:"id"`
	ExternalID string `tfsdk:"external_id"`
	UUID       string `tfsdk:"uuid"`
	Username   string `tfsdk:"username"`
	Email      string `tfsdk:"email"`
	FirstName  string `tfsdk:"first_name"`
	LastName   string `tfsdk:"last_name"`
	Language   string `tfsdk:"language"`
	RootAdmin  bool   `tfsdk:"root_admin"`
	Is2FA      bool   `tfsdk:"is_2fa"`
	CreatedAt  string `tfsdk:"created_at"`
	UpdatedAt  string `tfsdk:"updated_at"`
}

// NewUserDataSource is a helper function to simplify the provider implementation.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// userDataSource is the data source implementation.
type userDataSource struct {
	client *pterodactyl.Client
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the user.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the user.",
				Optional:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email of the user.",
				Optional:    true,
			},
			// Add other attributes as needed
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel

	// Get the attributes from the request
	var target struct {
		ID       int    `tfsdk:"id"`
		Username string `tfsdk:"username"`
		Email    string `tfsdk:"email"`
	}
	diags := req.Config.Get(ctx, &target)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the user from the API based on the provided attribute
	var user pterodactyl.User
	var err error
	if target.ID != 0 {
		user, err = d.client.GetUser(target.ID)
	} else if target.Username != "" {
		user, err = d.client.GetUserUsername(target.Username)
	} else if target.Email != "" {
		user, err = d.client.GetUserEmail(target.Email)
	} else {
		resp.Diagnostics.AddError(
			"Missing Attribute",
			"One of 'id', 'username', or 'email' must be specified.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Pterodactyl User",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = userDataSourceModel{
		ID:         user.ID,
		ExternalID: user.ExternalID,
		UUID:       user.UUID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Language:   user.Language,
		RootAdmin:  user.RootAdmin,
		Is2FA:      user.Is2FA,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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