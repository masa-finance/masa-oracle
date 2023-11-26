// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./MasaToken.sol"; // Import the MasaToken contract

contract StakingContract {
    IERC20 public token;
    mapping(address => uint256) public stakes;

    event Staked(address indexed user, uint256 amount);

    constructor(MasaToken tokenContract) { // Pass the MasaToken contract as a constructor argument
        token = IERC20(address(tokenContract)); // Cast the MasaToken contract to an IERC20 interface
    }

    function stakeTokens(uint256 amount) public {
        // Transfer the tokens to this contract
        require(token.transferFrom(msg.sender, address(this), amount), "Transfer failed");

        // Update the user's stake
        stakes[msg.sender] += amount;

        // Emit the Staked event
        emit Staked(msg.sender, amount);
    }
}