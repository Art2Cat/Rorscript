#!/bin/bash 

sudo apt remove --purge vim vim-runtime vim-gnome vim-tiny vim-common vim-gui-common

sudo apt install liblua5.1-dev luajit libluajit-5.1 python-dev \
	       python3-dev libperl-dev libncurses5-dev ruby-dev \
	       libx11-dev dbus-x11 libgnome2-dev libgnomeui-dev \
				 libgtk2.0-dev libatk1.0-dev libbonoboui2-dev \
				 libcairo2-dev libxpm-dev libxt-dev -y


git clone https://github.com/vim/vim.git ~/vim/

cd ~/vim/
echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
echo "configure vim"
echo
./configure --with-features=huge \
			--enable-largefile \
			--enable-rubyinterp=yes \
			--enable-perlinterp=yes \
			--enable-luainterp=yes \
			--with-luajit \
			--enable-pythoninterp \
			--enable-python3interp=yes \
			--enable-gui=auto \
			--enable-cscope \
			--prefix=/usr/local \
			--with-x \

echo
echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
sudo make
echo
echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
sudo make install

sudo update-alternatives --install /usr/bin/editor editor /usr/local/bin/vim 1
sudo update-alternatives --set editor /usr/local/bin/vim
sudo update-alternatives --install /usr/bin/vi vi /usr/local/bin/vim 1
sudo update-alternatives --set vi /usr/local/bin/vim
