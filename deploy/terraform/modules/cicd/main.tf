variable "codebuild_logging_bucket" {}
variable "tf_state_bucket" {}
variable "user_name" {}
variable "account_id" {}
variable "codepipeline_artifact_bucket" {}
variable "dynamodb_lock_state_table" {}

module "codebuild" {
  source                       = "./modules/Codebuild"
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
  user_name                    = var.user_name
  account_id                   = var.account_id
}

module "codedeploy" {
  source = "./modules/codedeploy"
}

variable "github_repo" {
  description = "Github repository hosting the source code of the project"
}

variable "webhook_secret" {
  description = "Secret used by webhooks to authenticate"
}

variable "github_oauth_token" {
  description = "GitHub OAuth token"
}

module "codepipeline" {
  source                       = "./modules/codepipeline"
  github_repo                  = var.github_repo
  codebuild_project            = module.codebuild.codebuild_project
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  codedeploy                   = module.codedeploy
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  webhook_secret               = var.webhook_secret
  github_oauth_token           = var.github_oauth_token
  account_id                   = var.account_id
}
