#!/bin/bash
#
# source: https://toutyrater.github.io/advanced/tls.html
#
# https://www.github.com/Neilpan/acme.sh
#

if [ -z $1 ]; then
	echo "please enter your domain"
	echo "for issue: ./install_tls.sh yourdomain.com i "
	echo "for renew: ./install_tls.sh yourdomain.com r "
	exit
fi

DOMAIN=$1
PKGM=''
#if grep 'Ubuntu' /etc/*-release
#then
#	PKGM='apt'
#elif grep 'Arch' /etc/*-release
#then
#	PKGM='pacman -S'
#fi

if [ -z "$(which curl)" ]; then
	sudo apt update
	sudo apt -y install curl
	sudo apt -y install netcat socat
fi

DIR=~/.acme.sh
if [ ! -d $DIR ]; then
	curl https://get.acme.sh | sh
fi

if [ $2 == "i" ]; then
	~/.acme.sh/acme.sh --issue -d $DOMAIN --standalone -k ec-256
elif [ $2 == "r" ]; then
	~/.acme.sh/acme.sh --renew -d $DOMAIN --force --ecc
fi

~/.acme.sh/acme.sh --installcert -d $DOMAIN --fullchainpath /etc/v2ray/v2ray.crt --keypath /etc/v2ray/v2ray.key --ecc

