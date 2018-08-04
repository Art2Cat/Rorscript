#!/bin/bash

git clone https://github.com/universal-ctags/ctags.git ~/ctags
cd ~/ctags
if [ -z $(which make) ]
then
	sudo apt -y install make
fi

if [ -z `which autoreconf` ]
then
	sudo apt -y install autoreconf
	sudo apt -y install dh-autoreconf
fi

if [ -z `which pkg-config` ]
then
	sudo apt -y install pkg-config
fi

./autogen.sh
./configure
make
sudo make install
