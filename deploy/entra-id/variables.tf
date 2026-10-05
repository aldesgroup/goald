variable "internal_tenant_id" {
  description = "Tenant ID (GUID) of the Microsoft Entra ID (workforce) tenant used to authenticate internal users."
  type        = string
}

variable "external_tenant_id" {
  description = "Tenant ID (GUID) of the Microsoft Entra External ID (CIAM) tenant used to authenticate external users."
  type        = string
}

variable "app_display_name_prefix" {
  description = "Prefix used for the display name of every app registration created by this module, e.g. your app's name."
  type        = string
  default     = "Goald"
}

variable "cli_public_client_redirect_uris" {
  description = "Redirect URIs allowed for the public/CLI test clients (only relevant if you also exercise interactive/browser flows against them)."
  type        = list(string)
  default     = ["http://localhost"]
}
