#!/bin/bash
set -eux pipefail

sudo systemctl daemon-reload
sudo systemctl start httpd
