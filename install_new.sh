#!/bin/bash
# Download and run this script to make a new deployment on a clean Ubuntu install.
# REQUIRED: Verified Heroku account.

# Get source
mkdir -p ~/go/src/github.com/iangudger/basicGoAPI
cd ~/go/src/github.com/iangudger/basicGoAPI
sudo apt-get update
sudo apt-get install -y git
git clone https://github.com/iangudger/basicGoAPI.git .

# Install dependencies
./setupenv.sh

# Create new Heroku app with required addons
heroku create -b https://github.com/kr/heroku-buildpack-go.git
heroku addons:add heroku-postgresql
heroku addons:add mandrill

# Setup default database
./load_schema.sh

# Deploy app
git push -u heroku master
./deploy.sh
