resource "awas_codedeploy_deploymentgroup" "deploy" {
	app_name 			  = "${var.app_name}"
	deployment_group_name = "${var.deployment_group_name}"
	service_role_arn 	  = "${var.service_role_arn}"

	ec2_tag_set {
		type  = "KEY_AND_VALUE"	
		key   = "FaaS"
		value = "Service"
	}
}