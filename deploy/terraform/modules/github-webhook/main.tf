resource "github_repository_webhook" {
	repository = var.repo_name
	name = "faas-cicd"

	configuration {
		url = var.webhook_url
		content_type = "form"
		insecure_ssl = true
		secret = var.webhook_secret
	}

	events = ["push"]
}