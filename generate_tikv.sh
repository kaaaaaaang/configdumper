#!/bin/bash

ver=$1
full_ver=$2
echo "$1 $2"
mkdir -p ./tmp
mkdir -p ./versions

echo "Generating tikv configs of version $ver"
tiup tikv:$full_ver --config-info json | go run main.go > ./versions/$ver/tikv.json
