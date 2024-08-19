---
id: quickstart
title: Quickstart Guide
---

This guide will help you set up and run a Masa Oracle node quickly. 

### 1. Clone and prepare the repository

```bash
git clone https://github.com/masa-finance/masa-oracle.git
```

### 2. Navigate to the project directory
```bash
cd masa-oracle
```

### 3. Build the node

```bash
make build
```

### 4. Install contract dependencies
Navigate to the contract directory:
```bash
cd contracts
```

Install dependencies using yarn
```bash
yarn install
```

Return to the root directory
```bash
cd ..
```

### 5. Set up basic environment variables

Create a `.env` file in the root directory with these essential variables:
```plaintext
# Default .env configuration
BOOTNODES=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa,/ip4/34.121.111.128/udp/4001/quic-v1/p2p/16Uiu2HAmKULCxKgiQn1EcfKnq1Qam6psYLDTM99XsZFhr57wLadF

API_KEY=
RPC_URL=https://ethereum-sepolia.publicnode.com
ENV=test
FILE_PATH=.
VALIDATOR=false
PORT=8080

# AI LLM
CLAUDE_API_KEY=
CLAUDE_API_URL=https://api.anthropic.com/v1/messages
CLAUDE_API_VERSION=2023-06-01
ELAB_URL=https://api.elevenlabs.io/v1/text-to-speech/ErXwobaYiN019PkySvjV/stream
ELAB_KEY=
OPENAI_API_KEY=
PROMPT="You are a helpful assistant."

# Twitter(X)
TWITTER_SCRAPER=false
TWITTER_USERNAME="yourusername"
TWITTER_PASSWORD="yourpassword"
TWITTER_2FA_CODE="your2fa"

# Discord
DISCORD_SCRAPER=false
DISCORD_BOT_TOKEN=

# Telegram
TELEGRAM_SCRAPER=false
TELEGRAM_APP_ID=
TELEGRAM_APP_HASH=

#Web
WEB_SCRAPER=false
```
Note: Full configuration (such as API keys, Twitter, Discord and Telegram credentials, and worker node settings) are not required for the initial node setup. We will add those later once we get a node running and staked.

### 6. Start the node to obtain your public key

```bash
make run
```

Look for your `Public Key` in the startup logs. You'll need this to stake and participate in the network. Example:

```bash
#######################################
#     __  __    _    ____    _        #
#    |  \/  |  / \  / ___|  / \       #
#    | |\/| | / _ \ \___ \ / _ \      #
#    | |  | |/ ___ \ ___) / ___ \     #
#    |_|  |_/_/   \_\____/_/   \_\    #
#                                     #
#######################################

Version:             v0.5.0
Multiaddress:        /ip4/192.168.1.8/udp/4001/quic-v1/p2p/16Uiu2HAmDXWNV9RXVoRsbt9z7pFSsKS2KdpN7HHFVLdFZmS7iCvo
IP Address:          /ip4/127.0.0.1/udp/4001/quic-v1
Public Key:          0x5dA36a3eB07fd1624B054b99D6417DdF2904e826
Is Staked:           false
Is Validator:        false
Is TwitterScraper:   false
Is DiscordScraper:   false
Is TelegramScraper:  false
```
Once you've noted your Public Key, you can stop the node (Ctrl+C).


### 7. Prepare for staking
To participate in the network, you need:

1. Sepolia ETH (about 0.015 ETH) to pay for transaction fees
2. 1,000 Sepolia MASA token to stake

#### a) Get Sepolia ETH:

- Use a Sepolia ETH testnet faucet 
- Send approximately 0.015 Sepolia ETH to your public key address

#### b) Get Sepolia MASA tokens:
Once you have Sepolia ETH, run: 

```bash
make faucet
```

### 8. Stake your Sepolia MASA
Stake your MASA tokens:

```bash
make stake
```
You'll see approval and staking transaction hashes in the output.



### 9. Verify staking and start the node again:
```bash
make run
```
Verify that the `Is Staked` flag has changed to `true` in the startup logs:
```bash
#######################################
#     __  __    _    ____    _        #
#    |  \/  |  / \  / ___|  / \       #
#    | |\/| | / _ \ \___ \ / _ \      #
#    | |  | |/ ___ \ ___) / ___ \     #
#    |_|  |_/_/   \_\____/_/   \_\    #
#                                     #
#######################################
Multiaddress:           /ip4/192.168.1.25/udp/4001/quic-v1/p2p/16Uiu2HAm28dTN2WVWD2y2bjzwPdym59XASDfQsSktCtejtNR9Vox
IP Address:             /ip4/127.0.0.1/udp/4001/quic-v1
Public Key:             0x065728510468A2ef48e6E8a860ff42D68Ca612ee
Is Staked:              true
Is Validator:              false
Is TwitterScraper:      false
Is WebScraper:          false
INFO[0001] Peer added to DHT: 16Uiu2HAmHpx13GPKZAP3WpgpYkZ39M5cwuvmXS5gGvrsa5ofLNoq 
INFO[0005] Successfully advertised protocol /masa/oracle_protocol/v0.0.9-beta-dev
```


### 10. Access the Swagger API
Open your browser and navigate to:

```bash
http://localhost:8080/swagger/index.html
```
Congratulations! You've completed the basic setup of your Masa Oracle node. 