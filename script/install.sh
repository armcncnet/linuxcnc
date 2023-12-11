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

cat <<'EOF' > /etc/set_mac_address.sh
#!/bin/bash
mac_file=/etc/network/mac_address

if [ -s ${mac_file} ] && [ -f ${mac_file} ]; then
    ifconfig eth0 down
    ifconfig eth0 hw ether $(cat ${mac_file})
    ifconfig eth0 up
else
    openssl rand -rand /dev/urandom:/sys/class/socinfo/soc_uid -hex 6 | sed -e 's/../&:/g;s/:$//' -e 's/^\(.\)[13579bdf]/\10/' > $mac_file
    ifconfig eth0 down
    ifconfig eth0 hw ether $(cat ${mac_file})
    ifconfig eth0 up
fi

if [ -e /etc/ethercat.conf ]; then
    systemctl restart ethercat.service
fi
EOF

sudo apt -y update && sudo apt -y upgrade
if [ ! -f "/usr/bin/linuxcnc" ]; then
    sudo apt install -y linuxcnc-uspace linuxcnc-uspace-dev
fi
sudo apt install -y armcnc

echo "Install ARMCNC successfully"