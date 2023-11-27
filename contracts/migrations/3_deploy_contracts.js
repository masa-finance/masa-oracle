const MasaToken = artifacts.require("MasaToken");

module.exports = function(deployer, network, accounts) {
  const admin = accounts[0];
  deployer.deploy(MasaToken, admin);
};