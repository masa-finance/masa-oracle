# Masa Node Installation Guide

This guide provides steps to install, set up, and run a Masa node on a Google Cloud Platform (GCP) instance.

## Prerequisites

- A GCP account with a Compute Engine instance already set up:
  - Machine type: e2-medium (2 vCPU, 4 GB memory)
  - Boot disk: 50 GB
  - Operating System: Ubuntu 20.04 LTS or later
- The gcloud CLI installed and configured on your local machine
  - For installation instructions, visit: [Install the gcloud CLI](https://cloud.google.com/sdk/docs/install)
- Configured firewall rules to allow:
  - Port 8080 (TCP/UDP)
  - Port 4001 (TCP/UDP)

To SSH into your GCP instance using the gcloud command line:

1. Open a terminal on your local machine.

2. Run the following command, replacing `[INSTANCE_NAME]` and `[ZONE]` with your specific instance details:

   ```bash
   gcloud compute ssh [INSTANCE_NAME] --zone=[ZONE]
   ```

   For example:
   ```bash
   gcloud compute ssh masa-node --zone=us-central1-a
   ```

3. If prompted, allow gcloud to create a SSH key pair.

4. You should now be connected to your GCP instance via SSH.

## Installation Steps

1. SSH into your GCP instance.

2. Create the installation script:
   ```
   nano ~/install_masa.sh
   ```
3. Configure the .env file in the script:
```bash
# Set up environment variables for the default Twitter configuration
cat << EOF > .env
BOOTNODES=
ENV=test
FILE_PATH=.
PORT=8080
RPC_URL=https://ethereum-sepolia.publicnode.com
TWITTER_SCRAPER=true
TWITTER_USERNAME=your_username
TWITTER_PASSWORD=your_password
TWITTER_2FA_CODE=""
EOF
```
You can add more environment variables to the .env file if you want to configure the node differently. Refer to the [.env.example](../../.env.example) file in the root of our repository for a comprehensive list of available configuration options and their descriptions.

3. Copy the contents of the [install_masa.sh](../../install_masa.sh) script from the root of our repository into this file. Save and exit (Ctrl+X, then Y, then Enter).

4. Make the script executable:
   ```
   chmod +x ~/install_masa.sh
   ```

5. Run the installation script:
   ```
   ~/install_masa.sh
   ```

6. The script will:
   - Update system packages
   - Install Go 1.22
   - Install Node.js and npm
   - Clone the Masa Oracle repository
   - Set up environment variables
   - Build and run the Masa node

## Post-Installation Setup

1. After installation, check your home directory:
   ```
   ls ~
   ```
   You should see:
   - `go1.22.0.linux-amd64.tar.gz`
   - `install_masa.sh`
   - `masa-node` directory

2. Navigate to the `.masa` directory:
   ```
   cd ~/.masa
   ```

3. List the contents:
   ```
   ls
   ```
   You should see `twitter_cookies.json`.

4. Edit the `twitter_cookies.json` file:
   ```
   nano twitter_cookies.json
   ```

5. Open the [`twitter_cookies.example.json`](../../twitter_cookies.example.json) file in the root of the repository and follow these instructions to fill out the `twitter_cookies.json` file:

   1. Log in to Twitter in your web browser
   2. Open the browser's developer tools (usually F12 or right-click > Inspect)
   3. Go to the "Application" or "Storage" tab
   4. In the left sidebar, expand "Cookies" and click on "https://twitter.com"
   5. Look for the cookie names listed in the example file and copy their values
   6. Replace the "X" placeholders in the "Value" field with the actual values
   7. Save the file as "twitter_cookies.json" (remove ".example" from the filename)

   Note: Most browsers only show the Name, Value, Domain, and Path fields in the developer tools.
   The other fields (Expires, MaxAge, Secure, HttpOnly, SameSite) may not be visible or editable.
   You can leave these fields as they are in the template.

   IMPORTANT: Be extremely cautious with your auth_token and other sensitive cookies. 
   Never share them publicly or commit them to version control.

After following these steps, paste the contents into your `twitter_cookies.json` file and save.

6. To view the contents of the file:
   ```
   cat twitter_cookies.json
   ```
## Build and Run Instructions

1. Navigate to the masa-node directory:
   ```
   cd ~/masa-node
   ```

2. Build the Masa node:
   ```
   make build
   ```

3. Run the Masa node to get the ETH address to send Sepolia ETH to:
   ```
   make run
   ```

4. To get Sepolia MASA tokens you can use the faucet:
   ```
   make faucet
   ```

5. To stake your node:
   ```
   make stake
   ```

6. To run your node:
   ```
   make run
   ```

7. Get your username:
   ```
   whoami
   ```
   Remember this username for the next step.

8. To run the Masa node and ensure it restarts if the container restarts, create a systemd service:
   ```
   sudo nano /etc/systemd/system/masa-node.service
   ```

9. Add the following content to the file:
   ```
   [Unit]
   Description=Masa Node
   After=network.target

   [Service]
   Environment="PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
   ExecStart=/usr/bin/make -C /home/your_username/masa-node run
   Restart=always
   User=your_username
   WorkingDirectory=/home/your_username/masa-node

   [Install]
   WantedBy=multi-user.target
   ```
   Replace `your_username` with the username you got from step 7.

10. Save and exit the editor.

11. Enable and start the service:
    ```
    sudo systemctl enable masa-node
    sudo systemctl start masa-node
    ```

12. Check the status of the service:
    ```
    sudo systemctl status masa-node
    ```

Certainly, I'll add information about restarting the service to the merged notes:

## Important Notes

- The Masa node files are located in `~/masa-node`
- Logs can be found in two ways:
  1. Using the systemd journal: `journalctl -u masa-node`
  2. In the log file: `~/masa-node.log`
- To tail the logs in real-time:
  - For systemd: `sudo journalctl -fu masa-node`
  - For log file: `tail -f ~/masa-node.log`
- Update the Twitter credentials in `~/masa-node/.env` if using the Twitter scraper functionality
- Ensure necessary ports (default 8080) are open in GCP firewall rules
- For production use, review and adjust security settings as needed

## Troubleshooting

If you encounter issues:
1. Check the logs:
   - Systemd service logs: `journalctl -u masa-node`
   - Log file: `cat ~/masa-node.log`
2. Ensure all prerequisites are correctly installed
3. Verify your GCP instance has sufficient resources
4. If Go is not found, try: `source ~/.bashrc`

## Restarting the Service

To restart the Masa node service:

1. Use the systemctl command:
   ```
   sudo systemctl restart masa-node
   ```

2. To check the status of the service after restarting:
   ```
   sudo systemctl status masa-node
   ```

3. If you need to stop the service:
   ```
   sudo systemctl stop masa-node
   ```

4. To start the service:
   ```
   sudo systemctl start masa-node
   ```

5. If you've made changes to the service configuration, reload the systemd manager:
   ```
   sudo systemctl daemon-reload
   ```
   Then restart the service as shown in step 1.

For further assistance, consult documentation or community Discord.