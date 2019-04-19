#!/bin/sh

sudo tar -C /usr/local -xzf /vagrant/go1.12.3.linux-amd64.tar.gz

cat << EOF >> /home/vagrant/.bashrc 
cd /vagrant
source .env
cd /vagrant/src/github.com/CiscoLive/telemetry/dial_in
EOF


