// contracts/contracts/stMASAToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Importing OpenZeppelin's ERC20PresetMinterPauser contract which provides basic ERC20 functionality
// and also includes minting, pausing, and access control mechanisms.
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetMinterPauser.sol";

contract stMasaToken is ERC20PresetMinterPauser {
    // Define a new role identifier for the burner role
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");

    // Define an event for logging
    event Log(string message, address account);

    // The constructor sets up the roles for the contract and emits a log event.
    constructor(address admin) ERC20PresetMinterPauser("Staked Masa Token", "stMASA") {
        // Set up the default admin, minter, pauser, and burner roles to the provided admin address
        _setupRole(DEFAULT_ADMIN_ROLE, admin);
        _setupRole(MINTER_ROLE, admin);
        _setupRole(PAUSER_ROLE, admin);
        _setupRole(BURNER_ROLE, admin);

        // Emit a log event
        emit Log("Admin address:", admin);

        // The following lines are commented out, they would remove all roles from the sender of the transaction
        // renounceRole(DEFAULT_ADMIN_ROLE, _msgSender());
        // renounceRole(MINTER_ROLE, _msgSender());
        // renounceRole(PAUSER_ROLE, _msgSender());
        // renounceRole(BURNER_ROLE, _msgSender());
    }

    // The burn function allows an account with the burner role to burn (destroy) tokens from their balance.
    // It emits a log event and then calls the burn function from the parent contract.
    function burn(uint256 amount) public virtual override {
        // Check that the caller has the burner role
        require(hasRole(BURNER_ROLE, _msgSender()), "Must have burner role to burn");
        
        // Emit a log event
        emit Log("Burner address:", _msgSender());

        // Call the burn function from the parent contract
        super.burn(amount);
    }

    // Override the burnFrom function to include burner role check
    function burnFrom(address account, uint256 amount) public virtual override {
        // Check that the caller has the burner role
        require(hasRole(BURNER_ROLE, _msgSender()), "Must have burner role to burn from");

        // Emit a log event
        emit Log("Burner address:", _msgSender());

        // Call the burnFrom function from the parent contract
        super.burnFrom(account, amount);
    }
}