# Masa Oracle Release Notes

[All Releases](https://github.com/masa-finance/masa-oracle/releases)

## [0.0.7-beta](https://github.com/masa-finance/masa-oracle/releases) (2024)

## Overview

This release of the Masa Oracle Node introduces new features, masa node will exit after running cli commands --stake n and --faucet

### Breaking Changes

* None

### Bug fixes

### New Features

* Added test Discord Bot Token to example.env for beta use
* Added masa token faucet with cli param --faucet

> compile and build the masa-node

```shell
make build
```

> will give 1000 masa tokens to the node and exit

```shell
make faucet 
```

> will stake 1000 masa tokens and exit

```shell
make stake 
```

> run as normal

```shell
make run
```

## Change Log

* Version: 0.0.7-beta
