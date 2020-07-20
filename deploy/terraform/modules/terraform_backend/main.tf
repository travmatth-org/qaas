resource "aws_s3_bucket" "tf_state_bucket" {
  bucket = "faas-terraform-state-bucket-${var.aws_account_id}"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  versioning {
    enabled = true
  }

  tags = {
    Terraform = "true"
		FaaS = "true"
  }
}

resource "aws_s3_bucket" "codebuild_logging_bucket" {
  bucket = "faas-codebuild-logging-bucket-${var.aws_account_id}"

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

resource "aws_s3_bucket" "codepipeline_artifact_bucket" {
  bucket = "faas-codepipeline-artifact-bucket-${var.aws_account_id}"

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

resource "aws_dynamodb_table" "terraform_lock_state_dynamodb" {
  name = "faas-dynamodb-terraform-locking"
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
