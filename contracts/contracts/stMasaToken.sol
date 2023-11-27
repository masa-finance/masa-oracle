// contracts/contracts/stMASAToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetMinterPauser.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract stMasaToken is ERC20PresetMinterPauser {
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");

    constructor(address admin) ERC20PresetMinterPauser("Staked Masa Token", "stMASA") {
        _setupRole(DEFAULT_ADMIN_ROLE, admin);
        _setupRole(MINTER_ROLE, admin);
        _setupRole(PAUSER_ROLE, admin);
        _setupRole(BURNER_ROLE, admin);

        renounceRole(DEFAULT_ADMIN_ROLE, _msgSender());
        renounceRole(MINTER_ROLE, _msgSender());
        renounceRole(PAUSER_ROLE, _msgSender());
        renounceRole(BURNER_ROLE, _msgSender());
    }

    function burn(uint256 amount) public virtual override {
        require(hasRole(BURNER_ROLE, _msgSender()), "Must have burner role to burn");
        super.burn(amount);
    }
}