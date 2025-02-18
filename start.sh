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

# Function to check node status
check_node_status() {
    echo "Starting node to check status..."
    
    # Start in detached mode
    docker compose up -d
    
    # Wait for node to initialize (max 30 seconds)
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for node status..."
    while [ $attempt -le $max_attempts ]; do
        if docker compose logs masa-node 2>&1 | grep -q "Is Staked:"; then
            local logs=$(docker compose logs masa-node 2>&1)
            local PUBLIC_KEY=$(echo "$logs" | grep "Public Key:" | tail -n1 | awk '{print $NF}')
            local IS_STAKED=$(echo "$logs" | grep "Is Staked:" | tail -n1 | grep -q "true" && echo "true" || echo "false")
            local HAS_TOKENS=$(check_token_balance "$logs" && echo "true" || echo "false")
            check_twitter_status "$logs"
            
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
                    echo "3. Run './start.sh --faucet' to request MASA tokens"
                    echo "4. Then run './start.sh --stake' to stake the node"
                else
                    echo "You have sufficient MASA tokens!"
                    echo "1. Run './start.sh --stake' to stake the node"
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
    echo "Now run './start.sh --stake' to stake the node"
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
    echo "✅ Staking successful!"
    echo "Starting miners..."
}

# Function to start a miner
start_miner() {
    local miner_id="miner_$1"
    local http_port=$((8080 + $1))
    
    # Create miner-specific directory
    mkdir -p .masa/$miner_id
    
    # Start the miner with specific ID and ports
    MINER_ID=$miner_id docker compose up -d \
        --project-name masa-$miner_id
}

# Function to check Twitter scraper status
check_twitter_status() {
    local logs=$1
    echo "Checking Twitter scraper status..."
    
    if echo "$logs" | grep -q "Is TwitterScraper:.*true"; then
        echo "✅ Twitter scraper is enabled"
        # Look for Twitter activity in logs
        if echo "$logs" | grep -q "Twitter login successful\|Scraping Twitter\|Twitter API"; then
            echo "✅ Twitter scraper is active"
            return 0
        else
            echo "⚠️  Twitter scraper enabled but no activity seen yet"
            echo "Waiting for Twitter activity..."
            return 0  # Still return success as it's enabled
        fi
    else
        echo "❌ Twitter scraper is not enabled"
        return 1
    fi
}

# Debug function
debug_node() {
    echo "=== Environment Check ==="
    echo "1. Checking .env file:"
    if [ -f .env ]; then
        echo "✅ .env file exists"
        echo "Contents:"
        cat .env
    else
        echo "❌ .env file not found"
    fi
    
    echo -e "\n2. Checking Docker configuration:"
    echo "Docker Compose Version:"
    docker compose version
    
    echo -e "\nDocker Compose Config:"
    docker compose config
    
    echo -e "\n3. Checking container logs:"
    if docker compose ps -q masa-node >/dev/null 2>&1; then
        echo "Container logs:"
        docker compose logs masa-node
    else
        echo "Container not running"
    fi
    
    echo -e "\n4. Checking .masa directory:"
    ls -la .masa/
}

# Update restart_node function with more logging
restart_node() {
    echo "=== Restarting Node ==="
    echo "1. Current configuration:"
    echo "TWITTER_SCRAPER: $(grep TWITTER_SCRAPER .env)"
    echo "TWITTER_ACCOUNTS: $(grep TWITTER_ACCOUNTS .env)"
    echo "USER_AGENTS: $(grep USER_AGENTS .env)"
    
    echo -e "\n2. Stopping containers..."
    docker compose down
    sleep 2
    
    echo -e "\n3. Starting containers..."
    docker compose up -d
    
    echo -e "\n4. Container status:"
    docker compose ps
    
    # Wait for node to initialize (max 30 seconds)
    local max_attempts=30
    local attempt=1
    
    echo -e "\n5. Waiting for node to initialize..."
    while [ $attempt -le $max_attempts ]; do
        if docker compose logs masa-node 2>&1 | grep -q "Is Staked:"; then
            local logs=$(docker compose logs masa-node 2>&1)
            local IS_STAKED=$(echo "$logs" | grep "Is Staked:" | tail -n1 | grep -q "true" && echo "true" || echo "false")
            
            echo -e "\n6. Node initialization complete!"
            echo "Full startup logs:"
            echo "$logs"
            
            check_twitter_status "$logs"
            if [ "$IS_STAKED" = "true" ]; then
                echo "✅ Node is staked and running"
                return 0
            else
                echo "❌ Node is not staked"
                return 1
            fi
        fi
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo -e "\nTimeout waiting for node to restart"
    echo "Last container logs:"
    docker compose logs masa-node
    return 1
}

# Function to check API status and get tweets
check_tweets() {
    echo "Checking Twitter API..."
    local max_attempts=5
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "http://localhost:8080/tweets" > /dev/null; then
            echo "Recent tweets:"
            curl -s "http://localhost:8080/tweets" | jq '.'
            return 0
        fi
        echo "Waiting for API to be ready..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo "❌ Could not connect to API"
    return 1
}

# Function to run with debug
run_debug() {
    echo "Starting node in debug mode..."
    docker compose down
    sleep 2
    
    echo "Starting with debug flags..."
    docker compose up  # Run without -d to see all output
}

rebuild_node() {
    echo "=== Rebuilding Node ==="
    docker compose down
    docker compose build
    docker compose up -d
}

# Function to clean up Docker resources
cleanup_docker() {
    echo "=== Cleaning up Docker resources ==="
    echo "Stopping all containers..."
    docker compose down
    
    echo "Cleaning up unused Docker resources..."
    docker system prune -a --volumes -f
    docker builder prune -a -f
}

# Main process
case "$1" in
    --faucet)
        run_faucet
        ;;
    --stake)
        run_stake
        if [ $? -eq 0 ] && [ $# -gt 1 ]; then
            shift
            for i in $(seq 1 $1); do
                start_miner $i
            done
        fi
        ;;
    --restart)
        restart_node
        ;;
    --rebuild)
        rebuild_node
        ;;
    --cleanup)
        cleanup_docker
        ;;
    --logs)
        docker compose logs -f masa-node
        ;;
    --tweets)
        check_tweets
        ;;
    --swagger)
        echo "Opening Swagger UI..."
        xdg-open http://localhost:8080/swagger/index.html 2>/dev/null || open http://localhost:8080/swagger/index.html 2>/dev/null || echo "Please open http://localhost:8080/swagger/index.html in your browser"
        ;;
    --debug)
        run_debug
        ;;
    --help)
        echo "Usage: ./start.sh [OPTION] [NUMBER_OF_MINERS]"
        echo
        echo "Start and manage your MASA node"
        echo
        echo "Options:"
        echo "  (no option)  Start node and check status"
        echo "  --faucet    Request MASA tokens from faucet"
        echo "  --stake     Stake your node and optionally start miners"
        echo "  --restart   Fully restart the node"
        echo "  --rebuild   Rebuild and restart the node"
        echo "  --cleanup   Clean up Docker resources"
        echo "  --logs      Follow node logs"
        echo "  --tweets    Check Twitter API status"
        echo "  --swagger   Open Swagger UI"
        echo "  --debug     Run node in debug mode"
        echo "  --help      Show this help message"
        echo
        echo "Examples:"
        echo "  ./start.sh          # Check node status"
        echo "  ./start.sh --faucet # Get MASA tokens (if needed)"
        echo "  ./start.sh --stake  # Stake your node"
        echo "  ./start.sh --stake 3 # Stake your node and start 3 miners"
        ;;
    *)
        echo "Usage: $0 [--faucet|--stake|--restart|--rebuild|--cleanup|--logs|--tweets|--swagger|--debug]"
        exit 1
        ;;
esac 