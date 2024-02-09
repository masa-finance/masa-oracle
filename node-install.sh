#!/bin/bash

# Create the 'masa' user and set up home directory
useradd -m masa

# set RPC_URL
RPC_URL=https://ethereum-sepolia.publicnode.com 

# Append the RPC_URL to the masa user's .bash_profile
echo "export RPC_URL=${RPC_URL}" | tee -a /home/masa/.bash_profile

# Set permissions for the masa user's home directory
chown masa:masa /home/masa/.bash_profile

# Build go binary
go build -v -o masa-node ./cmd/masa-node
cp masa-node /usr/local/bin/masa-node
chmod +x /usr/local/bin/masa-node

# Install Node.js and Yarn
curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor -o /usr/share/keyrings/yarn-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/yarn-archive-keyring.gpg] https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
apt-get update -y && apt-get install -y yarn nodejs jq
npm install -g npm@10.4.0

# Determine global npm modules path and set NODE_PATH
GLOBAL_NODE_MODULES=$(npm root -g)
export NODE_PATH=$GLOBAL_NODE_MODULES

# Install the contracts npm module
cd /home/masa/contracts/
npm install 

MASANODE_CMD="/usr/bin/masa-node --port=4001 --udp=true --tcp=false --start --bootnodes=${BOOTNODES}"

# Create a systemd service file for masa-node
cat <<EOF | sudo tee /etc/systemd/system/masa-node.service
[Unit]
Description=MASA Node Service
After=network.target

[Service]
User=masa
WorkingDirectory=/home/masa
Environment="RPC_URL=${RPC_URL}"
ExecStart=$MASANODE_CMD
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Ensure the service file is owned by root
sudo chown root:root /etc/systemd/system/masa-node.service

# Reload the systemd daemon
sudo systemctl daemon-reload

# Enable and start the masa-node service
sudo systemctl enable masa-node
sudo systemctl start masa-node
