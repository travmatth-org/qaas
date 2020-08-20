# shellcheck disable=SC2148
remove-host() {
	ssh-keygen -R "$(make tf.ip)";
}

login-ec2() {
	ssh -i protected/faas_ec2.key ec2-user@"$(make tf.ip)";
}