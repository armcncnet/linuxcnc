#!/bin/bash

set -e

uname -a

lsb_release -a

architecture=$(dpkg --print-architecture)

sudo apt install -y nodejs npm

sudo npm install -g yarn

if [ ! -d "/usr/local/go/bin/" ]; then
    sudo wget -q https://studygolang.com/dl/golang/go1.19.4.linux-"${architecture}".tar.gz && sudo tar -C /usr/local -xzf go1.19.4.linux-"${architecture}".tar.gz && sudo rm -rf go1.19.4.linux-"${architecture}".tar.gz
    sudo sh -c 'echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile'
    # shellcheck disable=SC1090
    source /etc/profile && source ~/.bashrc
fi

echo "Install development environment successfully"