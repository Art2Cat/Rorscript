#!/bin/bash
# install jdk on ubuntu.

function mkjvmdir() {
  jvm=/user/lib/jvm
  if [ -e $jvm ]; then
    sudo mkdir -p $jvm
  fi
}

function config_environment() {
  java_home=$(cat .bashrc | grep JAVA_HOME | head -n1)
  if [ -z $java_home ]; then
    echo -e "export JAVA_HOME=/usr/lib/jvm/jdk-9.0.4" >> ~/.bashrc
    echo -e "export JRE_HOME=${JAVA_HOME}/jre" >> ~/.bashrc
    echo -e "export CLASSPAHT=.:${JAVA_HOME}/lib:${JRE_HOME}/lib" >> ~/.bashrc
    echo -e "export PATH=${JAVA_HOME}/bin:$PATH" >> ~/.bashrc
  fi
  source ~/.bashrc
}

read -p "openjdk or oracle jdk(1/2)" response
if [[ $response -eq 1 ]]; then
  sudo apt update && sudo apt install default-jdk -y
else
  mkjvmdir
  wget -O 'http://download.oracle.com/otn-pub/java/jdk/9.0.4+11/c2514751926b4512b076cc82f959763f/jdk-9.0.4_linux-x64_bin.tar.gz?AuthParam=1516608341_3343898786f6f1b4fa420f46ce81b023'
  
  sudo tar -zxvf jdk-9.0.4_linux-x64_bin.tar.gz -C /usr/lib/jvm
fi

if [ -z $(which java) ]; then
  config_environment
fi

update-alternatives --config java

java -version
