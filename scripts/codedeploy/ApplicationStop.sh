#!/bin/bash -xe

if ( sudo systemctl status httpd | grep active ); then \
	sudo systemctl stop httpd;
fi