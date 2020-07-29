#!/bin/bash
echo "BeforeInstall.sh" | systemd-cat
# set -eux pipefail
# # send script output to /tmp/BeforeInstall.log for debugging
# # exec >> /tmp/BeforeInstall.log 2>&1

# # remove prev program 
# # TODO: Needed?
# # sudo rm -f /usr/sbin/httpd
# # sudo rm -f /usr/lib/systemd/system/httpd.service

# # install server assets
# sudo unzip -o /srv/assets.zip -d /srv
# # shellcheck disable=SC2046
# sudo chmod 0444 $(find /srv -type f)
# # clean
# sudo rm /srv/assets.zip
