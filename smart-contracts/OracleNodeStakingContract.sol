pragma solidity ^0.8.0;

contract OracleNodeStakingContract {
    function stakeTokens(uint256 amount) public {
        // Logic for staking tokens
    }

    function generateSignature(bytes32 hash, uint8 v, bytes32 r, bytes32 s) public view returns (address) {
        return ecrecover(hash, v, r, s);
    }
}