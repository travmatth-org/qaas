#!/bin/bash
set -eux pipefail
# if [[ $HOST =~ "*compute\.internal$" ]]; then
# 	exec 5>> >(logger -t $0)
# fi
# send script output to /tmp/AfterInstall.log for debugging
# exec >> /tmp/AfterInstall.log 2>&1

# give service appropriate permissions
sudo chmod 755 /usr/sbin/httpd

# allow service to run on privileged port
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/httpd

# enable faas service
sudo systemctl enable httpd