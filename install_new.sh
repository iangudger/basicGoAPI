#!/bin/bash
# Download and run this script to make a new deployment on a clean Ubuntu install.
# REQUIRED: Verified Heroku account.

# Get source
mkdir -p ~/go/src/github.com/iangudger/basicGoAPI
cd ~/go/src/github.com/iangudger/basicGoAPI
sudo apt-get update
sudo apt-get install -y git
git clone https://github.com/iangudger/basicGoAPI.git .

# Setup Git
GITNAME=$(git config --get user.name)
while [ -z "$GITNAME" ]
	do
		echo -n "Git name: "
		read GITNAME
		sleep 1;
done
git config --global user.name "$GITNAME"

GITEMAIL=$(git config --get user.email)
while [ -z "$GITEMAIL" ]
	do
		echo -n "Git email: "
		read GITEMAIL
		sleep 1;
done
git config --global user.email "$GITEMAIL"

# Install dependencies
./setupenv.sh

# Create new Heroku app with required addons
heroku create -b https://github.com/kr/heroku-buildpack-go.git
heroku addons:add heroku-postgresql
heroku addons:add mandrill

# Setup default database
./load_schema.sh

# Go variables
if [ -z "$GOPATH" ]
	then
		export GOPATH=$HOME/go
		export GOBIN=$GOPATH/bin
		export PATH=$PATH:$GOBIN
fi

# Deploy app
godep save
git add -A .
git commit -m "Added dependencies."
git push -u heroku master
