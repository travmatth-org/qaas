resource "aws_codepipeline_webhook" "faas" {
	name = "faas-codepipeline-webhook"
	authentication = "GITHUB_HMAC"
	target_action = "Source"
	target_pipeline = var.codepipeline_name

	authentication_configuration {
		secret_token = var.webhook_secret
	}

	filter {
		json_path = "$.ref"
		match_equals = "refs/heads/{Branch}"
	}
}