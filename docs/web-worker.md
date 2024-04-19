---
id: web-worker
title: Providing Compute for Web Scraping Requests
---

## Introduction

This guide is designed for oracle node workers who are interested in contributing compute resources to fulfill web scraping data requests within the Masa Oracle Node network. It outlines the specific requirements, setup, and operational guidelines for workers to efficiently process web scraping requests. By configuring your node to handle these requests, you play a pivotal role in the decentralized extraction of data from various websites, which is essential for applications such as content aggregation, market analysis, and competitive research.

## Getting Started: Worker's Role in Processing Web Scraping Data

As a worker in the Masa Oracle Node network, your primary function is to process web scraping data requests sent by clients. This involves extracting data from websites based on specified parameters and returning the data to the network. Here's a brief overview of the workflow:

### Worker's Workflow

1. **Initialization**: Your node, acting as a Worker, joins a pool managed by a Manager actor. This ensures efficient distribution and management of incoming web scraping requests.

2. **Receiving Requests**: When a web scraping request is received, the Manager actor assigns the task to you, the Worker, based on your availability and scraping capabilities.

3. **Processing Requests**: You execute the web scraping script to collect the required data from the target website and format the data for return to the network.

## Prerequisites for Web Scraping Workers

To become a worker focused on web scraping data requests, you need to:

- Have your Masa Oracle Node staked as per the [Staking Guide for Masa Oracle Node](staking-guide.md).
- Install web scraping tools and libraries on your node.
- Ensure your Masa Oracle Node is operational, with network accessibility for receiving and processing requests.

## Setting Up Your Node for Web Scraping Requests

BP TODO
