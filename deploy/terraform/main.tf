variable "github_oauth_token" {}
data "aws_caller_identity" "current" {}

terraform {
  required_version = ">=0.12.28"

  # backend "s3" {
  #   bucket         = "faas-terraform-state-bucket-${var.aws_account_id}"
  #   key            = "terraform.tfstate"
  #   region         = "us-west-1"
  #   dynamodb_table = "faas-dynamodb-terraform-locking"
  #   encrypt        = true
  # }
}

provider "aws" {
  version = "~>2.8"
  region = "us-west-1"

  # assume_role {
  #   role_arn     = "arn:aws:iam::${var.aws_account_id}:role/TerraformAssumedIamRole"
  #   session_name = "faas-session"
  # }
}

provider "github" {
  version       = "~>2.9"
  token         = var.github_oauth_token
  organization  = "travmatth-org"
}

resource "random_password" "github_secret" {
  length  = 22
  special = false
}

data "github_repository" "faas" {
  name = "faas"
}

resource "aws_ssm_parameter" "github_oauth_token" {
  name        = "/faas/ci/github_oauth_token"
  description = "GitHub OAuth Token"
  type        = "SecureString"
  value       = var.github_oauth_token

  tags        = {
    FaaS      = "true"
  }
}

resource "aws_ssm_parameter" "webhook_secret" {
  name        = "/faas/ci/webhook_secret"
  description = "CICD WebHook Secret"
  type        = "SecureString"
  value       = random_password.github_secret.result

  tags        = {
    FaaS      = "true"
  }
}

module "tf_backend" {
  source          = "./modules/terraform_backend"
  aws_account_id  = data.aws_caller_identity.current.account_id
}

module "cicd" {
  source                       = "./modules/cicd"
  github_oauth_token           = var.github_oauth_token
  webhook_secret               = random_password.github_secret.result
  codebuild_logging_bucket     = module.tf_backend.codebuild_logging_bucket
  tf_state_bucket              = module.tf_backend.tf_state_bucket
  codepipeline_artifact_bucket = module.tf_backend.codepipeline_artifact_bucket
  dynamodb_lock_state_table    = module.tf_backend.dynamodb_lock_state_table
  github_repo                  = data.github_repository.faas
}

module "network" {
  source = "./modules/network"
}

module "ec2" {
  source           = "./modules/ec2"
  public_vpc       = module.network.public_vpc
  public_subnet    = module.network.public_subnet
  internet_gateway = module.network.internet_gateway
  codepipeline_artifact_bucket = module.tf_backend.codepipeline_artifact_bucket
}