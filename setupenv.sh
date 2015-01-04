#!/bin/bash
# Setups up development environment and installs dependencies

# Heroku
wget -qO- https://toolbelt.heroku.com/install-ubuntu.sh | sh

# Golang
sudo add-apt-repository -y ppa:eugenesan/ppa
sudo apt-get update
GOEXE = $(which go)
if [ -z "$GOEXE" ]
	then
		sudo apt-get install -y golang
	else
		sudo apt-get dist-upgrade
fi

# Postgres
sudo apt-get install -y postgresql-client

# Go variables
if [ -z "$GOPATH" ]
	then
		mkdir -p ~/go/src
		echo  '' >> ~/.bashrc
		echo  '# For programming language Go' >> ~/.bashrc
		echo  'export GOPATH=$HOME/go' >> ~/.bashrc
		echo  'export GOBIN=$GOPATH/bin' >> ~/.bashrc
		echo  'export PATH=$PATH:$GOBIN' >> ~/.bashrc	
		echo  '' >> ~/.bashrc
		
		export GOPATH=$HOME/go
		export GOBIN=$GOPATH/bin
		export PATH=$PATH:$GOBIN
fi

# Required libraries
go get "github.com/keighl/mandrill"
go get "github.com/lib/pq"
go get "golang.org/x/crypto/bcrypt"
go get "github.com/kr/godep"
