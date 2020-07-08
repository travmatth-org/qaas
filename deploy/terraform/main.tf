locals {
  codepipeline_artifact_bucket_name = "faas-codebuild-artifact-bucket"
  codebuild_iam_role_name = "FaaSCodeBuildIamRole"
	s3_tfstate_bucket = "faas-terraform-tfstate"
  region = "us-west-1"
	dynamo_db_table_name = "faas-terraform-locking"
}
 
terraform {
  required_version = ">=0.12.28"

  # backend "s3" {
  #   bucket         = locals.s3_tfstate_bucket
  #   key            = "terraform.tfstate"
  #   region         = locals.region
  #   dynamodb_table = locals.dynamo_db_table_name
  #   encrypt        = true
  # }
}

provider "aws" {
  region = locals.region

  # assume_role {
  #   role_arn     = "arn:aws:iam:XXXXXXXX:role/TerraformAssumedIamRole"
  #   session_name = terraform
  # }
}

provider "github" {
  token = "${var.github_token}"
}
 
# secret to authenticate in github / codepipeline webhook requests
resource "random_password" "github_secret" {
  length  = 22
  special = false
}

# github address of repo
data "github_repository" "faas" {
  full_name = "travmatth/faas"
}

# store github oauth token in ssm parameter store
data "aws_ssm_parameter" "github_oauth_token" {
  name = "/ci/github_oauth_token"
  description = "GitHub OAuth Token"
  type = "SecureString"
  value = "${var.github_oauth_token}"

  tags = {
    FaaS = "true"
  }
}

# store webhook secret in ssm parameter store 
data "aws_ssm_parameter" "webhook_secret" {
  name = "/ci/webhook_secret"
  description = "CICD WebHook Secret"
  type = "SecureString"
  value = "${random_password.github_secret.result}"

  tags = {
    FaaS = "true"
  }
}

# Bootstrap backend

# s3 buckets for tf state, codepipeline logs
module "backend" {
	source = "./modules/terraform_backend"
	s3_tfstate_bucket = locals.s3_tfstate_bucket
	s3_logging_bucket_name = "faas-codebuild-logging-bucket"
	codepipeline_artifact_bucket_name = locals.codepipeline_artifact_bucket_name
}

# dynamodb for tf remote state locking
module "dynamodb_table" {
	source = "./modules/dynamodb-table"
	dynamo_db_table_name = locals.dynamo_db_table_name
}

# create codepipeline: github -> codebuild -> codedeploy

# role codebuild will assume
module "iam_codebuild_role" {
	source = "./modules/iam-codebuild-role"
	codebuild_iam_role_name = locals.codebuild_iam_role_name
}

# policy detailing resource access for codebuild role
module "iam_codebuild_role_policy" {
	source = "./modules/iam-codebuild-role-policy"
	codebuild_iam_role_policy_name = "CodeBuildIamRolePolicy"
  codebuild_iam_role_name = locals.codebuild_iam_role_name
  logging_bucket_arn = module.backend.logging_bucket.arn
  state_bucket_arn = module.backend.state_bucket.arn
  codepipeline_artifact_bucket_arn = module.backend.codepipeline_artifact_bucket.arn
  codebuild_iam_role_arn = module.iam_codebuild_role.codebuild_iam_role_arn
  dynamodb_lock_state_table_arn = module.dynamodb_table.lock_state.arn
}

# codebuild instance for building and testing code
module "codebuild" {
  source                 = "./modules/codebuild"
  project_name           = "TerraformPlan"
  s3_logging_bucket_id   = module.backend.logging_bucket_id
  s3_logging_bucket      = module.backend.logging_bucket
  codebuild_iam_role_arn = module.iam_codebuild_role.codebuild_iam_role_arn
}

# deploy compiled http server to ec2 instances via codedeploy
module "codedeploy" {
  source = "./modules/codedeploy"
  app_name = "faas"
}

# role codedeploy will assume
module "codedeploy_role" {
  source = "./modules/iam-codedeploy-role"
}

# attach codedeploy policy to role
module "codedeploy_role_policy_attachment" {
  source = "./modules/iam-codedeploy-role-policy"
  codedeploy_role_name = module.codedeploy_role.codedeploy_role_name
}

# deployment groups of codedeploy
module "codedeploy_deployment_group" {
  app_name = module.codedeploy.name
  deployment_group_name = "ec2 group"
  service_role_arn = module.codedeploy_role.codedeploy_role_arn
}

# role codepipeline to assume 
module "codepipeline_role" {
  codepipeline_role_name = "TerraformCodePipelineIamRole"
}

# policy detailing resource access for codepipeline role
module "codepipeline_role_policy" {
  codepipeline_role_policy_name = "TerraformCodePipelineIamRolePolicy"
  codepipeline_role_id = module.codepipeline_role.id
}

# cicd server that listens for github hook, tests, builds app & deploys to ec2
module "codepipeline" {
  source                       = "./modules/codepipeline"
  codepipeline_name            = "FaaSCodePipeline"
  codepipeline_role_arn        = module.codepipeline_role.codebuild_project_arn
  codepipeline_artifact_bucket = module.backend.artifact_bucket
  project_name                 = module.codebuild.project_name
  oauth_token                  = data.aws_ssm_parameter.github_token.value
  codedeploy_app_name          = module.codedeploy.app_name
  codedeploy_group_name        = module.codedeploy.group_name
  repo_name                    = data.github_repository.faas.name
  oauth_token                  = data.github_oauth_token.value
}

# direct codepipeline to listen for github hooks
module "codepipeline_webhook" {
  source            = "./modules/codepipeline-webhook"
  codepipeline_name = module.codepipeline.name
  webhook_url       = data.webhook_secret.value
}

# direct github to post to codepipeline on changes to master
module "github_webhook" {
  source         = "./modules/github-webhook"
  webhook_url    = data.github_oauth_token.url
  webhook_secret = data.webhook_secret.value
}