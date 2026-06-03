terraform {
  backend "s3" {
    bucket       = "x0lie-analytics-terraform-state"
    key          = "infra/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
    encrypt      = true
  }
}
