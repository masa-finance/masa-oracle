# Masa Oracle Release Notes

## [0.0.1-beta](https://github.com/masa-finance/masa-oracle/releases) (2024)

> Masa Oracle Node Release

### Breaking Changes

* Gossip version change to 0.0.1-beta

### New Features

* Added Llama3 LLM Model

### Bug Fixes

* Replaces Actor Worker Model with protoactor-go
* Handle RemoteUnreachableEvent to workers 
* Updated API endpoints
  * /api/v1/node/data
  * /api/v1/node/:peerid
  * /api/v1/node/status

### Performance Improvements

* Added node data status to pubsub
* Added worker status to pubsub

### ChangeLog

* Version update
