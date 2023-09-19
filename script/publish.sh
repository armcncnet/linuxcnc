#!/bin/bash

set -e

version=$1
# shellcheck disable=SC2037
architecture=$(dpkg --print-architecture)

echo "$PWD $architecture $version"
lsb_release -a

sudo apt install -y dpkg-dev gpg

cd ../

sudo chmod +x debian/DEBIAN/*
find ./debian -type f -name ".gitkeep" -exec rm -f {} +

sudo touch debian/DEBIAN/control && sudo chmod +x debian/DEBIAN/control
sudo sh -c 'echo "Package: armcnc" >> debian/DEBIAN/control'
version_str="Version: $version"
export version_str
sudo -E sh -c 'echo $version_str >> debian/DEBIAN/control'
sudo sh -c 'echo "Maintainer: ARMCNC <admin@geekros.com>" >> debian/DEBIAN/control'
sudo sh -c 'echo "Homepage: https://www.armcnc.net" >> debian/DEBIAN/control'
architecture_str="Architecture: $architecture"
export architecture_str
sudo -E sh -c 'echo $architecture_str >> debian/DEBIAN/control'
sudo sh -c 'echo "Installed-Size: 5000" >> debian/DEBIAN/control'
sudo sh -c 'echo "Section: utils" >> debian/DEBIAN/control'
sudo sh -c 'echo "Description: armcnc" >> debian/DEBIAN/control'

export GO111MODULE=on && export GOPROXY=https://goproxy.io && sudo rm -rf main
sudo -E /usr/local/go/bin/go build main.go
sudo cp main debian/usr/local/bin/armcnc
sudo rm -rf main

sudo dpkg --build debian/ && dpkg-name debian.deb

echo "Publish successfully"