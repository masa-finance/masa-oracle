# Running a Masa Protocol Node on RunPod

This guide explains how to deploy your Masa Oracle Node on RunPod using Go and access it via SSH.

## Prerequisites

1. A RunPod account (sign up at https://www.runpod.io/)
2. Funds added to your RunPod account
3. Your Masa Oracle Node repository
4. A GitHub account with access to the Masa Oracle repository

## Deployment Steps

1. Log in to your RunPod account.

2. Ensure you have sufficient funds:
   - Go to the Billing section and add funds via credit card or cryptocurrency.
   - For large transactions (over $5,000), consider business invoicing options.

3. Deploy a pod:
   - Go to "Pods" in the sidebar and click "Deploy".
   - Select a CPU-only template.
   - Choose Runpod Ubuntu as the base image.
   - Set Container Disk to 10 GB and Volume Disk to 20 GB (or as needed).

4. Once the pod is deployed, click on it to view details.

5. In the pod details, find the SSH connection information in the "Connect" section.

6. Use the provided SSH command to connect to your pod. It will look like:
   ```
   ssh root@ssh.runpod.io -p 12345
   ```
   Replace `12345` with the actual port number provided.

7. When prompted, enter the SSH password provided in the pod details.

8. Generate an SSH key and add it to your GitHub account:
   ```bash
   ssh-keygen -t ed25519 -C "your_email@example.com"
   cat /root/.ssh/id_ed25519.pub
   ```
   Copy the output and add it to your GitHub account under Settings > SSH and GPG keys > New SSH key

9. Clone your repository and set up the node:
    ```bash
    git clone git@github.com:masa-finance/masa-oracle.git /masa-node
    cd /masa-node
    ```

10. Install Go 1.22:
    ```bash
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
    go version  # Verify installation
    ```

11. Install Node.js and npm:
    ```bash
    # Update package lists
    sudo apt-get update

    # Install Node.js and npm
    curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
    sudo apt-get install -y nodejs

    # Verify the installation
    node --version
    npm --version
    ```

12. Install node smart contracts
    ```bash
    cd contracts
    npm install
    ```

13. Set up your environment variables:
    ```bash
    cat << EOF > .env
    BOOTNODES=/ip4/44.209.96.127/udp/4001/quic-v1/p2p/16Uiu2HAmFF8FCaciAiU3WiodmgxoZ5ibPo2azvR6DPgoftxwsHHA
    ENV=test
    FILE_PATH=.
    PORT=8080
    RPC_URL=https://ethereum-sepolia.publicnode.com
    TWITTER_SCRAPER=true
    TWITTER_USERNAME=username
    TWITTER_PASSWORD=password
    TWITTER_2FA_CODE=""
    EOF
    ```
    Edit this file with your actual values using a text editor like nano or vim.

14. Verify the .env contents:
    ```bash
    cat .env
    ```

15. Create and copy your `twitter_cookies.json` file to the appropriate directory:
    ```bash
    # Ensure you're in the /app directory
    cd /app

    # Create the .masa directory if it doesn't exist
    mkdir -p /root/.masa

    # Create the twitter_cookies.json file
    nano twitter-cookies.json
    copy the contents of twitter_cookies.json
    ctrl-o to save
    ctrl-x to exit

    # Copy the twitter_cookies.json file
    cp twitter-cookies.json /root/.masa/twitter_cookies.json

    # Verify the file has been copied
    ls -l /root/.masa/twitter_cookies.json
    ```

16. Build the Masa node:
    ```bash
    make build
    ```

17. Run the Masa node to get the ETH address to send Sepolia ETH to:
    ```bash
    make run
    ```

18. To get Sepolia MASA tokens you can use the faucet:
    ```bash
    make faucet
    ```

19. To stake your node:
    ```bash
    make stake
    ```

20. To run your node:
    ```bash
    make run
    ```

21. To get node API endpoints:

    Get the IP address of your pod and use the proxy to access the swagger endpoints:
    
    ```
    https://kgq7ouc1dp68ym-8080.proxy.runpod.net/swagger/#/
    ```

    Add the following to your .env file in you Bittensor config (this is an example URLfor the runpod proxy):

    ```
    ORACLE_BASE_URL="https://kgq7ouc1dp68ym-8080.proxy.runpod.net/api/v1"
    ```