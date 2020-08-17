data "aws_iam_policy_document" "assume_role" {
	statement {
		effect  = "Allow"
		actions = ["sts:AssumeRole"]

		principals {
			type = "Service"
			identifiers = [
				"ec2.amazonaws.com",
				"codedeploy.amazonaws.com",
				"s3.amazonaws.com",
			]
		}
	}
}

resource "aws_iam_role" "faas" {
  name  = "faas_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

output "faas_iam_role" {
	value = aws_iam_role.faas
}

data "aws_iam_policy_document" "faas_cicd_policy" {
	statement {
		sid = "AWS"
		effect = "Allow"
		actions = [
			"s3:*",
			"logs:*",
			"elasticloadbalancing:*",
			"codedeploy:*",
		]
		resources = ["*"]
	}
	
	# allow cloudwatch-agent to collect and send logs, metrics
	# https://github.com/SummitRoute/aws_managed_policies/blob/master/policies/CloudWatchAgentServerPolicy
	statement {
		effect  = "Allow"
		actions = [
			"cloudwatch:PutMetricData",
			"ec2:DescribeVolumes",
			"ec2:DescribeTags",
			"logs:PutLogEvents",
			"logs:DescribeLogStreams",
			"logs:DescribeLogGroups",
			"logs:CreateLogStream",
			"logs:CreateLogGroup"
		]
		resources = ["*"]
	}

	# allow x-ray to send traces
	# https://docs.aws.amazon.com/xray/latest/devguide/security_iam_id-based-policy-examples.html
	statement {
		effect  = "Allow"
		actions = [
			"xray:PutTraceSegments",
			"xray:PutTelemetryRecords",
			"xray:GetSamplingRules",
			"xray:GetSamplingTargets",
			"xray:GetSamplingStatisticSummaries"
		]
		resources = ["*"]
	}
}

resource "aws_iam_policy" "cicd" {
  name  = "faas_role_policy"
  description = "A policy for CodeDeploy"
  policy = data.aws_iam_policy_document.faas_cicd_policy.json
  depends_on = [aws_iam_role.faas]
}

resource "aws_iam_role_policy_attachment" "attach_custom_policy" {
  depends_on = [aws_iam_policy.cicd]
  role       = aws_iam_role.faas.name
  policy_arn = aws_iam_policy.cicd.arn
}

resource "aws_iam_instance_profile" "faas_service" {
  name = "faas_service"
  path = "/"
  role = aws_iam_role.faas.name
}