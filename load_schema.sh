#!/bin/bash
# Loads default schema and promotes database

cat schema.sql | heroku pg:psql
DATABASE=$(heroku pg:info | grep -oE 'HEROKU_POSTGRESQL_[A-Z]+_URL' | head -n 1)
heroku pg:promote $DATABASE
