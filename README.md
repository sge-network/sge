# SGE Network

The Sports, Gaming & Entertainment Network (SGE Network), is a blockchain
designed to support the future of sports betting & related gaming by
leveraging the modular Cosmos design. We believe the future will be heavily shaped by many of the values driving the recent wave of crypto and blockchain development: transparency, increased decentralization, and utility that benefits all stakeholders, especially the user-base.

Utilizing a sovereign blockchain uniquely enables:

- An adaptable framework to design custom applications.
- Enablement of features, tools and economic models where users can directly benefit from the value they help create.
- An unparalleled level of transparency.
- An efficiency of settlement and immediate payout to participants.

At launch, the SGE Network will be optimized to deploy an inaugural application: Six Sigma Sports, which is re-imagining the sports betting landscape and bringing a unique user experience with the benefit of blockchain technology.[Please visit to learn more about Six Sigma Sports.](https://sixsigmasports.io/)

---

## Hardware Requirements

- **Minimal**
  - 1 GB RAM
  - 25 GB SSD
  - 1.4 GHz CPU
- **Recommended**
  - 2 GB RAM
  - 100 GB SSD
  - 2.0 GHz x2 CPU

## Operating System

- Linux/Windows/MacOS(x86)
- **Recommended**
  - Linux(x86_64)

## Installation Steps
>
>Prerequisite: go1.23+ required. [ref](https://golang.org/doc/install)

Sge could be installed by two ways - downloading binary from releases page or build from source.

### Download from releases page

- Download from release required binary

- Check sha256 hash sum

- Place sged into /usr/local/sbin

```shell
sudo mv sged /usr/local/sbin/sged
```

### Building from source
>
>Optional requirement: git. [ref](https://github.com/git/git) and GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)

- Clone git repository

```shell
git clone https://github.com/sge-network/sge.git
```

- Checkout release tag

```shell
cd sge
git fetch --tags
git checkout [vX.X.X]
```

- Install

```shell
go mod tidy
make install
```

### Install system.d service file

```shell
nano /etc/systemd/system/sged.service
```

Please following contents(working dir may be changed as needed)

```systemd
[Unit]
Description=Sge Network node
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/usr/local/sbin/sged start
Restart=on-failure
RestartSec=10
LimitNOFILE=40960

[Install]
WantedBy=multi-user.target
```

Reload unit files in systemd

```shell
sudo systemctl daemon-reload
```

### Generate keys

`sged keys add [key_name]`

or

`sged keys add [key_name] --recover` to regenerate keys with your [BIP39](https://github.com/bitcoin/bips/tree/master/bip-0039) mnemonic

### Connect to a chain and start node

- [Install](#installation-steps) sge application
- Initialize node

```shell
sged init {{NODE_NAME}} --chain-id sgenet-1
```

Select network to join

- Replace `${HOME}/.sge/config/genesis.json` with the genesis file of the chain.
- Add `persistent_peers` or `seeds` in `${HOME}/.sge/config/config.toml`
- Start node

```shell
sged start
```

## Network Compatibility Matrix

| Version | Mainnet | Testnet | SDK Version |
|:-------:|:-------:|:-------:|:-----------:|
|  v1.7.0 |    ✓    |    ✓    |   v0.47.10  |

## Active Networks

### Mainnet

- [sgenet-1](https://github.com/sge-network/networks/tree/master/mainnet/sgenet-1)

- Place the genesis file with the genesis file of the chain.

```shell
wget https://github.com/sge-network/networks/blob/master/mainnet/sgenet-1/genesis.json -O ~/.sge/config/genesis.json
```

Verify genesis hash sum

```shell
sha256sum ~/.sge/config/genesis.json
```

Correct sha256 sum for sgenet-1 genesis file is 3beb0662ade1ad80d41d992bb196770d53a939863c1fed12fa01411dfb981e0b

- Add `persistent_peers` or `seeds` in `${HOME}/.sge/config/config.toml`

```shell
sed -i 's\persistent_peers = ""\persistent_peers = "55f83e1872c482caa102f54e3a73da6c6a146a3f@190.124.251.30:26656,8cb8fecf6470ceaba3f2e7b7c3442b19bd692dea@34.168.149.213:26656,be9721fb11f2ace5b59d26710b4a0d5467ddc8c9@136.243.67.44:17756,d09a5df7a13c758928ab1de0dc7342cab2e7b686@74.50.74.98:36656,401a4986e78fe74dd7ead9363463ba4c704d8759@38.146.3.183:17756,6aa15d14b1e7dadb1923e5701b22c6e370612c29@136.243.67.189:17756,033d3698baf8488429cf2af86ce7d7ad81780a39@[2001:bc8:702:1841::226]:26656,6e0bfbf0c69e60158b310783d129141f88a3c228@5.181.190.81:26656,af9d9bd15ca597eb77dab73c56b0ae51bafcbb28@142.132.202.86:16656,88f341a9670494c3d529934dc578eec1b00f4aa1@141.94.168.85:26656,a44284e563c31676f1c06ff08315d9642e0a6f59@103.230.87.171:26656,17da9d2fea9d6d431d390c3b9575547d8881da2b@185.16.39.190:11156"\g' $HOME/.sge/config/config.toml
```

- Start node

```shell
sged start --minimum-gas-prices [desired-gas-price(ex. 0.001usge)]
```

### Testnet

- [sge-network-4](https://github.com/sge-network/networks/tree/master/testnet/sge-network-4)

- Place the genesis file  with the genesis file of the chain.

```shell
wget https://github.com/sge-network/networks/blob/master/testnet/sge-network-4/genesis.json -O ~/.sge/config/genesis.json
```

Verify genesis hash sum

```shell
sha256sum ~/.sge/config/genesis.json
```

Correct sha256 sum for sge-network-4 genesis file is caa7f15bab24a87718bff96ffeee058373154f7701a1e8977fff46d2f620dbcb

- Add `persistent_peers` or `seeds` in `${HOME}/.sge/config/config.toml`

```shell
sed -i 's\persistent_peers = ""\persistent_peers = "51e4e7b04d2f669f5efa53e8d95891fa04e4c5b9@206.125.33.62:26656,59724f5c6232b1d10507e08b9a9f2ff14181a779@51.195.61.9:20656,7f06552a64b0eed2c4ebd15003a360dbb752e9ce@50.19.180.153:26656,1ae72dbbd1e0143cf2a69441e45eec6dc9212410@52.44.14.245:26656,1e5f1fa5725ab5e09209b7935c6ea3f57b2711ed@[2a01:4f9:1a:9462::3]:26656,13408a5d533afc428a235aa7f58915302c3fccb6@185.246.86.199:26656,7bd23b2967a99b19800282c34b5f509ada38c9ab@52.44.14.245:26656,a37dfffae53ba7a80ef1a54c6906c2072985a3ee@65.108.2.41:56656,476a6214e6abbf038f1e489a3062d62e243150b3@147.135.105.3:26656,1d8dd9667f7a5e83370603fc635a0f0ed7a360d1@50.19.180.153:26656,94f40d2af393be3751518e15818c445632a712a4@84.46.246.109:26656,f5a8e867ae61da981adfb2e142555064694ef541@57.128.37.47:26656,3819c7aebf9ec5f3694747ea3c061b91f555c590@148.251.177.108:17756,58556b5fb572e20d41ce686149ab7b1646ad63a9@65.108.15.170:26656,02ed7e4128bf0bc72a69696aa9157234e0f1e39e@38.146.3.184:17256,e6ad3d00958fafd19f15fa3f151dac8dd8d48c80@5.42.76.30:26656"\g' $HOME/.sge/config/config.toml
```

- Start node

```shell
sged start
```

### Initialize a new chain and start node

- Initialize: `sged init [node_name] --chain-id [chain_name]`
- Add key for genesis account `sged keys add [genesis_key_name]`
- Add genesis account `sged add-genesis-account [genesis_key_name] 1000000000usge`
- Create a validator at genesis `sged gentx [genesis_key_name] 500000000usge --chain-id [chain_name]`
- Collect genesis transactions `sged collect-gentxs`
- Start node `sged start --minimum-gas-prices [desired-gas-price(ex. 0.001usge)]`

### Reset chain

```shell
rm -rf ~/.sge
```

### Shutdown node

```shell
killall sged
```

### Check version

```shell
sged version
```

### Documentations

For the most up-to-date documentation please visit [Gitbook](https://sgenetwork.gitbook.io/documentation-1/)
