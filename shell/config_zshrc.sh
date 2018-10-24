#!/bin/bash

# if [ $UID -ne 0 ];
# then
#		echo "Not root user. Please run as root."
#		exit
# else
#		echo "Root user."
# fi
if [ -z "$(which zsh)" ]; then
    if [ "$(uname)" == "Darwin" ]; then
        brew install zsh
        echo "bash -c zsh" >> ~/.bash_profile
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        sudo apt update
        sudo apt install zsh -y
        chsh -s `which zsh`
        echo "bash -c zsh" >> ~/.bashrc
    fi
fi

sudo sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"

echo "Install zsh-syntax-hightlighting"
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

cp ../dotfile/.zshrc ~/

source ~/.bashrc
source ~/.zshrc
