#!/bin/bash

set -e

uname -a

lsb_release -a

architecture=$(dpkg --print-architecture)

if [ ! -f "/etc/apt/sources.list.d/armcnc.list" ]; then
    sudo touch /etc/apt/sources.list.d/armcnc.list
    sudo wget -q -O /tmp/Release.gpg https://mirrors.geekros.com/dists/focal/Release.public
    # shellcheck disable=SC2024
    sudo gpg --dearmor < /tmp/Release.gpg > /etc/apt/trusted.gpg.d/geekros-archive.gpg
    echo "deb [arch=${architecture}] https://mirrors.geekros.com/ bookworm main contrib non-free non-free-firmware" | sudo tee /etc/apt/sources.list.d/armcnc.list
fi

sudo apt -y update && sudo apt -y upgrade
sudo apt install linuxcnc-uspace=2.9.0~pre1+git20230208.f1270d6ed7-1 linuxcnc-uspace-dev=2.9.0~pre1+git20230208.f1270d6ed7-1

echo "Install ARMCNC successfully"