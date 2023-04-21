#!/bin/bash -ex
# Sorry, I don't have any steps for installing
# postgresql@15 at this time. Please do that 
# yourself.

mkdir -p pgdata

# This is not safe for production.
initdb --pgdata=pgdata --username=admin
