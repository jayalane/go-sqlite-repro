#! /bin/bash

set -x
set -e

docker build -t chlane/demo  -f ./Dockerfile .
docker run chlane/demo:latest
