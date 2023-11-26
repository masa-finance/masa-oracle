pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract StakingContract {
    IERC20 public token;
    mapping(address => uint256) public stakes;

    event Staked(address indexed user, uint256 amount);

    constructor(address tokenAddress) {
        token = IERC20(tokenAddress);
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