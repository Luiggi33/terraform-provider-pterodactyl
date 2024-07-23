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
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *pterodactyl.Client
}

// usersDataSourceModel maps the data source schema data.
type usersDataSourceModel struct {
	Users []User `tfsdk:"users"`
}

// Users schema data.
type User struct {
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

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"external_id": schema.StringAttribute{
							Computed: true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"first_name": schema.StringAttribute{
							Computed: true,
						},
						"last_name": schema.StringAttribute{
							Computed: true,
						},
						"language": schema.StringAttribute{
							Computed: true,
						},
						"root_admin": schema.BoolAttribute{
							Computed: true,
						},
						"is_2fa": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
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

	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Pterodactyl Users",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, user := range users {
		userState := User{
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

		state.Users = append(state.Users, userState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
