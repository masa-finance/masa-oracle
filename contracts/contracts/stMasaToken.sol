// contracts/contracts/stMASAToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetMinterPauser.sol";

contract stMASAToken is ERC20PresetMinterPauser {
    constructor(address admin) ERC20PresetMinterPauser("Staked Masa Token", "stMASA") {
        _setupRole(DEFAULT_ADMIN_ROLE, admin);
        _setupRole(MINTER_ROLE, admin);
        _setupRole(PAUSER_ROLE, admin);

        renounceRole(DEFAULT_ADMIN_ROLE, _msgSender());
        renounceRole(MINTER_ROLE, _msgSender());
        renounceRole(PAUSER_ROLE, _msgSender());
    }
}