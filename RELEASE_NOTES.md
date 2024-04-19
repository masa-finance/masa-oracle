# Masa Oracle Release Notes

## [0.0.11-alpha](https://github.com/masa-finance/masa-oracle/releases) (2024)

> Masa Oracle Node Release

### Breaking Changes

* None

### New Features

* Added Bytes Scraped to NodeData
* Protobuf for node<->node communications

### Bug Fixes

* Moved /status to top level vs api/v1
* Swagger bugs

### Performance Improvements

* Protobuffers for message and scraper worker communication

### ChangeLog

* Flags added
  * --twitterScraper
  * --webScraper
* Upgraded /contracts to @latest
