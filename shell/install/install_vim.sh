#!/bin/bash
# install vim8 on ubuntu.

sudo apt install software-properties-common

sudo add-apt-repository ppa:jonathonf/vim

sudo apt update

sudo apt install libncurses5-dev libgnome2-dev libgnomeui-dev \
    libgtk2.0-dev libatk1.0-dev libbonoboui2-dev \
    libcairo2-dev libx11-dev libxpm-dev libxt-dev python-dev \
	python3-dev ruby-dev mercurial

sudo apt install vim -y

