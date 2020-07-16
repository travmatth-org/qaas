resource "aws_iam_role" "codepipeline_role" {
  name = "TerraformCodePipelineIamRole"

  assume_role_policy =<<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Service": "codepipeline.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }
    EOF

  tags = {
    FaaS = "true"
  }
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name = "TerraformCodePipelineIamRolePolicy"
  role = aws_iam_role.codepipeline_role.id
  policy = <<-POLICY
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Resource": [
            "*"
          ],
          "Action": [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "s3:*"
          ],
          "Resource": [
            "${var.codebuild_logging_bucket.arn}",
            "${var.codebuild_logging_bucket.arn}/*",
            "${var.tf_state_bucket.arn}",
            "${var.tf_state_bucket.arn}/*",
            "arn:aws:s3:::codepipeline-us-east-1*",
            "arn:aws:s3:::codepipeline-us-east-1*/*",
            "${var.codepipeline_artifact_bucket.arn}",
            "${var.codepipeline_artifact_bucket.arn}/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "dynamodb:*"
          ],
          "Resource": "${var.dynamodb_lock_state_table.arn}"
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:Get*",
            "iam:List*"
          ],
          "Resource": "${var.codebuild_iam_role.arn}"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Resource": "${var.codebuild_iam_role.arn}"
        }
      ]
    }
    POLICY
    }

resource "aws_codepipeline" "codepipeline" {
  name     = "FaaSCodePipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = var.codepipeline_artifact_bucket.bucket
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
      configuration = {
        Owner = "travmatth"
        Repo = "faas"
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
        ProjectName = var.codebuild_project.name
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
    name = "Deploy"

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

resource "aws_codepipeline_webhook" "faas" {
	name = "faas-codepipeline-webhook"
	authentication = "GITHUB_HMAC"
	target_action = "Source"
	target_pipeline = aws_codepipeline.codepipeline.name
	authentication_configuration {
		secret_token = var.webhook_secret
	}
	filter {
		json_path = "$.ref"
		match_equals = "refs/heads/{Branch}"
	}
}

resource "github_repository_webhook" "test" {
  repository = var.repo_name
  configuration {
    url          = aws_codepipeline_webhook.faas.url
    content_type = "form"
    insecure_ssl = true
    secret       = var.webhook_secret
  }
  events = ["push"]
}
