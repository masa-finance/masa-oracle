#!/bin/bash

# Ensure .masa directory exists
mkdir -p .masa

# Function to check token balance
check_token_balance() {
    local logs=$1
    # Add debug output
    echo "Checking token balance in logs..."
    echo "$logs" | grep -i "balance\|masa" || true
    
    # Check for different possible balance formats
    if echo "$logs" | grep -qi "masa.*balance.*[1-9]"; then
        local balance=$(echo "$logs" | grep -i "masa.*balance" | tail -n1 | grep -o '[0-9]\+')
        echo "Found balance: $balance"
        if [ -n "$balance" ] && [ "$balance" -gt "0" ]; then
            return 0  # Has tokens
        fi
    fi
    
    # Also check for "Token balance" format
    if echo "$logs" | grep -q "Token balance:"; then
        local balance=$(echo "$logs" | grep "Token balance:" | tail -n1 | awk '{print $NF}')
        echo "Found token balance: $balance"
        if [ "$balance" != "0" ]; then
            return 0  # Has tokens
        fi
    fi
    
    return 1  # No tokens
}

# Function to start and check node status
check_node_status() {
    echo "Starting node to check status..."
    
    # Start in detached mode
    docker compose up -d
    
    # Wait for node to initialize (max 30 seconds)
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for node status..."
    while [ $attempt -le $max_attempts ]; do
        if docker compose logs masa-oracle 2>&1 | grep -q "Is Staked:"; then
            local logs=$(docker compose logs masa-oracle 2>&1)
            local PUBLIC_KEY=$(echo "$logs" | grep "Public Key:" | tail -n1 | awk '{print $NF}')
            local IS_STAKED=$(echo "$logs" | grep "Is Staked:" | tail -n1 | grep -q "true" && echo "true" || echo "false")
            local HAS_TOKENS=$(check_token_balance "$logs" && echo "true" || echo "false")
            
            echo "Node status check complete:"
            echo "Public Key: $PUBLIC_KEY"
            echo "Is Staked: $IS_STAKED"
            echo "Has Tokens: $HAS_TOKENS"
            
            if [ "$IS_STAKED" = "false" ]; then
                echo "Node is not staked. Shutting down..."
                docker compose down
                echo
                echo "To stake this node:"
                if [ "$HAS_TOKENS" = "false" ]; then
                    echo "1. Send 0.05 Sepolia ETH to: $PUBLIC_KEY"
                    echo "2. Get Sepolia ETH from: https://sepoliafaucet.com/"
                    echo "3. Run './stake.sh --faucet' to request MASA tokens"
                    echo "4. Then run './stake.sh --stake' to stake the node"
                else
                    echo "You have sufficient MASA tokens!"
                    echo "1. Run './stake.sh --stake' to stake the node"
                fi
                exit 1
            fi
            
            echo "✅ Node is staked and running!"
            return 0
        fi
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo "Timeout waiting for node status. Shutting down..."
    docker compose down
    exit 1
}

# Function to run faucet
run_faucet() {
    echo "Starting node to request tokens..."
    docker compose up -d
    sleep 10  # Give node time to initialize
    
    echo "Requesting funds from faucet..."
    if ! docker compose exec masa-node /usr/bin/masa-node --faucet --env sepolia --masaDir /home/masa/.masa; then
        echo "❌ Faucet request failed. Please check the output above."
        docker compose down
        exit 1
    fi
    echo "✅ Faucet request successful!"
    docker compose down
    echo "Now run './stake.sh --stake' to stake the node"
}

# Function to stake node
run_stake() {
    echo "Starting node to stake..."
    docker compose up -d
    sleep 10  # Give node time to initialize
    
    echo "Staking node..."
    if ! docker compose exec masa-node /usr/bin/masa-node --stake 1000 --env sepolia --masaDir /home/masa/.masa; then
        echo "❌ Staking failed. Please check the output above."
        docker compose down
        exit 1
    fi
    echo "✅ Staking successful! Node will continue running."
}

# Main process
case "$1" in
    --faucet)
        run_faucet
        ;;
    --stake)
        run_stake
        ;;
    *)
        check_node_status
        ;;
esac 