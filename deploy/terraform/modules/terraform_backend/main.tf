# Terraform state bucket
resource "aws_s3_bucket" "tf_state_bucket" {
  bucket = "faas-terraform-state-bucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  # Prevents Terraform from destroying or replacing this object
  lifecycle {
    prevent_destroy = true
  }

  # Tells AWS to keep a version history of the state file
  versioning {
    enabled = true
  }

  tags = {
    Terraform = "true"
		FaaS = "true"
  }
}

# CodeBuild logging bucket
resource "aws_s3_bucket" "codebuild_logging_bucket" {
  bucket = "faas-codebuild-logging-bucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Terraform = "true"
    Logging = "true"
  }
}

# S3 bucket for CodePipeline artifact storage
resource "aws_s3_bucket" "codepipeline_artifact_bucket" {
  bucket = "faas-codepipeline-artifact-bucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Terraform = "true"
    CodePipelineArtifacts = "true"
  }
}

# DynamoDB to use for terraform state locking
resource "aws_dynamodb_table" "terraform_lock_state_dynamodb" {
  name = "faas-dynamodb-terraform-locking"

  # Pay per request is cheaper for low-i/o applications, like our TF lock state
  billing_mode = "PAY_PER_REQUEST"

  # Hash key is required, and must be an attribute
  hash_key = "LockID"

  # Attribute LockID is required for TF to use this table for lock state
  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
  }
}

resource "aws_iam_policy" "tf_iam_assumed_policy" {
  name = "TerraformAssumedIamPolicy"

  policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "AllowAllPermissions",
          "Effect": "Allow",
          "Action": [
            "*"
          ],
          "Resource": "*"
        }
      ]
    }
    EOF

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_iam_role" "codebuild_iam_role" {
  name = "CodeBuildIamRole" 
  assume_role_policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Service": "codebuild.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }
    EOF

  tags = {
    Terraform = "true"
    FaaS      = "true"
  }
}

resource "aws_iam_role" "tf_iam_assumed_role" {
  name = "TerraformAssumedIamRole"
  assume_role_policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "AWS": "${aws_iam_role.codebuild_iam_role.arn}"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }
    EOF

  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
  }
}

resource "aws_iam_role_policy_attachment" "tf_iam_attach_assumed_role_to_permissions_policy" {
  role       = aws_iam_role.tf_iam_assumed_role.name
  policy_arn = aws_iam_policy.tf_iam_assumed_policy.arn

  lifecycle {
    prevent_destroy = true
  }
}