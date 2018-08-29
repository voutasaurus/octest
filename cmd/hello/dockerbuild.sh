#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"
cd $scriptdir

# main
GOOS=linux go build -o ochello
docker build -t ochello .

# cleanup
rm ochello
set +ex
