variable "repo_name" {
	description = "Name of the GitHub repo to create hook for"
}

variable "webhook_url" {
	description = "URL of the GitHub WebHook"
}

variable "webhook_secret" {
	description = "Secret for the GitHub Webhook"
}
