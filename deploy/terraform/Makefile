TF_PLAN	:= infra.tfplan

init:
	terraform init;

plan:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN);

apply:
	terraform apply $(TF_PLAN);

destroy:
	terraform destroy -var-file=".tfvars";

show:
	terraform show;

# makefile phony target
.PHONY: init plan apply destroy destroy.ec2 show show.eip
