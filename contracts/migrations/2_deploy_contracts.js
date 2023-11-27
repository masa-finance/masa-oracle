const stMasaToken = artifacts.require("stMasaToken");

module.exports = function (deployer, network, accounts) {
    const admin = accounts[0];

    deployer.deploy(stMasaToken, admin);
};