#!/bin/bash -xe

# send script output to /tmp/ApplicationStop.log for debugging
exec >> /tmp/ApplicationStop.log 2>&1

if ( sudo systemctl status httpd | grep active ); then \
	sudo systemctl stop httpd;
fi
