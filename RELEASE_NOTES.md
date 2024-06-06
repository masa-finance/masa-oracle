# Masa Oracle Release Notes

## [0.0.6-beta](https://github.com/masa-finance/masa-oracle/releases) (2024)

> Masa Oracle Node Release

### Breaking Changes

* None

### New Features

#### Worker storage metrics #314

* Store CID in DHT for Availability
* Store CID and Underlying Data in LevelDB for Persistence
* Metrics to Track
* Storage of BytesScraped by Peer
* PeerID Selection and Response

### Bug Fixes

* Fixed error response handling from workers

### Performance Improvements

* Added protobuf message.Response for worker responses
* REQUIREMENT: Port 4001 TCP inbound needs to be open

### ChangeLog

* version string 0.0.6-beta
