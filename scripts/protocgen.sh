#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto
buf mod update
cd ..
buf generate

# move proto files to the right places
cp -r ./github.com/sge-network/sge/x/* x/
cp -r ./github.com/sge-network/sge/types/* types/

rm -rf ./github.com