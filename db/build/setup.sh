#!/bin/bash -ex
# Make the db folder ready to run. Start the postgresql
# server.
#
# Sorry, I don't have any steps for installing
# postgresql@15 at this time. Please do that 
# yourself.

mkdir -p tmp

# This is not safe for production.
initdb --pgdata=tmp/pgdata
pg_ctl --pgdata=tmp/pgdata -l tmp/log start
createdb main

psql -d main -c "\
CREATE TABLE clients (
    name VARCHAR ( 50 ) NOT NULL PRIMARY KEY, 
    ip VARCHAR ( 50 ) NOT NULL
);"

cat tmp/log
