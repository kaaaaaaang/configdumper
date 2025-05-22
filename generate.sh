#!/bin/bash


mkdir -p ./tmp
mkdir -p ./versions

vers=("7.5" "8.1" "8.5")
for ver in "${vers[@]}"; do
  echo "Generating configs of version $ver"

  # pd repo
  mkdir -p ./tmp/${ver} && mkdir -p ./versions/v${ver}
  git clone --depth 1 -b release-$ver https://github.com/tikv/pd.git ./tmp/${ver}/pd
  cp ./pd.go.example ./tmp/${ver}/pd/main.go
  cd ./tmp/${ver}/pd && go mod tidy && go run main.go
  cd ../../..
  cp ./tmp/${ver}/pd/pd.json ./versions/v$ver/pd.json

  # tidb repo
  mkdir -p ./tmp/${ver} && mkdir -p ./versions/v${ver}
  git clone --depth 1 -b release-$ver https://github.com/pingcap/tidb.git ./tmp/${ver}/tidb
  cp ./tidb.go.example ./tmp/${ver}/tidb/main.go
  cd ./tmp/${ver}/tidb && go mod tidy && go run main.go
  cd ../../..
  cp ./tmp/${ver}/tidb/tidb.json ./versions/v$ver/tidb.json
  cp ./tmp/${ver}/tidb/tidb-lightning.json ./versions/v$ver/tidb-lightning.json
done
