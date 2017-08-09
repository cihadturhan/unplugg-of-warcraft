#!/usr/bin/env bash

docker stop mongo
docker run --rm --name mongo -d -p "27017:27017" mongo:latest
