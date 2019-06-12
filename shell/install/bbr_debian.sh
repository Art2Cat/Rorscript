#!/bin/bash
#
current_version=$(me -r | grep -Eo "[0-9\.]+" | xargs | awk '{print $1}' | sed 's/\.//g')
target_verison=490

if [[ $current_version < $target_verison ]];
then
	sudo mkdir kernel-tmp && cd kernel-tmp
	wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.9.40/linux-headers-4.9.40-040940_4.9.40-040940.201707271932_all.deb
	wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.9.40/linux-headers-4.9.40-040940-generic_4.9.40-040940.201707271932_amd64.deb
	wget http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.9.40/linux-image-4.9.40-040940-generic_4.9.40-040940.201707271932_amd64.deb
	sudo dpkg -i *.deb
fi

modprobe tcp_bbr

echo -e "tcp_bbr" >> /etc/modules-load.d/modules.conf
echo -e "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo -e "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf

sysctl -p
sysctl net.ipv4.tcp_available_congestion_control
sysctl net.ipv4.tcp_congestion_control

lsmod | grep bbr
