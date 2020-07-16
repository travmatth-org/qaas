module "codebuild" {
  source = "./modules/Codebuild"
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  tf_state_bucket              = var.tf_state_bucket
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
  codebuild_iam_role           = var.codebuild_iam_role
  buildspec_yml                = "../../build/ci/buildspec.yml"
}


# deploy compiled http server to ec2 instances via codedeploy
module "codedeploy" {
  source = "./modules/codedeploy"
}

module "codepipeline" {
  source                       = "./modules/codepipeline"
  repo_name                    = var.repo_name
  codebuild_project            = module.codebuild.codebuild_project
  dynamodb_lock_state_table    = var.dynamodb_lock_state_table
  tf_state_bucket              = var.tf_state_bucket
  codebuild_logging_bucket     = var.codebuild_logging_bucket
  codebuild_iam_role           = var.codebuild_iam_role
  codedeploy_app_name          = module.codedeploy.app_name
  codedeploy_group_name        = module.codedeploy.deployment_group_name
  codepipeline_artifact_bucket = var.codepipeline_artifact_bucket
  webhook_secret               = var.webhook_secret
  github_oauth_token           = var.github_oauth_token
}
