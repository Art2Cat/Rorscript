#!/bin/bash

sudo apt update

sudo apt install build-essential cmake -y

sudo apt install python-dev python3-dev -y

cd ~/.vim/plugged

git clone https://github.com/Valloric/YouCompleteMe.git

cd ~/.vim/plugged/YouCompleteMe

git submodule update --init --recursive

./install.py --clang-completer --rust-completer
#--go-completer
