// ------------------------------------------------------------------------------------------------
// Customers: Microsoft Entra External ID (CIAM tenant)
//
// The app registration mechanics are identical to the workforce tenant (same Microsoft Graph
// objects), so this mirrors colleague.tf. What Terraform can't do (yet) is configure the CIAM
// "user flow" (sign-up / sign-in experience, local account policy, password strength, etc.) - that
// part is still done once, by hand, in the Azure portal - see docs/authentication.md.
// ------------------------------------------------------------------------------------------------

data "azuread_client_config" "customer" {
  provider = azuread.customer
}

# --- The API itself -------------------------------------------------------------------------

resource "azuread_application" "customer_api" {
  provider     = azuread.customer
  display_name = "${var.app_display_name_prefix} API - Customers"
  owners       = [data.azuread_client_config.customer.object_id]
}

resource "azuread_application_identifier_uri" "customer_api" {
  provider       = azuread.customer
  application_id = azuread_application.customer_api.id
  identifier_uri = "api://${azuread_application.customer_api.client_id}"
}

resource "azuread_service_principal" "customer_api" {
  provider  = azuread.customer
  client_id = azuread_application.customer_api.client_id
  owners    = [data.azuread_client_config.customer.object_id]
}

# the delegated scope that clients request to call this API on behalf of a signed-in customer
resource "random_uuid" "customer_scope" {}

resource "azuread_application_permission_scope" "customer_access_as_user" {
  provider       = azuread.customer
  application_id = azuread_application.customer_api.id
  scope_id       = random_uuid.customer_scope.result
  value          = "access_as_user"
  type           = "User"

  admin_consent_description  = "Allows the app to access the Goald API on behalf of the signed-in customer."
  admin_consent_display_name = "Access the Goald API as a customer"
  user_consent_description   = "Allows the app to access the Goald API on your behalf."
  user_consent_display_name  = "Access the Goald API"
}

# --- A public-client app registration, for testing the login flow (ROPC via curl) ------------

resource "azuread_application" "customer_cli" {
  provider     = azuread.customer
  display_name = "${var.app_display_name_prefix} CLI - Customers (testing)"
  owners       = [data.azuread_client_config.customer.object_id]

  # this is what the Azure portal calls "Allow public client flows" - required for the Resource
  # Owner Password Credentials (ROPC) grant used by curl / goald.Login()
  fallback_public_client_enabled = true

  public_client {
    redirect_uris = var.cli_public_client_redirect_uris
  }

  required_resource_access {
    resource_app_id = azuread_application.customer_api.client_id

    resource_access {
      id   = azuread_application_permission_scope.customer_access_as_user.scope_id
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "customer_cli" {
  provider  = azuread.customer
  client_id = azuread_application.customer_cli.client_id
  owners    = [data.azuread_client_config.customer.object_id]
}

# admin-consenting the CLI app for the API's scope, so ROPC logins never hit an interactive consent screen
resource "azuread_service_principal_delegated_permission_grant" "customer_cli_consent" {
  provider                             = azuread.customer
  service_principal_object_id          = azuread_service_principal.customer_cli.object_id
  resource_service_principal_object_id = azuread_service_principal.customer_api.object_id
  claim_values                         = ["access_as_user"]
}
