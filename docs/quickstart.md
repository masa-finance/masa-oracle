---
id: quickstart
title: Quickstart
---

Follow these steps to get your Masa Oracle node up and running quickly:

### 1. Clone the repository

```bash
git clone https://github.com/masa-finance/masa-oracle.git
cd masa-oracle
```

### 2. Build the node

```bash
go build -v -o masa-node ./cmd/masa-node
```

### 3. Install contract dependencies

```bash
cd contracts/
yarn install
cd ../
```

### 4. Set up environment variables

Create a `.env` file in the root directory with the following content:

```plaintext
# Default .env configuration
BOOTNODES=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa

API_KEY=
RPC_URL=https://ethereum-sepolia.publicnode.com
ENV=dev
FILE_PATH=.
WRITER_NODE=false
CACHE_PATH=CACHE
PORT=8080

# AI LLM
CLAUDE_API_KEY=
CLAUDE_API_URL=https://api.anthropic.com/v1/messages
CLAUDE_API_VERSION=2023-06-01
ELAB_URL=https://api.elevenlabs.io/v1/text-to-speech/ErXwobaYiN019PkySvjV/stream
ELAB_KEY=
OPENAI_API_KEY=
PROMPT="You are a helpful assistant."

# X
TWITTER_USER="yourusername"
TWITTER_PASS="yourpassword"
TWITTER_2FA_CODE="your2fa"

# PG
PG_URL=
```

### 5. Start the node

```bash
./masa-node
```
Your Masa Oracle node should now be running and attempting to connect to the network. Check the logs to ensure it's functioning correctly. You will need your Public Key from the node startup logs to stake the node. Grab some testnet MASA from [Discord](https://discord.gg/masafinance).

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
Is Staked:              false
Is Writer:              false
Is TwitterScraper:      false
Is WebScraper:          false
```

### 6. Stake the node
Grab your Public Key and get some Sepolia MASA from Discord. Then use the following command to initiate staking. Make sure you restart your node once you have staked:

   ```bash
   ./masa-node --stake <amount>
   ```

   Replace `<amount>` with the number of tokens you want to stake. For example, to stake 1000 MASA tokens:
  
   ```bash
   ./masa-node --stake 1000
   ```


