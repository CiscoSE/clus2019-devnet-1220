#!/bin/bash
docker rm -f devnet_3000
docker run -d -p 12345:8080 --name devnet_3000 sfloresk/devnet_3000
curl -O https://dl.google.com/go/go1.12.3.linux-amd64.tar.gz
vagrant up
