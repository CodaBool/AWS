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

# module "emailer" {
#   source = "./modules/emailer"
# }
module "texter" {
  source = "./modules/texter"
}

module "scraper" {
  source = "./modules/scraper"
}

# module "cheapo" {
#   source = "./modules/ec2"
# }

module "key" {
  source = "./modules/key"
  key_name = "win"
}

module "actions" {
  source = "./modules/actions"
}

module "s3" {
  source = "./modules/s3"
}

module "cloudwatch" {
  source = "./modules/cloudwatch"
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

output "actions_role" {
  value       = module.actions.role_arn
  description = "github actions assume role"
}