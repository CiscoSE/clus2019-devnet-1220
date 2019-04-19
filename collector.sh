#!/bin/sh

sudo tar -C /usr/local -xzf /vagrant/go1.12.3.linux-amd64.tar.gz

cat << EOF >> /home/vagrant/.bashrc 
cd /vagrant
source .env
cd /vagrant/src/github.com/CiscoLive/telemetry/dial_in
sshpass -p vagrant ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null vagrant@192.0.2.2 'bash cat /misc/config/grpc/ems.pem' | grep -Pzo '\-----BE(.*\n)*' | sed $'s/[^[:print:]\t]//g' > /vagrant/src/github.com/CiscoLive/telemetry/dial_in/ems.pem
EOF

cd /vagrant
export GOPATH=$PWD
export GOBIN=$PWD/bin
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
mkdir pkg
mkdir bin
go get github.com/golang/protobuf/proto
go get github.com/nleiva/xrgrpc
sudo apt-get update 
sudo apt-get install sshpass -y
