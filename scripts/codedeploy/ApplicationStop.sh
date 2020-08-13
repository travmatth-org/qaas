#!/bin/bash
set -eux pipefail

if systemctl is-active --quiet httpd; then
  sudo systemctl stop httpd.socket httpd
fi