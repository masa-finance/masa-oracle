<!-- Title -->
<h1 align="center">
  <img src="https://emojicdn.elk.sh/🌐" width="30" /> Masa Oracle: Decentralized AI Data and LLM Network
</h1>

<!-- Description -->
<p align="center">The Masa Oracle provides infrastructure for AI developers to access decentralized data sets and decentralized LLMs. Oracle node workers can provide compute to the network by offering Twitter data, public web data, and LLM services.</p>

<!-- Table of Contents -->
<h2 align="center">
  Table of Contents
</h2>

<p align="center">
  <a href="#getting-started">Getting Started</a> •
  <a href="#staking-tokens">Staking Tokens</a> •
  <a href="#running-the-node">Running the Node</a> •
  <a href="#updates-additional-information">Updates & Additional Information</a>
</p>

<!-- Getting Started -->
<h2 align="center" id="getting-started">
  Getting Started
</h2>

<p align="center">Initiate your participation in the Masa Oracle by setting up and running your own node. Follow the steps below to integrate seamlessly into the decentralized data protocol.</p>

<h3 align="center">Prerequisites</h3>

<p align="center">Ensure these tools are installed:</p>

<ul align="center">
  <li><strong>Go</strong>: Download from <a href="https://golang.org/dl/">Go's official site</a>.</li>
  <li><strong>Yarn</strong>: Install via <a href="https://classic.yarnpkg.com/en/docs/install/">Yarn's official site</a>.</li>
  <li><strong>Git</strong>: Essential for cloning the repository.</li>
</ul>

<p align="center">For comprehensive instructions on building, staking, and running a node with Docker, see <a href="./DOCKER.md">DOCKER.md</a></p>

<h3 align="center">Installation</h3>

<p align="center">Choose your setup:</p>

<h4 align="center">Local Setup</h4>

<ol align="center">
  <li>Clone the repository:
    <pre><code>git clone https://github.com/masa-finance/masa-oracle.git</code></pre>
  </li>
  <li>Build the Go code into the masa-node binary:
    <pre><code>go build -v -o masa-node ./cmd/masa-node</code></pre>
  </li>
  <li>Set environment variables to join the testnet and optionally set the RPC URL:
    <pre><code>export ENV=test
export RPC_URL=https://1rpc.io/sepolia</code></pre>
  </li>
  <li>Start the masa-node with the necessary bootnodes:
    <pre><code>./masa-node --bootnodes=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa</code></pre>
  </li>
</ol>

<h3 align="center">Docker Setup</h3>

<p align="center">Refer to <a href="./DOCKER.md">DOCKER.md</a> for Docker-specific installation and operation.</p>

<!-- Staking Tokens -->
<h2 align="center" id="staking-tokens">
  Staking Tokens
</h2>

<p align="center">Secure your participation and potential rewards in the network by staking tokens. Detailed steps for both local and Docker setups are available.</p>

<!-- Running the Node -->
<h2 align="center" id="running-the-node">
  Running the Node
</h2>

<p align="center">Instructions for initiating your node and connecting it to the Masa network. Verify the correct functioning of your setup via the logs and ensure your node's continuous operation.</p>

<!-- Updates & Additional Information -->
<h2 align="center" id="updates-additional-information">
  Updates & Additional Information
</h2>

<p align="center">Keep your node and knowledge up-to-date by following the latest developments and updates. Regularly check back for new features and community insights.</p>

<p align="center">For more details, refer to the Masa Node <a href="CLI.md">CLI Documentation</a> and <a href="LLM.md">Twitter Sentiment Analysis Documentation</a>.</p>

<!-- API Swagger Docs -->
<h2 align="center">
  API Swagger Docs
</h2>

<p align="center">
  Access the Masa node's API documentation here:
  <pre><code>http://<masa-node>:8080/swagger/index.html</code></pre>
</p>

<!-- Consensus and Rewards -->
<h2 align="center">
  Consensus & Rewards
</h2>

<p align="center">Discover how consensus is maintained and how rewards are structured and distributed within the Masa network. These mechanisms ensure fairness and encourage active participation.</p>

<!-- Footer -->
<p align="center">Join the Masa community for support and networking. Visit our <a href="https://discord.gg/masafinance">Discord</a> for real-time interaction with fellow developers and the Masa team.</p>
