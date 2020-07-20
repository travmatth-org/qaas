variable "github_repo" {
  description = "Github repository hosting the project source"
}

variable "codepipeline_artifact_bucket" {
  description = "S3 bucket for Codepipeline artifacts"
}

variable "codedeploy_app_name" {
  description = "Name of App within Codedeploy"
}

variable "codedeploy_group_name" {
  description = "Name of Group within Codedeploy"
}

variable "codebuild_logging_bucket" {
  description = "Bucket containing Codebuild log files"
}

variable "tf_state_bucket" {
  description = "Bucket containing tf state files"
}

variable "dynamodb_lock_state_table" {
  description = "DynamoDB table controlling terraform lock state"
}

variable "codebuild_project" {
  description = "Codebuild project"
}

variable "webhook_secret" {
  description = "Secret used by webhooks to authenticate"
}

variable "github_oauth_token" {
  description = "GitHub OAuth token"
}
