#!/usr/bin/env bash

set -eo pipefail

# get protoc executions
go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v1.3.3-alpha.regen.1 2>/dev/null

# get cosmos sdk from github
go get github.com/cosmos/cosmos-sdk@v0.45.11 2>/dev/null

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./sge -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &>/dev/null; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# move proto files to the right places
#
# Note: Proto files are suffixed with the current binary version.
cp -r github.com/sge-network/sge/* ./
rm -rf github.com

go mod tidy -compat=1.18
