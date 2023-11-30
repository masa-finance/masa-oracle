const MasaToken = artifacts.require("MasaToken");
const stMasaToken = artifacts.require("stMasaToken");
const OracleNodeStakingContract = artifacts.require("OracleNodeStakingContract");

module.exports = async function(deployer, network, accounts) {
  // Get instances of already deployed MasaToken and stMasaToken
  const masaTokenInstance = await MasaToken.deployed();
  const stMasaTokenInstance = await stMasaToken.deployed();

  // Deploy OracleNodeStakingContract with the address of the deployed MasaToken and stMasaToken
  await deployer.deploy(OracleNodeStakingContract, masaTokenInstance.address, stMasaTokenInstance.address);
  const oracleNodeStakingContractInstance = await OracleNodeStakingContract.deployed();

  console.log('OracleNodeStakingContract deployed at address:', oracleNodeStakingContractInstance.address);
};