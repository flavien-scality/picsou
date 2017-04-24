#!/bin/bash

apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    git \
    software-properties-common \
    vim
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
apt-key fingerprint 0EBFCD88
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update && apt-get -y install docker-ce
