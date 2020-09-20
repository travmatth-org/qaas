variable "aws_account_id" {
  description = "Account ID of AWS account"
}

resource "aws_s3_bucket" "tf_state_bucket" {
  acl    = "private"
  bucket = "qaas-terraform-state-bucket-${var.aws_account_id}"

  # lifecycle {
  #   prevent_destroy = true
  # }

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
    qaas      = "true"
  }
}

output "tf_state_bucket" {
  value = aws_s3_bucket.tf_state_bucket
}

resource "aws_s3_bucket" "codebuild_logging_bucket" {
  bucket = "qaas-codebuild-logging-bucket-${var.aws_account_id}"

  # lifecycle {
  #   prevent_destroy = true
  # }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Terraform = "true"
    Logging   = "true"
  }
}

output "codebuild_logging_bucket" {
  value = aws_s3_bucket.codebuild_logging_bucket
}

resource "aws_s3_bucket" "codepipeline_artifact_bucket" {
  bucket = "qaas-codepipeline-artifact-bucket-${var.aws_account_id}"

  # lifecycle {
  #   prevent_destroy = true
  # }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Terraform             = "true"
    CodePipelineArtifacts = "true"
  }
}

output "codepipeline_artifact_bucket" {
  value = aws_s3_bucket.codepipeline_artifact_bucket
}

resource "aws_dynamodb_table" "terraform_lock_state_dynamodb" {
  name         = "qaas-dynamodb-terraform-locking"
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
    qaas      = "true"
  }
}

output "dynamodb_lock_state_table" {
  value = aws_dynamodb_table.terraform_lock_state_dynamodb
}
