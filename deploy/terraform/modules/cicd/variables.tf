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

variable "buildspec_yml" {
	description = "Location of the buildspec file to use"
}

variable "repo_name" {
	description = "Name of the repository"
}

variable "codebuild_iam_role" {
  description = "IAM role to use for Codebuild"
}

variable "webhook_secret" {
  description = "Secret used by webhooks to authenticate"
}

variable "github_oauth_token" {
  description = "GitHub OAuth token"
}
