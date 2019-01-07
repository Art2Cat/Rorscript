#!/bin/bash
# install maven on ubuntu.
wget http://www.eu.apache.org/dist/maven/maven-3/3.5.4/binaries/apache-maven-3.5.4-bin.tar.gz
tar -xvzf apache-maven-3.5.4-bin.tar.gz
sudo mkdir /usr/local/maven
sudo mv apache-maven-3.5.4/ /usr/local/maven/
sudo update-alternatives --install /usr/bin/mvn mvn /usr/local/maven/apache-maven-3.5.4/bin/mvn 1
sudo update-alternatives --config mvn
rm apache-maven-3.5.4-bin.tar.gz
