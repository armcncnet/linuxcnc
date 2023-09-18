#!/bin/bash

set -e

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
