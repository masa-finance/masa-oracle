const MasaToken = artifacts.require("MasaToken");

module.exports = function(deployer, network, accounts) {
  const admin = accounts[0]; // for example, let's take the first account as admin
  deployer.deploy(MasaToken, admin);
};