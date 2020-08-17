data "aws_iam_policy_document" "codebuild_role_policy" {
	statement {
		actions = ["sts:AssumeRole"]

		principals {
			type = "Service"
			identifiers = ["codebuild.amazonaws.com"]
		}
	}
}

resource "aws_iam_role" "codebuild_role" {
	name = "FaaSCoddeBuildIamRole"
	assume_role_policy = data.aws_iam_policy_document.codebuild_role_policy.json
}

data "aws_iam_policy_document" "codebuild_policy" {
	statement {
		effect  = "Allow"
		actions = ["s3:*"]
		resources = [
			var.codebuild_logging_bucket.arn,
			"${var.codebuild_logging_bucket.arn}/*",
			var.codepipeline_artifact_bucket.arn,
			"${var.codepipeline_artifact_bucket.arn}/*",
			"arn:aws:s3:::codepipeline-us-west-1*",
			"arn:aws:s3:::codepipeline-us-west-1*/*",
		]
	}

	statement {
		effect  = "Allow"
		resources = ["*"]
		actions = [
			"logs:CreateLogGroup",
			"logs:CreateLogStream",
			"logs:PutLogEvents",
		]
	}
}

resource "aws_iam_policy" "codebuild_policy" {
	name = "CodeBuildIamPolicy"
	path = "/service-role/"
	policy = data.aws_iam_policy_document.codebuild_policy.json
}

resource "aws_iam_policy_attachment" "codebuild_policy_attachment" {
	name = "CodeBuildPolicyAttachment"
	policy_arn = aws_iam_policy.codebuild_policy.arn
	roles = [aws_iam_role.codebuild_role.id]
}