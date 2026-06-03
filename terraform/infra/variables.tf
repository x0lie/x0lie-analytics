variable "cloudflare_api_token" {
  description = "Cloudflare API token with DNS edit permissions"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for x0lie.com"
  type        = string
}

variable "db_password" {
  description = "Password for RDS"
  type        = string
  sensitive   = true
}

variable "image_tag" {
  description = "x0lie-analytics container image tag"
  type        = string
}
