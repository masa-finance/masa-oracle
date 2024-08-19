# Quickstart: Deploying Masa Oracle Node on RunPod

## Prerequisites

1. A RunPod account (sign up at https://www.runpod.io/)
2. Basic knowledge of Docker and command-line interfaces

## Step 1: Use the Official Masa Finance Docker Image

Instead of building your own Docker image, we'll use the official Masa Finance image from Docker Hub.

## Step 2: Create a RunPod Template

1. Log in to your RunPod account.
2. Go to "Templates" in the left sidebar.
3. Click "New Template".
4. Fill in the template details:
   - Name: Masa Protocol Node
   - Image: masafinance/masa-node:latest
   - Container Disk: 10 GB (or as needed)
   - Volume Disk: 20 GB (or as needed)
   - Ports: 4001/tcp, 4001/udp, 8080/tcp

5. In the "Docker Command" field, enter:
   ```bash
   /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --validator="$VALIDATOR" --cachePath="$CACHE_PATH"
   ```

6. Add environment variables:
   - BOOTNODES
   - ENV
   - RPC_URL
   - FILE_PATH
   - VALIDATOR
   - CACHE_PATH
   - TWITTER_PASSwORD
   - TWITTER_USERNAME
   - TWITTER_2FA_CODE
   - TWITTER_SCRAPER

7. Save the template.

## Step 3: Deploy Your Node on RunPod

1. Go to "Pods" in the left sidebar.
2. Click "Deploy".
3. Select your "Masa Oracle Node" template.
4. Choose a GPU type (CPU-only might be sufficient for basic node operation).
5. Set the values for your environment variables.
6. Deploy the pod.

## Step 4: Access Your Node

1. Once deployed, go to the "Pods" page.
2. Find your Masa Oracle Node pod and click on it.
3. You can access the node's logs and terminal from this page.

## Step 5: Verify Node Operation

1. Use the provided terminal to check if your node is running correctly:
   ```bash
   docker logs masa-node
   ```

2. Look for the startup message indicating successful connection to the network.

## Step 6: Stake Your Node (if not already staked)

If your node isn't staked, you can stake it using the RunPod terminal:

```bash
docker exec masa-node /usr/bin/masa-node --stake 1000
```

Replace 1000 with the amount of MASA tokens you want to stake.

## Conclusion

You've now successfully deployed your Masa Oracle Node on RunPod using the official Masa Finance Docker image. Monitor your node's performance and logs regularly to ensure it's operating correctly. Remember to keep your environment variables and node software up to date.

For more detailed information on node operation and troubleshooting, refer to the main Masa Oracle documentation.

[Source: https://hub.docker.com/repositories/masafinance]