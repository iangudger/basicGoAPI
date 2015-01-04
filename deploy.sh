#!/bin/bash
# Deploys a new version of the server.
# Assumes the git remote push repository is Heroku
# Does not make any changes to the database.

git rm -rf Godeps
./format.sh
godep save
git add .
git commit
git push
