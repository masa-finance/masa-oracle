#!/bin/bash

# Update and install dependencies
sudo apt-get update
sudo apt-get install -y git wget curl

# Install Go 1.22
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js and npm
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Clone the repository
git clone https://github.com/masa-finance/masa-oracle.git ~/masa-node
cd ~/masa-node

# Install node smart contracts
cd contracts
npm install
cd ..

# Set up environment variables
cat << EOF > .env
BOOTNODES=/ip4/44.209.96.127/udp/4001/quic-v1/p2p/16Uiu2HAmFF8FCaciAiU3WiodmgxoZ5ibPo2azvR6DPgoftxwsHHA
ENV=test
FILE_PATH=.
PORT=8080
RPC_URL=https://ethereum-sepolia.publicnode.com
TWITTER_SCRAPER=true
TWITTER_USERNAME=your_username
TWITTER_PASSWORD=your_password
TWITTER_2FA_CODE=""
EOF

# Create .masa directory and twitter_cookies.json
mkdir -p ~/.masa
echo '{}' > ~/.masa/twitter_cookies.json

# Build the Masa node
make build