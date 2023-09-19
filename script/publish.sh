#!/bin/bash

set -e

version=$1
# shellcheck disable=SC2037
architecture=$(dpkg --print-architecture)

cd ../armcnc && echo "$PWD $architecture $version"
lsb_release -a

sudo apt install -y dpkg-dev gpg

export GO111MODULE=on && export GOPROXY=https://goproxy.io && sudo rm -rf main
sudo -E /usr/local/go/bin/go build main.go
sudo cp main debian/usr/local/bin/armcnc
sudo rm -rf main

echo "Publish successfully"