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
    echo "deb [arch=${architecture}] https://mirrors.geekros.com/ focal main" | sudo tee /etc/apt/sources.list.d/armcnc.list
fi

sudo apt -y update

if [ ! -f "/usr/bin/linuxcnc" ]; then
    sudo apt install -y linuxcnc-uspace linuxcnc-uspace-dev
fi

if [ ! -f "/etc/ethercat.conf" ]; then
    sudo apt install -y ethercat-master libethercat-dev linuxcnc-ethercat
    MAC_ADDRESS=$(cat /etc/network/mac_address)
    sed -i "s/^MASTER0_DEVICE=\".*\"/MASTER0_DEVICE=\"$MAC_ADDRESS\"/" /etc/ethercat.conf
    sed -i "s/^DEVICE_MODULES=\".*\"/DEVICE_MODULES=\"generic\"/" /etc/ethercat.conf
    sudo touch /etc/udev/rules.d/99-ethercat.rules
    cat <<-EOF > /etc/udev/rules.d/99-ethercat.rules
KERNEL=="EtherCAT[0-9]", MODE="0777"
EOF
    sudo ldconfig
    # shellcheck disable=SC2046
    sudo hobot-sign-file $(modinfo -n ec_master)
    # shellcheck disable=SC2046
    sudo hobot-sign-file $(modinfo -n ec_generic)
    sudo systemctl enable ethercat.service
    sudo systemctl restart ethercat.service
fi

sudo apt install -y armcnc

echo "Install ARMCNC successfully"