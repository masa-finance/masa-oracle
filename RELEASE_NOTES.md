# Masa Oracle Release Notes

## [0.0.3-beta](https://github.com/masa-finance/masa-oracle/releases) (2024)

> Masa Oracle Node Release

### Breaking Changes

* Normalized all command line params to camelCase
* Changed default model to ollama/*

### New Features

* Added Cloudflare AI Workers for LLM compute
* Implemented Record and OracleData struct to save events to persisted storage
* Discord Scraper as pkg
* Twitter Scraper as pkg
* Reddit Scraper as pkg
* Web Scraper as pkg

### Bug Fixes

* Worker channel race condition

### Performance Improvements

* Increased worker time syncing between peers

### ChangeLog

* Updated swagger docs
* version 0.0.3-beta
