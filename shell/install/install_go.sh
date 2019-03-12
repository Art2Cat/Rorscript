#!/bin/bash

if [ -z $1 ]; then
	echo "please enter the go version"
	echo "example: ./install_go.sh 1.11.1 "
	exit
fi

VERSION=$1
if [ -z "$(which wget)" ]; then
    sudo apt-get update
	sudo apt install -y wget
fi

ARCH=$(uname -m)
if [ ${ARCH} == 'x86_64' ]; then
	ARCH=amd64
elif [ ${ARCH} == 'x86' ]; then
	ARCH=386
elif [ ${ARCH} == 'ARMv6' ]; then
	ARCH=armv6l
fi

OS=$(echo $(uname -s) | awk '{print tolower($0)}')
go_archive_file=go$VERSION.$OS-$ARCH.tar.gz
echo
echo "Downloading $go_archive_file"

wget https://dl.google.com/go/$go_archive_file

if [ ! -f ./$go_archive_file ]; then
	echo "go binary file did not exists."
	exit
fi

if [ -d /usr/local/go ]; then
	sudo rm -rf /usr/local/go
fi

sudo tar -C /usr/local -xvf $go_archive_file

if [[ -z $(cat ~/.zshrc | grep GOROOT | head -n1) ]]; then
	echo "Did not found GOROOT in the .zshrc. configuring..."
	echo -e "export GOROOT=/usr/local/go" >> ~/.zshrc
	echo -e "export PATH=\$PATH:\$GOROOT/bin" >> ~/.zshrc
	echo -e "export GOPROXY=https://goproxy.io" >> ~/.zshrc
	source ~/.zshrc
fi

echo
go version

rm -rf $go_archive_file
