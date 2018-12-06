#!/bin/bash


YCM=~/.vim/plugged/YouCompleteMe
cd ~/.vim/plugged

if [ ! -d $YCM ]; then
    sudo apt update
    sudo apt install build-essential cmake -y
    sudo apt install python-dev python3-dev -y
    git clone https://github.com/Valloric/YouCompleteMe.git $YCM
fi

cd $YCM

git submodule update --init --recursive

python3 ./install.py --clang-completer --go-completer
