variable "codepipeline_name" {
  description = "CodePipeline Name"
}

variable "codepipeline_role_arn" {
  description = "ARN for CodePipeline Role"
}

variable "codepipeline_artifact_bucket" {
  description = "S3 bucket for CodePipeline artifacts"
}

variable "codebuild_project_name" {
  description = "CodeBuild project name"
}

variable "github_oauth_token" {
  description = "OAuth token for GitHub"
}

variable "codedeploy_app_name" {
  description = ""
}

variable "codedeploy_group_name" {
  description = ""
}
