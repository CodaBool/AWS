provider "aws" {
  region              = "us-east-1"
  allowed_account_ids = ["919759177803"]
}

terraform {
  required_version = ">= 1.3.6, < 2.0.0"
  backend "s3" {
    bucket = "codabool-tf"
    key    = "shared.tfstate"
    region = "us-east-1"
  }
}

module "emailer" {
  source = "./modules/emailer"
}

module "scraper" {
  source = "./modules/scraper"
}