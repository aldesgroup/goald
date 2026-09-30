// ------------------------------------------------------------------------------------------------
// Terraform providers required to provision the Entra ID (colleagues) and Entra External ID
// (customers) app registrations used by Goald's authentication feature.
//
// Usage: see docs/authentication.md for the full picture. In short:
//   terraform init
//   terraform apply -var-file=my.tfvars
//
// You need to be authenticated against BOTH tenants for a single `apply` to work (e.g. via
// `az login` with a user/service principal that has access to each tenant, or by setting
// ARM_* / AZURE_* environment variables per provider alias). If that's not practical, split this
// into 2 separate `terraform apply` runs (one per tenant), commenting out the other tenant's
// resources, or move colleague.tf / customer.tf into 2 separate Terraform root modules.
// ------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.5"

  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

# the Entra ID (workforce) tenant, where your colleagues live
provider "azuread" {
  alias     = "colleague"
  tenant_id = var.colleague_tenant_id
}

# the Entra External ID (CIAM) tenant, where your customers live
provider "azuread" {
  alias     = "customer"
  tenant_id = var.customer_tenant_id
}
