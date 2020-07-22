#!/bin/bash -xe

# give service appropriate permissions
sudo chmod 755 /usr/sbin/httpd

# allow service to run on privileged port
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/httpd