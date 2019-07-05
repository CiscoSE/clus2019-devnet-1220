# Cisco Live 2019 - Devnet Workshop 3000

This repo contains source code and instructions for the DevNet workshop 3000 - Real Time Telemetry with go delivered at Cisco Live 2019.

# Lab resources

Please contact sfloresk@cisco.com if you need access to the vagrant images and the guide to perform the lab.

# Running locally 

If you want to run the app in your laptop instead of in the Ubuntu VM, you will need to go over these steps: (assuming either Mac or Linux OS)

1- Install Go (Tested with go1.12.3) - https://golang.org/doc/install?download=go1.12.4.linux-amd64.tar.gz 

2- Clone the repo

```bash
git clone https://github.com/CiscoSE/devnet_3000.git $HOME/devnet_3000
```

3- Source enviroment and download libraries

```bash
cd $HOME/devnet_3000
source .env
mkdir pkg
mkdir bin
go get github.com/golang/protobuf/proto
go get github.com/nleiva/xrgrpc
```

4- Build and run

```bash
cd src/github.com/CiscoLive/telemetry/dial_in
go build
/dial_in
```

