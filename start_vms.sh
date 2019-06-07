#!/bin/bash
docker rm -f devnet_3000
docker run -d -p 12345:8080 --name devnet_3000 sfloresk/devnet_3000
cp /Volumes/CLUS2019/go1.12.3.linux-amd64.tar.gz $HOME/devnet_3000/
vagrant box add --name ubuntu/xenial64.sfloresk /Volumes/CLUS2019/xenial-server-cloudimg-amd64-vagrant.box
vagrant box add --name iosxrv-fullk9-x64.snapshot.6.4.2.sfloresk /Volumes/CLUS2019/iosxrv-fullk9-x64.snapshot.6.4.2.box
vagrant up
