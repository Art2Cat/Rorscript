#!/usr/bin/bash

wget https://nginx.org/keys/nginx_signing.key
apt add nginx_signing.key

apt install apt-trasport-https -y

CODENAME=`awk -F"[)(]+" '/VERSION=/ {print }' /etc/os-release`

echo "deb https://nginx.org/packages/mainline/debian/ $CODENAME nginx" >> /etc/apt/source.list 
echo "deb-src https://nginx.org/packages/mainline/debian/ $CODENAME nginx" >> /etc/apt/source.list 

apt update 
apt install nginx

ngnix -v
