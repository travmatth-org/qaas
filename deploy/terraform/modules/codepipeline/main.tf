resource "aws_codepipeline" "codepipeline" {
  name     = var.codepipeline_name
  role_arn = var.codepipeline_role_arn

  artifact_store {
    location = var.codepipeline_artifact_bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name = "Git"
      category = "Source"
      owner = "ThirdParty"
      provider = "GitHub"
      version = "1"
      output_artifacts = ["source_artifact"]
      configuration {
        Owner = "travmatth"
        Repo = var.repo_name
        Branch = "master"
        OAuthToken = var.github_oauth_token
        PollForSourceChanges = false
      }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_artifact"]
      output_artifacts = ["build_artifact"]
      version          = "1"

      configuration = {
        ProjectName = var.codebuild_project_name
      }
    }
  }

  stage {
    name = "Manual_Approval"

    action {
      name = "Manual-Approval"
      category = "Approval"
      owner = "AWS"
      provider = "Manual"
      version = "1"
    }
  }

  # https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html
  stage {
    name = Deploy

    action {
      name = "Deploy"
      category = "Deploy"
      owner = "AWS"
      provider = "CodeDeploy"
      version = "1"
      input_artifacts = ["build_artifact"]
      configuration = {
        ApplicationName = var.codedeploy_app_name
        DeploymentGroupName = var.codedeploy_group_name
        AppSpecTemplateArtifact = "build_artifact"
      }
    }
  }
}
