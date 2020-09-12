TF_PLAN	:= infra.tfplan

init:
	terraform init;

plan:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN);

plan.cicd:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN) \
		-target=module.tf_backend \
		-target=module.cicd;

plan.network:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN) \
		-target=module.tf_backend \
		-target=module.cicd \
		-target=module.network;

plan.asg:
	terraform plan -var-file=".tfvars" -out $(TF_PLAN) \
		-target=module.tf_backend \
		-target=module.cicd \
		-target=module.network \
		-target=module.asg;

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

ip:
	@terraform output ec2_ip

# makefile phony target
.PHONY: init plan apply destroy destroy.ec2 show show.eip
