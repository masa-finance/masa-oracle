require("dotenv").config();

const {
  MNEMONIC,
  PROJECT_ID,
  POLYGONSCANAPIKEY,
  BSCSCANAPIKEY,
  BASESCANAPIKEY,
  ETHERSCANAPIKEY,
  CELOSCANAPIKEY,
} = process.env;

const HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 7545,
      network_id: "*",
    },
    mumbai: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `https://polygon-mumbai.infura.io/v3/20c37320d65e45d2b8f314d8fdec0a5e`
        ),
      network_id: 80001,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    bnbtestnet: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `https://data-seed-prebsc-1-s1.binance.org:8545/`
        ),
      network_id: 97,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    bnbmainnet: {
      provider: () =>
        new HDWalletProvider(MNEMONIC, `https://bsc-dataseed1.binance.org/`),
      network_id: 56,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    opbnbtestnet: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `wss://opbnb-testnet.nodereal.io/ws/v1/99613329b67d43e3a52f5ebe7c666efc`
        ),
      network_id: 5611,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
      baseUrl: "https://api-opbnb-testnet.bscscan.com/api",
      gas: 15000000, // Add this line to increase the gas limit
    },
    opbnbmainnet: {
      provider: () =>
        new HDWalletProvider(MNEMONIC, `https://opbnb-mainnet-rpc.bnbchain.org`),
      network_id: 204,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    baseMainnet: {
      provider: () =>
        new HDWalletProvider(MNEMONIC, `https://mainnet.base.org`),
      network_id: 8453,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
      baseUrl: "https://api.basescan.org/"
    },
    baseTestnet: {
      provider: () => new HDWalletProvider(MNEMONIC, `https://goerli.base.org`),
      network_id: 84531,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
      baseUrl: "https://api-goerli.basescan.org/"
    },
    celoMainnet: {
      provider: () => new HDWalletProvider(MNEMONIC, `https://forno.celo.org`),
      network_id: 42220,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    celoTestnet: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `https://alfajores-forno.celo-testnet.org`
        ),
      network_id: 44787,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    ethereumMainnet: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `https://mainnet.infura.io/v3/${PROJECT_ID}`
        ),
      network_id: 1,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
    goerliTestnet: {
      provider: () =>
        new HDWalletProvider(
          MNEMONIC,
          `https://goerli.infura.io/v3/${PROJECT_ID}`
        ),
      network_id: 5,
      confirmations: 2,
      timeoutBlocks: 200,
      skipDryRun: true,
    },
  },
  plugins: ["truffle-plugin-verify"],
  api_keys: {
    polygonscan: POLYGONSCANAPIKEY,
    bscscan: BSCSCANAPIKEY,
    basescan: BASESCANAPIKEY,
    etherscan: ETHERSCANAPIKEY,
    celoscan: CELOSCANAPIKEY,
  },
  compilers: {
    solc: {
      version: "0.8.7",
      settings: {
        optimizer: {
          enabled: true,
          runs: 200,
        },
      },
    },
  },
};
