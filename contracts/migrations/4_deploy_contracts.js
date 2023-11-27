const MasaToken = artifacts.require("MasaToken");
const OracleNodeStakingContract = artifacts.require("OracleNodeStakingContract");

module.exports = async function(deployer) {
  // Get the deployed instance of MasaToken
  const instance = await MasaToken.deployed();
  
  // Deploy OracleNodeStakingContract with the address of the deployed MasaToken
  const stakingInstance = await deployer.deploy(OracleNodeStakingContract, instance.address);
  
  console.log('OracleNodeStakingContract deployed at address:', stakingInstance.address);
};