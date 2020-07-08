# Terraform state bucket
resource "aws_s3_bucket" "state_bucket" {
  bucket = var.s3_tfstate_bucket

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

# Build an AWS S3 bucket for logging
resource "aws_s3_bucket" "logging_bucket" {
  bucket = var.s3_logging_bucket_name

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
  bucket = var.codepipeline_artifact_bucket_name

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
