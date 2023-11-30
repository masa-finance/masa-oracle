// contracts/contracts/OracleNodeStakingContract.sol
// SPDX-License-Identifier: MIT
// This contract is used for staking tokens in the Oracle Node
pragma solidity ^0.8.7;

// Importing required libraries and contracts
import "@openzeppelin/contracts/token/ERC20/IERC20.sol"; // Interface for ERC20 tokens
import "@openzeppelin/contracts/security/ReentrancyGuard.sol"; // To prevent re-entrancy attacks
import "./MasaToken.sol"; // The token to be staked
import "./stMasaToken.sol"; // The token representing the stake

// The contract inherits from ReentrancyGuard to prevent re-entrancy attacks
contract OracleNodeStakingContract is ReentrancyGuard {
    IERC20 public stakingToken; // The token to be staked
    stMasaToken public stakingTokenRepresentation; // The token representing the stake

    // Mapping to keep track of the stakes of each address
    mapping(address => uint256) public stakes;

    // Events to be emitted when tokens are staked or withdrawn
    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount);

    // The constructor sets the staking token and the token representing the stake
    constructor(address _stakingToken, address _stakingTokenRepresentation) {
        stakingToken = IERC20(_stakingToken);
        stakingTokenRepresentation = stMasaToken(_stakingTokenRepresentation);
    }

    // Function to stake tokens. It updates the stake, transfers the tokens to the contract, mints the token representing the stake, and emits the Staked event
    function stake(uint256 amount) external nonReentrant {
        stakes[msg.sender] += amount;
        stakingToken.transferFrom(msg.sender, address(this), amount);
        stakingTokenRepresentation.mint(msg.sender, amount);
        emit Staked(msg.sender, amount);
    }

    // Function to withdraw staked tokens. It checks if the amount to be withdrawn is less than or equal to the staked amount, updates the stake, transfers the tokens back to the user, burns the token representing the stake, and emits the Withdrawn event
    function withdraw(uint256 amount) external nonReentrant {
        require(stakes[msg.sender] >= amount, "Withdraw amount exceeds staked amount");
        stakes[msg.sender] -= amount;
        stakingToken.transfer(msg.sender, amount);
        stakingTokenRepresentation.burn(amount);
        emit Withdrawn(msg.sender, amount);
    }

    // Function to check the balance of a particular account - this is called by the oracle node to check isStaked status
    function balanceOf(address account) external view returns (uint256) {
        return stakes[account];
    }
}