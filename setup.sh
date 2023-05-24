rm -rf ~/.sge

go mod tidy

make install

sged init local --chain-id saage

sed -i -e 's/stake/usge/' ~/.sge/config/genesis.json

#sged keys add valnode1 > mnemonic

sged add-genesis-account valnode1 1136900000000000usge
