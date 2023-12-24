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
    sudo echo "deb [arch=${architecture}] https://mirrors.geekros.com/ focal main" | sudo tee /etc/apt/sources.list.d/armcnc.list

    sudo cat <<-EOF > /etc/apt/sources.list
deb http://mirrors.ustc.edu.cn/debian/ bookworm main contrib non-free non-free-firmware
deb-src http://mirrors.ustc.edu.cn/debian/ bookworm main contrib non-free non-free-firmware
deb http://mirrors.ustc.edu.cn/debian/ sid main contrib non-free non-free-firmware
EOF

    if [ ! -f /etc/sudoers.d/armcnc ]; then
        echo "armcnc ALL=(ALL:ALL) ALL" | sudo tee /etc/sudoers.d/armcnc
        sudo chmod 0440 /etc/sudoers.d/armcnc
    fi

    sudo apt -y update

    sudo apt install -y python3.11 python3-pip

    sudo cat <<-EOF > /etc/apt/sources.list
deb http://mirrors.ustc.edu.cn/debian/ bookworm main contrib non-free non-free-firmware
deb-src http://mirrors.ustc.edu.cn/debian/ bookworm main contrib non-free non-free-firmware
EOF

    sudo sed -i 's/#\?PermitRootLogin .*/PermitRootLogin yes/' /etc/ssh/sshd_config
    sudo sed -i 's/#\?PubkeyAuthentication .*/PubkeyAuthentication yes/' /etc/ssh/sshd_config

    sudo cat <<-EOF > /etc/gdm3/daemon.conf
# GDM configuration storage
#
# See /usr/share/gdm/gdm.schemas for a list of available options.

[daemon]
# Uncomment the line below to force the login screen to use Xorg
#WaylandEnable=false

# Enabling automatic login
AutomaticLoginEnable = true
AutomaticLogin = armcnc

# Enabling timed login
#  TimedLoginEnable = true
#  TimedLogin = user1
#  TimedLoginDelay = 10

[security]

[xdmcp]

[chooser]

[debug]
# Uncomment the line below to turn on debugging
# More verbose logs
# Additionally lets the X server dump core if it crashes
#Enable=true
EOF

    sudo cp -aRf /etc/skel/. /root/

    wired_card=$(nmcli device status | grep 'ethernet' | awk '{print $1}')

    if [ ! -f "/etc/set_mac_address.sh" ]; then
        sudo touch /etc/network/mac_address
        sudo cat <<EOF > /etc/set_mac_address.sh
#!/bin/bash
mac_file=/etc/network/mac_address

if [ -s \${mac_file} ] && [ -f \${mac_file} ]; then
    ifconfig $wired_card down
    ifconfig $wired_card hw ether \$(cat \${mac_file})
    ifconfig $wired_card up
else
    openssl rand -hex 6 | sed 's/\(..\)/\1:/g; s/:$//' > \$mac_file
    ifconfig $wired_card down
    ifconfig $wired_card hw ether \$(cat \$mac_file)
    ifconfig $wired_card up
fi

if [ -e /etc/ethercat.conf ]; then
    systemctl restart ethercat.service
fi
EOF

        sudo cat <<EOF > /etc/network/interfaces
# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d
auto $wired_card
iface $wired_card inet static
        pre-up /etc/set_mac_address.sh
        address 192.168.1.10
        netmask 255.255.255.0
        #network
        #broadcast
        gateway 192.168.1.1
        metric 700
EOF

        sudo chmod +x /etc/set_mac_address.sh
    fi
fi

sudo apt -y update
sudo apt install -y vim git wget curl make cmake net-tools htop geany chromium

if [ ! -f "/usr/bin/linuxcnc" ]; then
    sudo apt install -y linuxcnc-uspace linuxcnc-uspace-dev
fi

if [ ! -f "/etc/ethercat.conf" ]; then
    sudo apt install -y ethercat-master libethercat-dev linuxcnc-ethercat
    sudo openssl rand -hex 6 | sed 's/\(..\)/\1:/g; s/:$//' > /etc/network/mac_address
    MAC_ADDRESS=$(cat /etc/network/mac_address)
    sed -i "s/^MASTER0_DEVICE=\".*\"/MASTER0_DEVICE=\"$MAC_ADDRESS\"/" /etc/ethercat.conf
    sed -i "s/^DEVICE_MODULES=\".*\"/DEVICE_MODULES=\"generic\"/" /etc/ethercat.conf
    sudo touch /etc/udev/rules.d/99-ethercat.rules
    cat <<-EOF > /etc/udev/rules.d/99-ethercat.rules
KERNEL=="EtherCAT[0-9]", MODE="0777"
EOF
    sudo ldconfig
    sudo /etc/set_mac_address.sh
    sudo systemctl enable ethercat.service
    sudo systemctl restart ethercat.service
fi

sudo apt install -y armcnc

echo "Install ARMCNC successfully"