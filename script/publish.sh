#!/bin/bash

set -e

version=$1
# shellcheck disable=SC2037
architecture=$(dpkg --print-architecture)

echo "$PWD $architecture $version"
lsb_release -a

sudo apt install -y dpkg-dev gpg

cd ../

# shellcheck disable=SC2035
sudo rm -rf *.deb

if [ ! -f "/usr/lib/linuxcnc/modules/armcncio.so" ]; then
    echo "Not build armcncio.so"
    exit 0
else
    cp /usr/lib/linuxcnc/modules/armcncio.so debian/usr/lib/linuxcnc/modules/
fi

if [ ! -f "/root/desktop/template/dist/index.html" ]; then
    echo "Not build desktop"
    exit 0
else
    cp -r /root/desktop/template/dist/* debian/opt/armcnc/touch/
fi

sudo chmod +x debian/DEBIAN/*

sudo rm -rf debian/DEBIAN/control
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

cd ./armcnc
export GO111MODULE=on && export GOPROXY=https://proxy.golang.com.cn,direct && sudo rm -rf main
sudo -E /usr/local/go/bin/go build main.go
sudo cp main ../debian/usr/local/bin/armcnc
sudo rm -rf main

cd ../

sudo dpkg --build debian/ && dpkg-name debian.deb

sudo rm -rf debian/usr/local/bin/armcnc
sudo rm -rf debian/usr/lib/linuxcnc/modules/armcncio.so
sudo rm -rf debian/opt/armcnc/touch/index.html
sudo rm -rf debian/opt/armcnc/touch/favicon.ico
sudo rm -rf debian/opt/armcnc/touch/assets
sudo rm -rf debian/opt/armcnc/touch/monacoeditorwork
sudo rm -rf debian/opt/armcnc/touch/static

echo "Publish successfully"