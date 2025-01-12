provider "aws" {
  region              = "us-east-1"
  allowed_account_ids = ["919759177803"]
}

# terraform = https://github.com/hashicorp/terraform/releases
# aws provider = https://registry.terraform.io/providers/hashicorp/aws/latest
terraform {
  required_version = ">= 1.10.4, < 2.0.0"
  required_providers {
    aws = {
      version = "< 6.0"
    }
  }
  backend "s3" {
    bucket = "codabool-tf"
    key    = "shared.tfstate"
    region = "us-east-1"
  }
}

data "aws_caller_identity" "current" {}

module "lambda_emailer" {
  source = "./modules/emailer"
}

module "lambda_scraper" {
  source = "./modules/scraper"
}

module "lambda_discord" {
  source = "./modules/discord"
}

module "lambda_discord_slash" {
  source = "./modules/discord_slash"
}

module "lambda_discord_reminder" {
  source = "./modules/discord_reminder"
}

module "key" {
  source   = "./modules/key"
  key_name = "win"
}

module "actions" {
  source  = "./modules/actions"
  account = data.aws_caller_identity.current.account_id
}

resource "aws_ssm_parameter" "all_env" {
  name        = "/env"
  description = "A comma seperated list of all aws envs"
  type        = "SecureString"
  value       = data.external.read_all_env.result.env
}

data "external" "read_all_env" {
  program = ["bash", "readenv.sh"]
}

# module "test" {
#   source = "./modules/delete"
# }
