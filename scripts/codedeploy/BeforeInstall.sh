#!/bin/bash -xe

# remove prev program 
sudo rm -f /usr/sbin/httpd

# install server assets
unzip /srv/assets.zip -d /srv
