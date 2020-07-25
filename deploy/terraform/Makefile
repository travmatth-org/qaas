TF_PLAN	:= infra.tfplan

init:
	terraform init;

plan:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN);

apply:
	terraform apply $(TF_PLAN);

destroy:
	terraform destroy -var-file=".tfvars";

destroy.ec2:
	terraform destroy -target "module.ec2" -var-file=.tfvars;

destroy.cicd:
	terraform destroy -target "module.cicd" -var-file=.tfvars;

show:
	terraform show;

show.eip:
	terraform state show 'module.ec2.aws_eip.faas_eip';

# makefile phony target
.PHONY: init plan apply destroy destroy.ec2 show show.eip
