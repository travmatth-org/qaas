module "codebuild" {
  source = "./modules/Codebuild"
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  tf_state_bucket              = var.tf_state_bucket
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
}


module "codedeploy" {
  source = "./modules/codedeploy"
}

module "codepipeline" {
  source                       = "./modules/codepipeline"
  github_repo                  = var.github_repo
  codebuild_project            = module.codebuild.codebuild_project
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
  tf_state_bucket              = var.tf_state_bucket
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  codedeploy_app_name          = module.codedeploy.app_name
  codedeploy_group_name        = module.codedeploy.deployment_group_name
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  webhook_secret               = var.webhook_secret
  github_oauth_token           = var.github_oauth_token
}
