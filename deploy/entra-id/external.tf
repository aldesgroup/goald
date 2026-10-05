// ------------------------------------------------------------------------------------------------
// External users: Microsoft Entra External ID (CIAM tenant)
//
// The app registration mechanics are identical to the workforce tenant (same Microsoft Graph
// objects), so this mirrors internal.tf. What Terraform can't do (yet) is configure the CIAM
// "user flow" (sign-up / sign-in experience, local account policy, password strength, etc.) - that
// part is still done once, by hand, in the Azure portal - see docs/authentication.md.
// ------------------------------------------------------------------------------------------------

data "azuread_client_config" "external" {
  provider = azuread.external
}

# --- The API itself -------------------------------------------------------------------------

resource "azuread_application" "external_api" {
  provider     = azuread.external
  display_name = "${var.app_display_name_prefix} API - External"
  owners       = [data.azuread_client_config.external.object_id]
}

resource "azuread_application_identifier_uri" "external_api" {
  provider       = azuread.external
  application_id = azuread_application.external_api.id
  identifier_uri = "api://${azuread_application.external_api.client_id}"
}

resource "azuread_service_principal" "external_api" {
  provider  = azuread.external
  client_id = azuread_application.external_api.client_id
  owners    = [data.azuread_client_config.external.object_id]
}

# the delegated scope that clients request to call this API on behalf of a signed-in external user
resource "random_uuid" "external_scope" {}

resource "azuread_application_permission_scope" "external_access_as_user" {
  provider       = azuread.external
  application_id = azuread_application.external_api.id
  scope_id       = random_uuid.external_scope.result
  value          = "access_as_user"
  type           = "User"

  admin_consent_description  = "Allows the app to access the Goald API on behalf of the signed-in external user."
  admin_consent_display_name = "Access the Goald API as an external user"
  user_consent_description   = "Allows the app to access the Goald API on your behalf."
  user_consent_display_name  = "Access the Goald API"
}

# --- A public-client app registration, for testing the sign-in flow (auth code + PKCE) -------

resource "azuread_application" "external_cli" {
  provider     = azuread.external
  display_name = "${var.app_display_name_prefix} CLI - External (testing)"
  owners       = [data.azuread_client_config.external.object_id]

  # no password (ROPC) grant: users sign in on Microsoft's own page
  fallback_public_client_enabled = false

  public_client {
    redirect_uris = var.cli_public_client_redirect_uris
  }

  required_resource_access {
    resource_app_id = azuread_application.external_api.client_id

    resource_access {
      id   = azuread_application_permission_scope.external_access_as_user.scope_id
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "external_cli" {
  provider  = azuread.external
  client_id = azuread_application.external_cli.client_id
  owners    = [data.azuread_client_config.external.object_id]
}

# admin-consenting the CLI app for the API's scope, so sign-ins never hit an interactive consent screen
resource "azuread_service_principal_delegated_permission_grant" "external_cli_consent" {
  provider                             = azuread.external
  service_principal_object_id          = azuread_service_principal.external_cli.object_id
  resource_service_principal_object_id = azuread_service_principal.external_api.object_id
  claim_values                         = ["access_as_user"]
}
