# Coreth and the C-Chain

[Avalanche](https://www.avax.network/) is a network composed of multiple blockchains.
Each blockchain is an instance of a Virtual Machine (VM), much like an object in an object-oriented language is an instance of a class.
That is, the VM defines the behavior of the blockchain.
Coreth (from core Ethereum) is the [Virtual Machine (VM)](https://docs.avax.network/learn/virtual-machines) that defines the Contract Chain (C-Chain).
This chain implements the Ethereum Virtual Machine and supports Solidity smart contracts as well as most other Ethereum client functionality.

## Building

Coreth is a dependency of AvalancheGo which is used to implement the EVM based Virtual Machine for the Avalanche C-Chain. In order to run with a local version of Coreth, users must update their Coreth dependency within AvalancheGo to point to their local Coreth directory. If Coreth and AvalancheGo are at the standard location within your GOPATH, this will look like the following:

```bash
cd $GOPATH/src/github.com/ava-labs/avalanchego
go mod edit -replace github.com/ava-labs/coreth=../coreth
```

Now that AvalancheGo depends on the local version of Coreth, we can build with the normal build script:

```bash
./scripts/build.sh
./build/avalanchego
```

Note: the C-Chain originally ran in a separate process from the main AvalancheGo process and communicated with it over a local gRPC connection. When this was the case, AvalancheGo's build script would download Coreth, compile it, and place the binary into the `avalanchego/build/plugins` directory.

### Optional Dev Shell

Some activities, such as collecting metrics and logs from the nodes targeted by an e2e
test run, require binary dependencies. One way of making these dependencies available is
to use a nix shell which will give access to the dependencies expected by the test
tooling:

 - Install [nix](https://nixos.org/). The [determinate systems
   installer](https://github.com/DeterminateSystems/nix-installer?tab=readme-ov-file#install-nix)
   is recommended.
 - Use `./scripts/dev_shell.sh` to start a nix shell
 - Execute the dependency-requiring command (e.g. `./scripts/tests.e2e.sh --start-collectors`)

This repo also defines a `.envrc` file to configure [devenv](https://direnv.net/). With
`devenv` and `nix` installed, a shell at the root of the repo will automatically start a nix
dev shell.

## API

The C-Chain supports the following API namespaces:

- `eth`
- `personal`
- `txpool`
- `debug`

Only the `eth` namespace is enabled by default. 
To enable the other namespaces see the instructions for passing the C-Chain config to AvalancheGo [here.](https://docs.avax.network/nodes/configure/chain-config-flags#enabling-evm-apis)
Full documentation for the C-Chain's API can be found [here.](https://docs.avax.network/reference/avalanchego/c-chain/api)

## Compatibility

The C-Chain is compatible with almost all Ethereum tooling, including [Core,](https://docs.avax.network/build/dapp/launch-dapp#through-core) [Metamask,](https://docs.avax.network/build/dapp/launch-dapp#through-metamask) [Remix](https://docs.avax.network/dapps/smart-contract-dev/deploy-with-remix-ide) and [Truffle.](https://docs.avax.network/build/tutorials/smart-contracts/using-truffle-with-the-avalanche-c-chain)

## Differences Between Avalanche C-Chain and Ethereum

### Atomic Transactions

As a network composed of multiple blockchains, Avalanche uses *atomic transactions* to move assets between chains. Coreth modifies the Ethereum block format by adding an *ExtraData* field, which contains the atomic transactions.

### Block Timing

Blocks are produced asynchronously in Snowman Consensus, so the timing assumptions that apply to Ethereum do not apply to Coreth. To support block production in an async environment, a block is permitted to have the same timestamp as its parent. Since there is no general assumption that a block will be produced every 10 seconds, smart contracts built on Avalanche should use the block timestamp instead of the block number for their timing assumptions.

A block with a timestamp more than 10 seconds in the future will not be considered valid. However, a block with a timestamp more than 10 seconds in the past will still be considered valid as long as its timestamp is greater than or equal to the timestamp of its parent block.

## Difficulty and Random OpCode

Snowman consensus does not use difficulty in any way, so the difficulty of every block is required to be set to 1. This means that the DIFFICULTY opcode should not be used as a source of randomness.

Additionally, with the change from the DIFFICULTY OpCode to the RANDOM OpCode (RANDOM replaces DIFFICULTY directly), there is no planned change to provide a stronger source of randomness. The RANDOM OpCode relies on the Eth2.0 Randomness Beacon, which has no direct parallel within the context of either Coreth or Snowman consensus. Therefore, instead of providing a weaker source of randomness that may be manipulated, the RANDOM OpCode will not be supported. Instead, it will continue the behavior of the DIFFICULTY OpCode of returning the block's difficulty, such that it will always return 1.

## Block Format

To support these changes, there have been a number of changes to the C-Chain block format compared to what exists on Ethereum.

### Block Body

* `Version`: provides version of the `ExtData` in the block. Currently, this field is always 0.
* `ExtData`: extra data field within the block body to store atomic transaction bytes.

### Block Header

* `ExtDataHash`: the hash of the bytes in the `ExtDataHash` field
* `BaseFee`: Added by EIP-1559 to represent the base fee of the block (present in Ethereum as of EIP-1559)
* `ExtDataGasUsed`: amount of gas consumed by the atomic transactions in the block
* `BlockGasCost`: surcharge for producing a block faster than the target rate

## Releasing

See [docs/releasing/README.md](docs/releasing/README.md) for the release process.
