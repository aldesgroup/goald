// ------------------------------------------------------------------------------------------------
// Internal users: Microsoft Entra ID (workforce tenant)
//
// Creates:
//  - the "API" app registration representing this Goald backend, exposing 1 delegated scope
//  - a public-client "CLI" app registration, pre-consented for that scope, so you can test the
//    whole login flow right away with curl (Resource Owner Password Credentials grant)
// ------------------------------------------------------------------------------------------------

data "azuread_client_config" "internal" {
  provider = azuread.internal
}

# --- The API itself -------------------------------------------------------------------------

resource "azuread_application" "internal_api" {
  provider     = azuread.internal
  display_name = "${var.app_display_name_prefix} API - Internal"
  owners       = [data.azuread_client_config.internal.object_id]
}

resource "azuread_application_identifier_uri" "internal_api" {
  provider       = azuread.internal
  application_id = azuread_application.internal_api.id
  identifier_uri = "api://${azuread_application.internal_api.client_id}"
}

resource "azuread_service_principal" "internal_api" {
  provider  = azuread.internal
  client_id = azuread_application.internal_api.client_id
  owners    = [data.azuread_client_config.internal.object_id]
}

# the delegated scope that clients request to call this API on behalf of a signed-in internal user
resource "random_uuid" "internal_scope" {}

resource "azuread_application_permission_scope" "internal_access_as_user" {
  provider       = azuread.internal
  application_id = azuread_application.internal_api.id
  scope_id       = random_uuid.internal_scope.result
  value          = "access_as_user"
  type           = "User"

  admin_consent_description  = "Allows the app to access the Goald API on behalf of the signed-in internal user."
  admin_consent_display_name = "Access the Goald API as an internal user"
  user_consent_description   = "Allows the app to access the Goald API on your behalf."
  user_consent_display_name  = "Access the Goald API"
}

# --- A public-client app registration, for testing the login flow (ROPC via curl) ------------

resource "azuread_application" "internal_cli" {
  provider     = azuread.internal
  display_name = "${var.app_display_name_prefix} CLI - Internal (testing)"
  owners       = [data.azuread_client_config.internal.object_id]

  # this is what the Azure portal calls "Allow public client flows" - required for the Resource
  # Owner Password Credentials (ROPC) grant used by curl / goald.Login()
  fallback_public_client_enabled = true

  public_client {
    redirect_uris = var.cli_public_client_redirect_uris
  }

  required_resource_access {
    resource_app_id = azuread_application.internal_api.client_id

    resource_access {
      id   = azuread_application_permission_scope.internal_access_as_user.scope_id
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "internal_cli" {
  provider  = azuread.internal
  client_id = azuread_application.internal_cli.client_id
  owners    = [data.azuread_client_config.internal.object_id]
}

# admin-consenting the CLI app for the API's scope, so ROPC logins never hit an interactive consent screen
resource "azuread_service_principal_delegated_permission_grant" "internal_cli_consent" {
  provider                             = azuread.internal
  service_principal_object_id          = azuread_service_principal.internal_cli.object_id
  resource_service_principal_object_id = azuread_service_principal.internal_api.object_id
  claim_values                         = ["access_as_user"]
}
