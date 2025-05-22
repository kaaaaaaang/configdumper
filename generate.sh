#!/bin/bash


mkdir -p ./tmp
mkdir -p ./versions

vers=("6.5" "7.1" "7.5" "8.1" "8.5")
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
  # old config package path 
  if [[ $ver == "6.1" || $ver == "6.5" || $ver == "7.1" ]]; then
    sed -i 's|github.com/pingcap/tidb/pkg/config|github.com/pingcap/tidb/config|g' ./tmp/${ver}/tidb/main.go
    sed -i 's|github.com/pingcap/tidb/pkg/lightning/config|github.com/pingcap/tidb/br/pkg/lightning/config|g' ./tmp/${ver}/tidb/main.go
  elif [[ $ver == "7.5" ]]; then
    sed -i 's|github.com/pingcap/tidb/pkg/lightning/config|github.com/pingcap/tidb/br/pkg/lightning/config|g' ./tmp/${ver}/tidb/main.go
  fi
  cd ./tmp/${ver}/tidb && go mod tidy && go run main.go
  cd ../../..
  cp ./tmp/${ver}/tidb/tidb.json ./versions/v$ver/tidb.json
  cp ./tmp/${ver}/tidb/tidb-lightning.json ./versions/v$ver/tidb-lightning.json
done
