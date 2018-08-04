#!/bin/bash
# install and config git on ubuntu.
sudo apt update

echo ""
if [[ -z $(which git) ]]; then
	sudo apt install -y git
fi

echo ""
read -r -t 60 -p "do you want to set git config[Y/n]:" response
response=${response,,}
if [[ $response == "y" ]]; then
	echo "all configuration will be set as global"
	read -r -t 60 -p "enter your user name: " username
	if [[ -n $username ]]; then
		git config --global user.name $username
	fi
	read -r -t 60 -p "enter your email: " email
	if [[ -n $email ]]; then
		git config --global user.email $email
		echo ""
	fi
	read -r -t 60 -p "enter your core editor: " editor
	if [[ -n $email ]]; then
		git config --global core.editor $editor
		echo ""
	fi
	read -r -t 60 -p "enter your http proxy host: " host
	echo ""
	read -r -t 60 -p "enter your http proxy port: " port
	echo ""
	read -r -t 60 -p "enter your sock5 proxy port: " sock5port
	if [[ -n $host && -n $port ]]; then
		git config --global http.https://github.com.proxy  http://$host:$port
		git config --global https.https://github.com.proxy https://$host:$port
	fi
	if [[ -n $host && -n $sock5port ]]; then
		git config --global http.https://github.com.proxy  http://$host:$sock5port
		git config --global https.https://github.com.proxy https://$host:$sock5port
	fi
	git config --global core.autocrlf input
	git config --global http.postBuffer 524288000
	echo ""
	echo "######################################################"
	git config -l
fi

echo "######################################################"
version=`git --version`
echo "installed $version"
