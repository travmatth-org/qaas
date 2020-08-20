remove_host() {
	ssh-keygen -R $(make tf.ip);
}