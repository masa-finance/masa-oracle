<!-- Title -->
<h1>
  <img src="https://emojicdn.elk.sh/🌐" width="30" /> Masa Oracle: Decentralized AI Data and LLM Network
</h1>

<!-- Description -->
<p>The Masa Oracle provides infrastructure for AI developers to access decentralized data sets and decentralized LLMs. Oracle node workers can provide compute to the network by offering Twitter data, public web data, and LLM services.</p>

<!-- Table of Contents -->
<h2>
  Table of Contents
</h2>

<ul>
  <li><a href="#getting-started">Getting Started</a></li>
  <li><a href="#staking-tokens">Staking Tokens</a></li>
  <li><a href="#running-the-node">Running the Node</a></li>
  <li><a href="#updates-additional-information">Updates & Additional Information</a></li>
</ul>

<!-- Getting Started -->
<h2 id="getting-started">
  Getting Started
</h2>

<p>Initiate your participation in the Masa Oracle by setting up and running your own node. Follow the steps below to integrate seamlessly into the decentralized data protocol.</p>

<h3>Prerequisites</h3>

<p>Ensure these tools are installed:</p>

<ul>
  <li><strong>Go</strong>: Download from <a href="https://golang.org/dl/">Go's official site</a>.</li>
  <li><strong>Yarn</strong>: Install via <a href="https://classic.yarnpkg.com/en/docs/install/">Yarn's official site</a>.</li>
  <li><strong>Git</strong>: Essential for cloning the repository.</li>
</ul>

<p>For comprehensive instructions on building, staking, and running a node with Docker, see <a href="./DOCKER.md">DOCKER.md</a></p>

<h3>Installation</h3>

<p>Choose your setup:</p>

<h4>Local Setup</h4>

<ol>
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

<h3>Docker Setup</h3>

<p>Refer to <a href="./DOCKER.md">DOCKER.md</a> for Docker-specific installation and operation.</p>

<!-- Staking Tokens -->
<h2 id="staking-tokens">
  Staking Tokens
</h2>

<p>Secure your participation and potential rewards in the network by staking tokens. Detailed steps for both local and Docker setups are available.</p>

<!-- Running the Node -->
<h2 id="running-the-node">
  Running the Node
</h2>

<p>Instructions for initiating your node and connecting it to the Masa network. Verify the correct functioning of your setup via the logs and ensure your node's continuous operation.</p>

<!-- Updates & Additional Information -->
<h2 id="updates-additional-information">
  Updates & Additional Information
</h2>

<p>Keep your node and knowledge up-to-date by following the latest developments and updates. Regularly check back for new features and community insights.</p>

<p>For more details, refer to the Masa Node <a href="CLI.md">CLI Documentation</a> and <a href="LLM.md">Twitter Sentiment Analysis Documentation</a>.</p>

<!-- API Swagger Docs -->
<h2>
  API Swagger Docs
</h2>

<p>
  Access the Masa node's API documentation here:
  <pre><code>http://<masa-node>:8080/swagger/index.html</code></pre>
</p>

<!-- Consensus and Rewards -->
<h2>
  Consensus & Rewards
</h2>

<p>Discover how consensus is maintained and how rewards are structured and distributed within the Masa network. These mechanisms ensure fairness and encourage active participation.</p>

<!-- Footer -->
<p
