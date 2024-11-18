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

# Check bootnodes addresses in the Masa documentation https://developers.masa.ai/docs/welcome-to-masa
BOOTNODES=

API_KEY=
RPC_URL=https://ethereum-sepolia.publicnode.com
ENV=test
FILE_PATH=.
VALIDATOR=false
PORT=8080

# X
TWITTER_USERNAME="yourusername"
TWITTER_PASSWORD="yourpassword"
TWITTER_2FA_CODE="your2fa"

# Worker node config; default = false
TWITTER_SCRAPER=false
DISCORD_SCRAPER=false
WEB_SCRAPER=false
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
Is Validator:              false
Is TwitterScraper:      false
Is WebScraper:          false
INFO[0001] Peer added to DHT: 16Uiu2HAmHpx13GPKZAP3WpgpYkZ39M5cwuvmXS5gGvrsa5ofLNoq 
INFO[0005] Successfully advertised protocol /masa/oracle_protocol/v0.0.9-beta-dev
```

### 6. Stake the node with 1000 Sepolia MASA minimum

Grab your Public Key and get some Sepolia MASA from Discord. Then use the following command to initiate staking. Make sure you restart your node once you have staked:

   ```bash
   ./masa-node --stake <amount>
   ```

   Replace `<amount>` with the number of tokens you want to stake. For example, to stake 1000 MASA tokens:
  
   ```bash
   ./masa-node --stake 1000
   ```

Your tokens will approve and stake:

```bash
Approving staking contract to spend tokens.....
0x8de79f5111b185fe67090f904b72f3dda7814a8aa81494cd177241549c213ba3
Approve transaction hash: 0x8de79f5111b185fe67090f904b72f3dda7814a8aa81494cd177241549c213ba3
Staking tokens.....
0xea3e9f779b56a6972ce393d44cbfb4a72e74f5ef00c9b5ddfa6b86bdecf4eecb
Stake transaction hash: 0xea3e9f779b56a6972ce393d44cbfb4a72e74f5ef00c9b5ddfa6b86bdecf4eecb
```

### 7. Start the staked node

```bash
./masa-node
```

The Is staked flag will change to `true`

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

### 7. View swagger API

```bash
http://localhost:8080/swagger/index.html
```
