const MasaToken = artifacts.require("MasaToken");
const StakingContract = artifacts.require("StakingContract");

module.exports = function(deployer) {
  // Get the deployed instance of MasaToken
  MasaToken.deployed().then(function(instance) {
    // Deploy StakingContract with the address of the deployed MasaToken
    return deployer.deploy(StakingContract, instance.address);
  }).then(function(stakingInstance) {
    console.log('StakingContract deployed at address:', stakingInstance.address);
  });
};