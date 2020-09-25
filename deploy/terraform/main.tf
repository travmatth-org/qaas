variable "github_oauth_token" {}
data "aws_caller_identity" "current" {}

locals {
  user_name = basename(data.aws_caller_identity.current.arn)
}

terraform {
  required_version = ">=0.12.28"

  # backend "s3" {
  #   bucket         = "qaas-terraform-state-bucket-${var.aws_account_id}"
  #   key            = "terraform.tfstate"
  #   region         = "us-west-1"
  #   dynamodb_table = "qaas-dynamodb-terraform-locking"
  #   encrypt        = true
  # }
}

provider "aws" {
  version = "~> 2.8"
  region  = "us-west-1"

  # assume_role {
  #   role_arn     = "arn:aws:iam::${var.aws_account_id}:role/TerraformAssumedIamRole"
  #   session_name = "qaas-session"
  # }
}

provider "github" {
  version      = "~> 2.9"
  token        = var.github_oauth_token
  organization = "travmatth-org"
}

resource "random_password" "github_secret" {
  length  = 22
  special = false
}

data "github_repository" "qaas" {
  name = "qaas"
}

resource "aws_ssm_parameter" "github_oauth_token" {
  name        = "/qaas/ci/github-oauth-token"
  description = "GitHub OAuth Token"
  type        = "SecureString"
  value       = var.github_oauth_token

  tags = {
    qaas = "true"
  }
}

resource "aws_ssm_parameter" "webhook_secret" {
  name        = "/qaas/ci/webhook-secret"
  description = "CICD WebHook Secret"
  type        = "SecureString"
  value       = random_password.github_secret.result

  tags = {
    qaas = "true"
  }
}

module "tf_backend" {
  source         = "./modules/terraform_backend"
  aws_account_id = data.aws_caller_identity.current.account_id
}

# module "cicd" {
#   source                       = "./modules/cicd"
#   github_oauth_token           = var.github_oauth_token
#   webhook_secret               = random_password.github_secret.result
#   codebuild_logging_bucket     = module.tf_backend.codebuild_logging_bucket
#   tf_state_bucket              = module.tf_backend.tf_state_bucket
#   codepipeline_artifact_bucket = module.tf_backend.codepipeline_artifact_bucket
#   dynamodb_lock_state_table    = module.tf_backend.client_lock_state_table
#   github_repo                  = data.github_repository.qaas
#   account_id                   = data.aws_caller_identity.current.account_id
#   user_name                    = local.user_name
# }

# module "network" {
#   source = "./modules/network"
# }

# module "alb" {
#   source                        = "./modules/alb"
#   vpc                           = module.network.vpc
#   public_subnets                = module.network.public_subnets
#   internet_gateway              = module.network.internet_gateway
#   codepipeline_artifact_bucket  = module.tf_backend.codepipeline_artifact_bucket
#   account_id                    = data.aws_caller_identity.current.account_id
#   user_name                     = local.user_name
# }

# output "alb_ip" {
#   value = module.alb.lb_dns_name
# }

module "dynamodb" {
  source = "./modules/dynamodb"
}
