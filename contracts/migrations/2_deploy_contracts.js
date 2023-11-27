const stMASAToken = artifacts.require("stMASAToken");

module.exports = function (deployer, network, accounts) {
    const admin = accounts[0];

    deployer.deploy(stMASAToken, admin);
};