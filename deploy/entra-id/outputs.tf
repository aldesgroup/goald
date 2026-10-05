// ------------------------------------------------------------------------------------------------
// These map directly onto the "Auth" section of Goald's config file - see docs/authentication.md
// ------------------------------------------------------------------------------------------------

output "internal_tenant_id" {
  value = var.internal_tenant_id
}

output "internal_api_client_id" {
  description = "Use this as Auth.internal.clientId and (once logged in) as the 'aud' to expect - Auth.internal.audience."
  value       = azuread_application.internal_api.client_id
}

output "internal_cli_client_id" {
  description = "The public/test client's ID - use this to try the login flow with curl, see docs/authentication.md."
  value       = azuread_application.internal_cli.client_id
}

output "internal_scope" {
  value = "api://${azuread_application.internal_api.client_id}/access_as_user"
}

output "external_tenant_id" {
  value = var.external_tenant_id
}

output "external_api_client_id" {
  description = "Use this as Auth.external.clientId and (once logged in) as the 'aud' to expect - Auth.external.audience."
  value       = azuread_application.external_api.client_id
}

output "external_cli_client_id" {
  description = "The public/test client's ID - use this to try the login flow with curl, see docs/authentication.md."
  value       = azuread_application.external_cli.client_id
}

output "external_scope" {
  value = "api://${azuread_application.external_api.client_id}/access_as_user"
}
