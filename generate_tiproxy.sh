#!/bin/bash

mkdir -p ./tmp
mkdir -p ./versions

vers=("1.1" "1.2" "1.3")
for ver in "${vers[@]}"; do
  echo "Generating configs of version $ver"

  # tiproxy repo
  mkdir -p ./tmp/${ver} && mkdir -p ./versions/v${ver}
  git clone --depth 1 -b v$ver.0 https://github.com/pingcap/tiproxy.git ./tmp/${ver}/tiproxy
  cp ./tiproxy.go.example ./tmp/${ver}/tiproxy/main.go
  cd ./tmp/${ver}/tiproxy && go mod tidy && go run main.go
  cd ../../..
  cp ./tmp/${ver}/tiproxy/tiproxy.json ./versions/v$ver/tiproxy.json
done

