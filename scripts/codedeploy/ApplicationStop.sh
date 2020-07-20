#!/bin/bash

set -ex

if ( sudo systemctl status launchpad | grep active ); then \
	sudo systemctl stop launchpad;
fi