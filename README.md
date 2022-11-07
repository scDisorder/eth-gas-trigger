# Ethereum Gas Trigger

Simple program to watch ethereum network gas prices

```shell
# set provider url
export ETHEREUM_PROVIDER_URL=https://eth-goerli.g.alchemy.com/v2/SECRET_API_KEY

# build executable binary
go build -o eth-gas-trigger

# Run for one execution 
./eth-gas-trigger run --gwei 10 --cmd 'npx hardhat run scripts/deploy.js --network goerli'

# Run repeatable execution
./eth-gas-trigger --repeatable run --gwei 10 --cmd 'npx hardhat run scripts/whoami.js --network mainnet'
```
