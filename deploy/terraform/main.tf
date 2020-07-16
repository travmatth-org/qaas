terraform {
  required_version = ">=0.12.28"

  # backend "s3" {
  #   bucket         = "faas-terraform-tfstate"
  #   key            = "terraform.tfstate"
  #   region         = "us-west-1"
  #   dynamodb_table = "faas-terraform-locking"
  #   encrypt        = true
  # }
}

provider "aws" {
  version = "~>2.8"
  region = "us-west-1"

  # assume_role {
  #   role_arn     = "arn:aws:iam:XXXXXXXX:role/TerraformAssumedIamRole"
  #   session_name = terraform
  # }
}

provider "github" {
  version      = "~>2.9"
  token        = var.github_oauth_token
  individual   =  true
}

# Secret to authenticate in github / codepipeline webhook requests
resource "random_password" "github_secret" {
  length  = 22
  special = false
}

# Github address of repo
# data "github_repository" "faas" {
#   full_name = "travmatth/faas"
# }

# Store github oauth token in ssm parameter store
resource "aws_ssm_parameter" "github_oauth_token" {
  name        = "/faas/ci/github_oauth_token"
  description = "GitHub OAuth Token"
  type        = "SecureString"
  value       = "var.github_oauth_token"
  tags        = {
    FaaS      = "true"
  }
}

# Store webhook secret in ssm parameter store
resource "aws_ssm_parameter" "webhook_secret" {
  name        = "/faas/ci/webhook_secret"
  description = "CICD WebHook Secret"
  type        = "SecureString"
  value       = random_password.github_secret.result
  tags        = {
    FaaS      = "true"
  }
}


# Bootstrap backend
module "tf_backend" {
	source = "./modules/terraform_backend"
}

# CICD uses codepipeline to retrieve source from github,
# codebuild to test and construct artifacts, and codedeploy to deploy
# to ec2 servers
module "cicd" {
  source                       = "./modules/cicd"
  github_oauth_token           = var.github_oauth_token
  webhook_secret               = random_password.github_secret.result
  codebuild_iam_role           = module.tf_backend.codebuild_iam_role
  codebuild_logging_bucket     = module.tf_backend.codebuild_logging_bucket
  tf_state_bucket              = module.tf_backend.tf_state_bucket
  codepipeline_artifact_bucket = module.tf_backend.codepipeline_artifact_bucket
  dynamodb_lock_state_table    = module.tf_backend.dynamodb_lock_state_table
  buildspec_yml                = "${path.cwd}/build/ci/buildspec.yml"
  repo_name                    = "faas"#data.github_repository.faas.name
}

# security group, vpc, ec2 instance
module "http" {
  source = "./modules/http_ec2"
  
}