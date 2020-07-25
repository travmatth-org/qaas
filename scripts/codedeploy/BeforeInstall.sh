#!/bin/bash -eux pipefail

# send script output to /tmp/BeforeInstall.log for debugging
exec >> /tmp/BeforeInstall.log 2>&1

# remove prev program 
sudo rm -f /usr/sbin/httpd
sudo rm -f /usr/lib/systemd/system/httpd.service

# install server assets
unzip /srv/assets.zip -d /srv
