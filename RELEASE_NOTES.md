# Masa Oracle Release Notes

## [0.0.2-beta](https://github.com/masa-finance/masa-oracle/releases) (2024)

> Masa Oracle Node Release

### Breaking Changes

* None

### New Features

* Allow Nodes to participate in twitter scraping w/o bringing their own creds
* Added get bootnodes from deployment json on s3

### Bug Fixes

* None

### Performance Improvements

* Increased worker time syncing between peers

### ChangeLog
* 
* Removed obsolete pg integration for new data persistence architecture
* Added LLM_TWITTER_PROMPT to .env (optional)
* Added LLM_SCRAPER_PROMPT to .env (optional)
