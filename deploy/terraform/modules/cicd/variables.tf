variable "codebuild_logging_bucket" {
  description = "S3 bucket used by CodeBuild for logging"
}

variable "tf_state_bucket" {
  description = "S3 bucket used by CodeBuild for state management"
}

variable "codepipeline_artifact_bucket" {
  description = "ARN for S3 bucket used by CodePipeline for artifacts"
}

variable "dynamodb_lock_state_table" {
  description = "DynamoDB table controlling lock state"
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
