// ------------------------------------------------------------------------------------------------
// These map directly onto the "Auth" section of Goald's config file - see docs/authentication.md
// ------------------------------------------------------------------------------------------------

output "colleague_tenant_id" {
  value = var.colleague_tenant_id
}

output "colleague_api_client_id" {
  description = "Use this as Auth.colleague.clientId and (once logged in) as the 'aud' to expect - Auth.colleague.audience."
  value       = azuread_application.colleague_api.client_id
}

output "colleague_cli_client_id" {
  description = "The public/test client's ID - use this to try the login flow with curl, see docs/authentication.md."
  value       = azuread_application.colleague_cli.client_id
}

output "colleague_scope" {
  value = "api://${azuread_application.colleague_api.client_id}/access_as_user"
}

output "customer_tenant_id" {
  value = var.customer_tenant_id
}

output "customer_api_client_id" {
  description = "Use this as Auth.customer.clientId and (once logged in) as the 'aud' to expect - Auth.customer.audience."
  value       = azuread_application.customer_api.client_id
}

output "customer_cli_client_id" {
  description = "The public/test client's ID - use this to try the login flow with curl, see docs/authentication.md."
  value       = azuread_application.customer_cli.client_id
}

output "customer_scope" {
  value = "api://${azuread_application.customer_api.client_id}/access_as_user"
}
