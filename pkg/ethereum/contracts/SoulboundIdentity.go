// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PaymentGatewayPaymentParams is an auto generated low-level Go binding around an user-defined struct.
type PaymentGatewayPaymentParams struct {
	SwapRouter            common.Address
	WrappedNativeToken    common.Address
	StableCoin            common.Address
	MasaToken             common.Address
	ProjectFeeReceiver    common.Address
	ProtocolFeeReceiver   common.Address
	ProtocolFeeAmount     *big.Int
	ProtocolFeePercent    *big.Int
	ProtocolFeePercentSub *big.Int
}

// EthereumMetaData contains all meta data concerning the Ethereum contract.
var EthereumMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseTokenURI\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"swapRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wrappedNativeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableCoin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"masaToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"projectFeeReceiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"protocolFeeReceiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePercent\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePercentSub\",\"type\":\"uint256\"}],\"internalType\":\"structPaymentGateway.PaymentParams\",\"name\":\"paymentParams\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyAdded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"IdentityAlreadyCreated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"}],\"name\":\"InvalidPaymentMethod\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc20token\",\"type\":\"address\"}],\"name\":\"NonExistingErc20Token\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotLinkedToAnIdentitySBT\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentParamsNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProtocolFeeReceiverNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RefundFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SameValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SoulNameContractNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UserMustHaveProtocolOrProjectAdminRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROJECT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkPriceMASA\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_paymentMethod\",\"type\":\"address\"}],\"name\":\"disablePaymentMethod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_paymentMethod\",\"type\":\"address\"}],\"name\":\"enablePaymentMethod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"enabledPaymentMethod\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"enabledPaymentMethods\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEnabledPaymentMethods\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExtension\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getIdentityId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"}],\"name\":\"getMintPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"}],\"name\":\"getMintPriceWithProtocolFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getProtocolFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getProtocolFeeSub\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSoulName\",\"outputs\":[{\"internalType\":\"contractISoulName\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getSoulNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"sbtNames\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getSoulNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"sbtNames\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getTokenData\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"sbtName\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"linked\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"identityId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationDate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"isAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"available\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"masaToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"yearsPeriod\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"mintIdentityWithName\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentMethod\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"yearsPeriod\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"mintIdentityWithName\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mintPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mintPriceMASA\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"projectFeeReceiver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeePercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeePercentSub\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeeReceiver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLinkPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLinkPriceMASA\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_addLinkPrice\",\"type\":\"uint256\"}],\"name\":\"setAddLinkPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_addLinkPriceMASA\",\"type\":\"uint256\"}],\"name\":\"setAddLinkPriceMASA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_masaToken\",\"type\":\"address\"}],\"name\":\"setMasaToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"}],\"name\":\"setMintPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_mintPriceMASA\",\"type\":\"uint256\"}],\"name\":\"setMintPriceMASA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_projectFeeReceiver\",\"type\":\"address\"}],\"name\":\"setProjectFeeReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_protocolFeeAmount\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeeAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_protocolFeePercent\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeePercent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_protocolFeePercentSub\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeePercentSub\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_protocolFeeReceiver\",\"type\":\"address\"}],\"name\":\"setProtocolFeeReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_queryLinkPrice\",\"type\":\"uint256\"}],\"name\":\"setQueryLinkPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_queryLinkPriceMASA\",\"type\":\"uint256\"}],\"name\":\"setQueryLinkPriceMASA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISoulName\",\"name\":\"_soulName\",\"type\":\"address\"}],\"name\":\"setSoulName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_soulboundIdentity\",\"type\":\"address\"}],\"name\":\"setSoulboundIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stableCoin\",\"type\":\"address\"}],\"name\":\"setStableCoin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_swapRouter\",\"type\":\"address\"}],\"name\":\"setSwapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wrappedNativeToken\",\"type\":\"address\"}],\"name\":\"setWrappedNativeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"soulName\",\"outputs\":[{\"internalType\":\"contractISoulName\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"soulboundIdentity\",\"outputs\":[{\"internalType\":\"contractISoulboundIdentity\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stableCoin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"tokenOfOwner\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wrappedNativeToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162004d0138038062004d01833981016040819052620000349162000592565b8484848460008585858585858584848783620000518683620001ae565b8051600180546001600160a01b03199081166001600160a01b039384161790915560208084015160028054841691851691909117905560408401516003805484169185169190911790556060840151600480548416918516919091179055608084015160078054841691851691909117905560a084015160088054909316931692909217905560c082015160095560e0820151600a5561010090910151600b558351620001059250600c9185019062000262565b5080516200011b90600d90602084019062000262565b506200012d91506000905087620001ae565b82516200014290601490602086019062000262565b5050601580546001600160a01b0319166001600160a01b0392909216919091179055506200019792507f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a69150889050620001ae565b50506001601d5550620006b2975050505050505050565b620001ba828262000237565b62000233576000828152602081815260408083206001600160a01b03851684529091529020805460ff19166001179055620001f23390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b6000828152602081815260408083206001600160a01b038516845290915290205460ff165b92915050565b828054620002709062000681565b90600052602060002090601f016020900481019282620002945760008555620002df565b82601f10620002af57805160ff1916838001178555620002df565b82800160010185558215620002df579182015b82811115620002df578251825591602001919060010190620002c2565b50620002ed929150620002f1565b5090565b5b80821115620002ed5760008155600101620002f2565b60006001600160a01b0382166200025c565b620003258162000308565b81146200033157600080fd5b50565b80516200025c816200031a565b601f01601f191690565b634e487b7160e01b600052604160045260246000fd5b6200036c8262000341565b81018181106001600160401b03821117156200038c576200038c6200034b565b6040525050565b60006200039f60405190565b9050620003ad828262000361565b919050565b60006001600160401b03821115620003ce57620003ce6200034b565b620003d98262000341565b60200192915050565b60005b83811015620003ff578181015183820152602001620003e5565b838111156200040f576000848401525b50505050565b60006200042c6200042684620003b2565b62000393565b905082815260208101848484011115620004495762000449600080fd5b62000456848285620003e2565b509392505050565b600082601f830112620004745762000474600080fd5b81516200048684826020860162000415565b949350505050565b8062000325565b80516200025c816200048e565b60006101208284031215620004ba57620004ba600080fd5b620004c761012062000393565b90506000620004d7848462000334565b8252506020620004ea8484830162000334565b6020830152506040620005008482850162000334565b6040830152506060620005168482850162000334565b60608301525060806200052c8482850162000334565b60808301525060a0620005428482850162000334565b60a08301525060c0620005588482850162000495565b60c08301525060e06200056e8482850162000495565b60e083015250610100620005858482850162000495565b6101008301525092915050565b60008060008060006101a08688031215620005b057620005b0600080fd5b6000620005be888862000334565b95505060208601516001600160401b03811115620005df57620005df600080fd5b620005ed888289016200045e565b94505060408601516001600160401b038111156200060e576200060e600080fd5b6200061c888289016200045e565b93505060608601516001600160401b038111156200063d576200063d600080fd5b6200064b888289016200045e565b92505060806200065e88828901620004a2565b9150509295509295909350565b634e487b7160e01b600052602260045260246000fd5b6002810460018216806200069657607f821691505b60208210811415620006ac57620006ac6200066b565b50919050565b61463f80620006c26000396000f3fe6080604052600436106103615760003560e01c8062bdfde51461036657806301ffc9a7146103885780630513c3e9146103be57806306fdde03146103eb5780630f2e68af1461040d578063102005191461043a578063126ed01c1461045c57806313150b4814610489578063135f470c1461049f57806317fcb39b146104b557806318160ddd146104d55780631830e881146104ea5780631f37c12414610500578063217a2c7b1461051657806323af4e1714610536578063248a9ca314610556578063289c686b14610576578063294cdf0d146105965780632f2ff15d146105b65780632f745c59146105d657806336568abe146105f657806339a51be5146106165780633ad3033e146106365780633c72ae7014610656578063412736571461067657806341c04d5e1461069657806342966c68146106b857806346877b1a146106d857806346b2b087146106f85780634962a1581461072a5780634cf12d261461074a5780634f558e791461076a5780634f6ccce71461078a5780635141453e146107aa5780636352211e146107bd5780636817c76c146107dd5780636a627842146107f35780636bfd499f1461080657806370a0823114610826578063719d0f2b1461084657806376ad199714610866578063776ce6a114610886578063776d1a541461089b57806377bed5ed146108b15780637a0d1646146108d15780637db8cb68146109015780637e669891146109215780638d0184611461094e5780638ec9c93b1461096e57806391d1485414610984578063920ffa26146109a457806393702f33146109c457806394a665e9146109e457806395d89b4114610a04578063965306aa14610a1957806398acb9a914610a39578063992642e514610a4c57806399b589cb14610a6c578063a217fddf14610a8c578063a498342114610aa1578063b507d48114610ac1578063b79636b614610adf578063b97d6b2314610aff578063c1177d1914610b15578063c31c9c0714610b35578063c86aadb614610b55578063c87b56dd14610b75578063d539139314610b95578063d547741f14610bb7578063d6e6eb9f14610bd7578063da058ae314610bed578063eb93e85514610c0d578063ebda439614610c3b578063ee1fe2ad14610c5b578063ee7a9ec514610c6e578063f4a0a52814610c8e578063fd48ac8314610cae575b600080fd5b34801561037257600080fd5b5061038661038136600461364c565b610cce565b005b34801561039457600080fd5b506103a86103a3366004613688565b610d02565b6040516103b591906136b3565b60405180910390f35b3480156103ca57600080fd5b506103de6103d936600461364c565b610d13565b6040516103b591906136e1565b3480156103f757600080fd5b50610400610d3d565b6040516103b59190613759565b34801561041957600080fd5b50601e5461042d906001600160a01b031681565b6040516103b591906137a2565b34801561044657600080fd5b5061044f610dcf565b6040516103b5919061380d565b34801561046857600080fd5b5061047c61047736600461364c565b610e30565b6040516103b59190613824565b34801561049557600080fd5b5061047c601b5481565b3480156104ab57600080fd5b5061047c600b5481565b3480156104c157600080fd5b506002546103de906001600160a01b031681565b3480156104e157600080fd5b5060125461047c565b3480156104f657600080fd5b5061047c60175481565b34801561050c57600080fd5b5061047c60185481565b34801561052257600080fd5b5061047c610531366004613846565b610e3b565b34801561054257600080fd5b50610386610551366004613883565b610e4e565b34801561056257600080fd5b5061047c61057136600461364c565b610eab565b34801561058257600080fd5b5061038661059136600461364c565b610ec0565b3480156105a257600080fd5b5061047c6105b1366004613883565b610f33565b3480156105c257600080fd5b506103866105d13660046138a4565b610f40565b3480156105e257600080fd5b5061047c6105f1366004613846565b610f61565b34801561060257600080fd5b506103866106113660046138a4565b610fbc565b34801561062257600080fd5b506008546103de906001600160a01b031681565b34801561064257600080fd5b50610386610651366004613883565b610ff2565b34801561066257600080fd5b5061038661067136600461364c565b61104f565b34801561068257600080fd5b50610386610691366004613883565b6110c2565b3480156106a257600080fd5b5061047c6000805160206145ca83398151915281565b3480156106c457600080fd5b506103866106d336600461364c565b61111f565b3480156106e457600080fd5b506103866106f3366004613883565b611151565b34801561070457600080fd5b506107186107133660046139c5565b6111ae565b6040516103b5969594939291906139ff565b34801561073657600080fd5b5061038661074536600461364c565b61127d565b34801561075657600080fd5b506104006107653660046139c5565b6112f0565b34801561077657600080fd5b506103a861078536600461364c565b6113ba565b34801561079657600080fd5b5061047c6107a536600461364c565b6113c5565b61047c6107b8366004613a53565b611413565b3480156107c957600080fd5b506103de6107d836600461364c565b611458565b3480156107e957600080fd5b5061047c60165481565b61047c610801366004613883565b611463565b34801561081257600080fd5b5061038661082136600461364c565b611470565b34801561083257600080fd5b5061047c610841366004613883565b6114a4565b34801561085257600080fd5b5061047c610861366004613883565b6114e8565b34801561087257600080fd5b50610386610881366004613883565b6115e3565b34801561089257600080fd5b50610400611640565b3480156108a757600080fd5b5061047c60195481565b3480156108bd57600080fd5b5060155461042d906001600160a01b031681565b3480156108dd57600080fd5b506103a86108ec366004613883565b60056020526000908152604090205460ff1681565b34801561090d57600080fd5b5061038661091c36600461364c565b6116c6565b34801561092d57600080fd5b5061094161093c36600461364c565b611739565b6040516103b59190613b5d565b34801561095a57600080fd5b50610386610969366004613883565b6117e9565b34801561097a57600080fd5b5061047c60095481565b34801561099057600080fd5b506103a861099f3660046138a4565b611885565b3480156109b057600080fd5b506103de6109bf3660046139c5565b6118ae565b3480156109d057600080fd5b506104006109df366004613883565b61196f565b3480156109f057600080fd5b506103866109ff366004613883565b611987565b348015610a1057600080fd5b50610400611af4565b348015610a2557600080fd5b506103a8610a343660046139c5565b611b03565b61047c610a47366004613b6e565b611baf565b348015610a5857600080fd5b506003546103de906001600160a01b031681565b348015610a7857600080fd5b506007546103de906001600160a01b031681565b348015610a9857600080fd5b5061047c600081565b348015610aad57600080fd5b50610386610abc36600461364c565b611c8e565b348015610acd57600080fd5b50601e546001600160a01b031661042d565b348015610aeb57600080fd5b50610941610afa366004613883565b611cc2565b348015610b0b57600080fd5b5061047c601a5481565b348015610b2157600080fd5b5061047c610b3036600461364c565b611d1e565b348015610b4157600080fd5b506001546103de906001600160a01b031681565b348015610b6157600080fd5b50610386610b70366004613883565b611dd6565b348015610b8157600080fd5b50610400610b9036600461364c565b611e82565b348015610ba157600080fd5b5061047c6000805160206145ea83398151915281565b348015610bc357600080fd5b50610386610bd23660046138a4565b611ee8565b348015610be357600080fd5b5061047c600a5481565b348015610bf957600080fd5b50610386610c08366004613883565b611f04565b348015610c1957600080fd5b50610c2d610c28366004613883565b611f61565b6040516103b5929190613c1a565b348015610c4757600080fd5b506004546103de906001600160a01b031681565b61047c610c69366004613c35565b611f83565b348015610c7a57600080fd5b50610386610c89366004613c76565b611fba565b348015610c9a57600080fd5b50610386610ca936600461364c565b61203e565b348015610cba57600080fd5b50610386610cc936600461364c565b6120b1565b6000610cd981612124565b600954821415610cfc5760405163c23f6ccb60e01b815260040160405180910390fd5b50600955565b6000610d0d8261212e565b92915050565b60068181548110610d2357600080fd5b6000918252602090912001546001600160a01b0316905081565b6060600c8054610d4c90613cad565b80601f0160208091040260200160405190810160405280929190818152602001828054610d7890613cad565b8015610dc55780601f10610d9a57610100808354040283529160200191610dc5565b820191906000526020600020905b815481529060010190602001808311610da857829003601f168201915b5050505050905090565b60606006805480602002602001604051908101604052809291908181526020018280548015610dc557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e09575050505050905090565b6000610d0d82612153565b6000610e478383612186565b9392505050565b6000610e5981612124565b6003546001600160a01b0383811691161415610e885760405163c23f6ccb60e01b815260040160405180910390fd5b50600380546001600160a01b0319166001600160a01b0392909216919091179055565b60009081526020819052604090206001015490565b610ecb600033611885565b158015610eed5750610eeb6000805160206145ca83398151915233611885565b155b15610f0b576040516326f0f48160e01b815260040160405180910390fd5b806018541415610f2e5760405163c23f6ccb60e01b815260040160405180910390fd5b601855565b6000610d0d826000610f61565b610f4982610eab565b610f5281612124565b610f5c83836121ed565b505050565b6000610f6c836114a4565b8210610f935760405162461bcd60e51b8152600401610f8a90613d22565b60405180910390fd5b506001600160a01b03919091166000908152601060209081526040808320938352929052205490565b6001600160a01b0381163314610fe45760405162461bcd60e51b8152600401610f8a90613d7e565b610fee8282612271565b5050565b6000610ffd81612124565b6015546001600160a01b038381169116141561102c5760405163c23f6ccb60e01b815260040160405180910390fd5b50601580546001600160a01b0319166001600160a01b0392909216919091179055565b61105a600033611885565b15801561107c575061107a6000805160206145ca83398151915233611885565b155b1561109a576040516326f0f48160e01b815260040160405180910390fd5b8060195414156110bd5760405163c23f6ccb60e01b815260040160405180910390fd5b601955565b60006110cd81612124565b6001546001600160a01b03838116911614156110fc5760405163c23f6ccb60e01b815260040160405180910390fd5b50600180546001600160a01b0319166001600160a01b0392909216919091179055565b61112933826122d6565b6111455760405162461bcd60e51b8152600401610f8a90613dc5565b61114e816122f9565b50565b600061115c81612124565b6008546001600160a01b038381169116141561118b5760405163c23f6ccb60e01b815260040160405180910390fd5b50600880546001600160a01b0319166001600160a01b0392909216919091179055565b601e5460609060009081908190819081906001600160a01b03166111e557604051636d9e949f60e01b815260040160405180910390fd5b601e546040516346b2b08760e01b81526001600160a01b03909116906346b2b08790611215908a90600401613759565b60006040518083038186803b15801561122d57600080fd5b505afa158015611241573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526112699190810190613e4b565b949c939b5091995097509550909350915050565b611288600033611885565b1580156112aa57506112a86000805160206145ca83398151915233611885565b155b156112c8576040516326f0f48160e01b815260040160405180910390fd5b8060175414156112eb5760405163c23f6ccb60e01b815260040160405180910390fd5b601755565b601e546060906001600160a01b031661131c57604051636d9e949f60e01b815260040160405180910390fd5b601e546040516346b2b08760e01b81526000916001600160a01b0316906346b2b0879061134d908690600401613759565b60006040518083038186803b15801561136557600080fd5b505afa158015611379573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526113a19190810190613e4b565b505050925050506113b181611e82565b9150505b919050565b6000610d0d82612393565b60006113d060125490565b82106113ee5760405162461bcd60e51b8152600401610f8a90613f34565b6012828154811061140157611401613f44565b90600052602060002001549050919050565b601e546000906001600160a01b031661143f57604051636d9e949f60e01b815260040160405180910390fd5b61144d600086868686611baf565b90505b949350505050565b6000610d0d826123b0565b6000610d0d600083611f83565b600061147b81612124565b600b5482141561149e5760405163c23f6ccb60e01b815260040160405180910390fd5b50600b55565b60006001600160a01b0382166114cc5760405162461bcd60e51b8152600401610f8a90613f9d565b506001600160a01b03166000908152600f602052604090205490565b600060165460001480156114fc5750601754155b1561150957506000919050565b6004546001600160a01b03838116911614801561153e57506001600160a01b03821660009081526005602052604090205460ff165b801561154c57506000601754115b1561155957505060175490565b6003546001600160a01b03838116911614801561158e57506001600160a01b03821660009081526005602052604090205460ff165b1561159b57505060165490565b6001600160a01b03821660009081526005602052604090205460ff16156115c857610d0d826016546123e5565b81604051630ac29ab760e31b8152600401610f8a91906136e1565b60006115ee81612124565b6004546001600160a01b038381169116141561161d5760405163c23f6ccb60e01b815260040160405180910390fd5b50600480546001600160a01b0319166001600160a01b0392909216919091179055565b601e546040805163776ce6a160e01b815290516060926001600160a01b03169163776ce6a1916004808301926000929190829003018186803b15801561168557600080fd5b505afa158015611699573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526116c19190810190613fad565b905090565b6116d1600033611885565b1580156116f357506116f16000805160206145ca83398151915233611885565b155b15611711576040516326f0f48160e01b815260040160405180910390fd5b80601b5414156117345760405163c23f6ccb60e01b815260040160405180910390fd5b601b55565b601e546060906001600160a01b031661176557604051636d9e949f60e01b815260040160405180910390fd5b601e54604051637e66989160e01b81526001600160a01b0390911690637e66989190611795908590600401613824565b60006040518083038186803b1580156117ad57600080fd5b505afa1580156117c1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d0d91908101906140a5565b6117f4600033611885565b15801561181657506118146000805160206145ca83398151915233611885565b155b15611834576040516326f0f48160e01b815260040160405180910390fd5b6007546001600160a01b03828116911614156118635760405163c23f6ccb60e01b815260040160405180910390fd5b600780546001600160a01b0319166001600160a01b0392909216919091179055565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b601e546000906001600160a01b03166118da57604051636d9e949f60e01b815260040160405180910390fd5b601e546040516346b2b08760e01b81526000916001600160a01b0316906346b2b0879061190b908690600401613759565b60006040518083038186803b15801561192357600080fd5b505afa158015611937573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261195f9190810190613e4b565b505050925050506113b1816123b0565b6060600061197c83610f33565b90506113b181611e82565b600061199281612124565b6001600160a01b03821660009081526005602052604090205460ff166119cd57816040516318317bd560e01b8152600401610f8a91906136e1565b6001600160a01b0382166000908152600560205260408120805460ff191690555b600654811015610f5c57826001600160a01b031660068281548110611a1557611a15613f44565b6000918252602090912001546001600160a01b03161415611ae25760068054611a40906001906140f5565b81548110611a5057611a50613f44565b600091825260209091200154600680546001600160a01b039092169183908110611a7c57611a7c613f44565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506006805480611abb57611abb61410c565b600082815260209020810160001990810180546001600160a01b0319169055019055505050565b80611aec81614122565b9150506119ee565b6060600d8054610d4c90613cad565b601e546000906001600160a01b0316611b2f57604051636d9e949f60e01b815260040160405180910390fd5b601e54604051634b29835560e11b81526001600160a01b039091169063965306aa90611b5f908590600401613759565b60206040518083038186803b158015611b7757600080fd5b505afa158015611b8b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d0d919061413d565b601e546000906001600160a01b0316611bdb57604051636d9e949f60e01b815260040160405180910390fd5b611be3612578565b6000611bef8787611f83565b601e546040516303dd904360e41b81529192506001600160a01b031690633dd9043090611c2690899089908990899060040161415e565b602060405180830381600087803b158015611c4057600080fd5b505af1158015611c54573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c78919061419f565b509050611c856001601d55565b95945050505050565b6000611c9981612124565b600a54821415611cbc5760405163c23f6ccb60e01b815260040160405180910390fd5b50600a55565b601e546060906001600160a01b0316611cee57604051636d9e949f60e01b815260040160405180910390fd5b601e54604051635bcb1b5b60e11b81526001600160a01b039091169063b79636b6906117959085906004016136e1565b6015546000906001600160a01b0316611d4a57604051630d7fe67b60e41b815260040160405180910390fd5b6000611d55836123b0565b60155460405163294cdf0d60e01b81529192506001600160a01b03169063294cdf0d90611d869084906004016136e1565b60206040518083038186803b158015611d9e57600080fd5b505afa158015611db2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113b1919061419f565b6000611de181612124565b6001600160a01b03821660009081526005602052604090205460ff1615611e1b5760405163f411c32760e01b815260040160405180910390fd5b506001600160a01b03166000818152600560205260408120805460ff191660019081179091556006805491820181559091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b0319169091179055565b6060611e8d826125a2565b6000611e976125c7565b90506000815111611eb757604051806020016040528060008152506113b1565b80611ec1846125d6565b604051602001611ed29291906141e2565b6040516020818303038152906040529392505050565b611ef182610eab565b611efa81612124565b610f5c8383612271565b6000611f0f81612124565b6002546001600160a01b0383811691161415611f3e5760405163c23f6ccb60e01b815260040160405180910390fd5b50600280546001600160a01b0319166001600160a01b0392909216919091179055565b600080611f6d836114e8565b915081611f7a8484612186565b91509150915091565b600080611f8f836114a4565b1115611fb057816040516312d5c31d60e01b8152600401610f8a91906136e1565b610e478383612672565b6000611fc581612124565b6001600160a01b038216611fec5760405163d92e233d60e01b815260040160405180910390fd5b601e546001600160a01b038381169116141561201b5760405163c23f6ccb60e01b815260040160405180910390fd5b50601e80546001600160a01b0319166001600160a01b0392909216919091179055565b612049600033611885565b15801561206b57506120696000805160206145ca83398151915233611885565b155b15612089576040516326f0f48160e01b815260040160405180910390fd5b8060165414156120ac5760405163c23f6ccb60e01b815260040160405180910390fd5b601655565b6120bc600033611885565b1580156120de57506120dc6000805160206145ca83398151915233611885565b155b156120fc576040516326f0f48160e01b815260040160405180910390fd5b80601a54141561211f5760405163c23f6ccb60e01b815260040160405180910390fd5b601a55565b61114e81336126d6565b60006001600160e01b0319821663780e9d6360e01b1480610d0d5750610d0d8261272f565b600b546000901561217e57610d0d6064612178600b548561276f90919063ffffffff16565b9061277b565b506000919050565b6009546000908190156121c1576003546001600160a01b03858116911614156121b257506009546121c1565b6121be846009546123e5565b90505b600a5415610e47576114506121e66064612178600a548761276f90919063ffffffff16565b8290612787565b6121f78282611885565b610fee576000828152602081815260408083206001600160a01b03851684529091529020805460ff1916600117905561222d3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b61227b8282611885565b15610fee576000828152602081815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6000806122e2836123b0565b6001600160a01b0385811691161491505092915050565b6000612304826123b0565b905061231281600084612793565b6001600160a01b0381166000908152600f6020526040812080546001929061233b9084906140f5565b90915550506000828152600e602052604080822080546001600160a01b03191690555183916001600160a01b038416917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca59190a35050565b6000908152600e60205260409020546001600160a01b0316151590565b6000818152600e60205260408120546001600160a01b031680610d0d5760405162461bcd60e51b8152600401610f8a9061423c565b60008160008111801561240157506001546001600160a01b0316155b1561241f5760405163fca2174f60e01b815260040160405180910390fd5b60008111801561243857506002546001600160a01b0316155b156124565760405163fca2174f60e01b815260040160405180910390fd5b60008111801561246f57506003546001600160a01b0316155b1561248d5760405163fca2174f60e01b815260040160405180910390fd5b6000811180156124a657506007546001600160a01b0316155b156124c45760405163fca2174f60e01b815260040160405180910390fd5b6001600160a01b03841660009081526005602052604090205460ff1615806124f957506003546001600160a01b038581169116145b15612519578360405163961c9a4f60e01b8152600401610f8a91906136e1565b826125275760009150612571565b6001600160a01b03841661255957600254600354612552916001600160a01b0390811691168561279e565b9150612571565b6003546125529085906001600160a01b03168561279e565b5092915050565b6002601d54141561259b5760405162461bcd60e51b8152600401610f8a90614280565b6002601d55565b6125ab81612393565b61114e5760405162461bcd60e51b8152600401610f8a9061423c565b606060148054610d4c90613cad565b606060006125e38361285c565b60010190506000816001600160401b03811115612602576126026138d7565b6040519080825280601f01601f19166020018201604052801561262c576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846126655761266a565b612636565b509392505050565b60006000805160206145ea83398151915261268c81612124565b60008061269886611f61565b915091506126a7868383612932565b60006126b2601c5490565b90506126c2601c80546001019055565b6126cc8682612de6565b9695505050505050565b6126e08282611885565b610fee576126ed81612ec2565b6126f8836020612ed4565b6040516020016127099291906142a6565b60408051601f198184030181529082905262461bcd60e51b8252610f8a91600401613759565b60006001600160e01b031982166313f2a32f60e01b148061276057506001600160e01b03198216635b5e139f60e01b145b80610d0d5750610d0d8261303f565b6000610e4782846142f8565b6000610e478284614317565b6000610e47828461432b565b610f5c838383613074565b60006060806127ad868661312c565b6001546040516307c0329d60e21b81529192506001600160a01b031690631f00ca74906127e09087908590600401614343565b60006040518083038186803b1580156127f857600080fd5b505afa15801561280c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261283491908101906143d8565b91508160008151811061284957612849613f44565b6020026020010151925050509392505050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b831061289b5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6904ee2d6d415b85acef8160201b83106128c5576904ee2d6d415b85acef8160201b830492506020015b662386f26fc1000083106128e357662386f26fc10000830492506010015b6305f5e10083106128fb576305f5e100830492506008015b612710831061290f57612710830492506004015b60648310612921576064830492506002015b600a8310610d0d5760010192915050565b61293c8282612787565b60008111801561295557506001546001600160a01b0316155b156129735760405163fca2174f60e01b815260040160405180910390fd5b60008111801561298c57506002546001600160a01b0316155b156129aa5760405163fca2174f60e01b815260040160405180910390fd5b6000811180156129c357506003546001600160a01b0316155b156129e15760405163fca2174f60e01b815260040160405180910390fd5b6000811180156129fa57506007546001600160a01b0316155b15612a185760405163fca2174f60e01b815260040160405180910390fd5b82158015612a24575081155b15612a2e57612de0565b6000612a3984612153565b90506000831180612a4a5750600081115b8015612a5f57506008546001600160a01b0316155b15612a7d5760405163910af6f560e01b815260040160405180910390fd5b6001600160a01b03851660009081526005602052604090205460ff16612ab85784604051630ac29ab760e31b8152600401610f8a91906136e1565b6001600160a01b038516612d5a57612ad08484612787565b341015612afb57612ae18484612787565b60405163091a6d0f60e01b8152600401610f8a9190613824565b6000612b0785836132be565b1115612b98576007546000906001600160a01b0316612b2686846132be565b604051612b3290614412565b60006040518083038185875af1925050503d8060008114612b6f576040519150601f19603f3d011682016040523d82523d6000602084013e612b74565b606091505b5050905080612b96576040516312171d8360e31b815260040160405180910390fd5b505b8215612c22576008546040516000916001600160a01b0316908590612bbc90614412565b60006040518083038185875af1925050503d8060008114612bf9576040519150601f19603f3d011682016040523d82523d6000602084013e612bfe565b606091505b5050905080612c20576040516312171d8360e31b815260040160405180910390fd5b505b8015612cac576008546040516000916001600160a01b0316908390612c4690614412565b60006040518083038185875af1925050503d8060008114612c83576040519150601f19603f3d011682016040523d82523d6000602084013e612c88565b606091505b5050905080612caa576040516312171d8360e31b815260040160405180910390fd5b505b612cb68484612787565b341115612d55576000612cd3612ccc8686612787565b34906132be565b90506000336001600160a01b031682604051612cee90614412565b60006040518083038185875af1925050503d8060008114612d2b576040519150601f19603f3d011682016040523d82523d6000602084013e612d30565b606091505b5050905080612d5257604051633c31275160e21b815260040160405180910390fd5b50505b612dde565b6000612d6685836132be565b1115612d9a57600754612d9a9033906001600160a01b0316612d8887856132be565b6001600160a01b0389169291906132ca565b8215612dbc57600854612dbc906001600160a01b0387811691339116866132ca565b8015612dde57600854612dde906001600160a01b0387811691339116846132ca565b505b50505050565b6001600160a01b038216612e0c5760405162461bcd60e51b8152600401610f8a90614451565b612e1581612393565b15612e325760405162461bcd60e51b8152600401610f8a90614491565b612e3e60008383612793565b6001600160a01b0382166000908152600f60205260408120805460019290612e6790849061432b565b90915550506000818152600e602052604080822080546001600160a01b0319166001600160a01b038616908117909155905183927f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d412139688591a35050565b6060610d0d6001600160a01b03831660145b60606000612ee38360026142f8565b612eee90600261432b565b6001600160401b03811115612f0557612f056138d7565b6040519080825280601f01601f191660200182016040528015612f2f576020820181803683370190505b509050600360fc1b81600081518110612f4a57612f4a613f44565b60200101906001600160f81b031916908160001a905350600f60fb1b81600181518110612f7957612f79613f44565b60200101906001600160f81b031916908160001a9053506000612f9d8460026142f8565b612fa890600161432b565b90505b6001811115613020576f181899199a1a9b1b9c1cb0b131b232b360811b85600f1660108110612fdc57612fdc613f44565b1a60f81b828281518110612ff257612ff2613f44565b60200101906001600160f81b031916908160001a90535060049490941c93613019816144a1565b9050612fab565b508315610e475760405162461bcd60e51b8152600401610f8a906144ea565b60006001600160e01b03198216637965db0b60e01b1480610d0d57506301ffc9a760e01b6001600160e01b0319831614610d0d565b6001600160a01b0383166130cf576130ca81601280546000838152601360205260408120829055600182018355919091527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34440155565b6130f2565b816001600160a01b0316836001600160a01b0316146130f2576130f28382613322565b6001600160a01b03821661310957610f5c816133bf565b826001600160a01b0316826001600160a01b031614610f5c57610f5c828261346e565b6002546060906001600160a01b038481169116148061315857506002546001600160a01b038381169116145b1561322257604080516002808252606082018352600092602083019080368337019050506002549091506001600160a01b0385811691161461319a57836131a7565b6002546001600160a01b03165b816000815181106131ba576131ba613f44565b6001600160a01b0392831660209182029290920101526002548482169116146131e357826131f0565b6002546001600160a01b03165b8160018151811061320357613203613f44565b6001600160a01b03909216602092830291909101909101529050610d0d565b6040805160038082526080820190925260009160208201606080368337019050509050838160008151811061325957613259613f44565b6001600160a01b03928316602091820292909201015260025482519116908290600190811061328a5761328a613f44565b60200260200101906001600160a01b031690816001600160a01b031681525050828160028151811061320357613203613f44565b6000610e4782846140f5565b612de0846323b872dd60e01b8585856040516024016132eb939291906144fa565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526134b2565b6000600161332f846114a4565b61333991906140f5565b60008381526011602052604090205490915080821461338c576001600160a01b03841660009081526010602090815260408083208584528252808320548484528184208190558352601190915290208190555b5060009182526011602090815260408084208490556001600160a01b039094168352601081528383209183525290812055565b6012546000906133d1906001906140f5565b600083815260136020526040812054601280549394509092849081106133f9576133f9613f44565b90600052602060002001549050806012838154811061341a5761341a613f44565b60009182526020808320909101929092558281526013909152604080822084905585825281205560128054806134525761345261410c565b6001900381819060005260206000200160009055905550505050565b6000613479836114a4565b6001600160a01b039093166000908152601060209081526040808320868452825280832085905593825260119052919091209190915550565b6000613507826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166135449092919063ffffffff16565b9050805160001480613528575080806020019051810190613528919061413d565b610f5c5760405162461bcd60e51b8152600401610f8a90614569565b6060611450848460008585600080866001600160a01b0316858760405161356b9190614579565b60006040518083038185875af1925050503d80600081146135a8576040519150601f19603f3d011682016040523d82523d6000602084013e6135ad565b606091505b50915091506135be878383876135c9565b979650505050505050565b606083156136055782516135fe576001600160a01b0385163b6135fe5760405162461bcd60e51b8152600401610f8a906145b9565b5081611450565b611450838381511561361a5781518083602001fd5b8060405162461bcd60e51b8152600401610f8a9190613759565b805b811461114e57600080fd5b8035610d0d81613634565b60006020828403121561366157613661600080fd5b60006114508484613641565b6001600160e01b03198116613636565b8035610d0d8161366d565b60006020828403121561369d5761369d600080fd5b6000611450848461367d565b8015155b82525050565b60208101610d0d82846136a9565b6001600160a01b031690565b6000610d0d826136c1565b6136ad816136cd565b60208101610d0d82846136d8565b60005b8381101561370a5781810151838201526020016136f2565b83811115612de05750506000910152565b601f01601f191690565b600061372f825190565b8084526020840193506137468185602086016136ef565b61374f8161371b565b9093019392505050565b60208082528101610e478184613725565b6000610d0d61377e61377b846136c1565b90565b6136c1565b6000610d0d8261376a565b6000610d0d82613783565b6136ad8161378e565b60208101610d0d8284613799565b60006137bc83836136d8565b505060200190565b60006137ce825190565b80845260209384019383018060005b838110156138025781516137f188826137b0565b9750602083019250506001016137dd565b509495945050505050565b60208082528101610e4781846137c4565b806136ad565b60208101610d0d828461381e565b613636816136cd565b8035610d0d81613832565b6000806040838503121561385c5761385c600080fd5b6000613868858561383b565b925050602061387985828601613641565b9150509250929050565b60006020828403121561389857613898600080fd5b6000611450848461383b565b600080604083850312156138ba576138ba600080fd5b60006138c68585613641565b92505060206138798582860161383b565b634e487b7160e01b600052604160045260246000fd5b6138f68261371b565b81018181106001600160401b0382111715613913576139136138d7565b6040525050565b600061392560405190565b90506113b582826138ed565b60006001600160401b0382111561394a5761394a6138d7565b6139538261371b565b60200192915050565b82818337506000910152565b600061397b61397684613931565b61391a565b90508281526020810184848401111561399657613996600080fd5b61266a84828561395c565b600082601f8301126139b5576139b5600080fd5b8135611450848260208601613968565b6000602082840312156139da576139da600080fd5b81356001600160401b038111156139f3576139f3600080fd5b611450848285016139a1565b60c08082528101613a108189613725565b9050613a1f60208301886136a9565b613a2c604083018761381e565b613a39606083018661381e565b613a46608083018561381e565b6135be60a08301846136a9565b60008060008060808587031215613a6c57613a6c600080fd5b6000613a78878761383b565b94505060208501356001600160401b03811115613a9757613a97600080fd5b613aa3878288016139a1565b9350506040613ab487828801613641565b92505060608501356001600160401b03811115613ad357613ad3600080fd5b613adf878288016139a1565b91505092959194509250565b6000610e478383613725565b6000613b01825190565b80845260208401935083602082028501613b1b8560200190565b8060005b85811015613b505784840389528151613b388582613aeb565b94506020830160209a909a0199925050600101613b1f565b5091979650505050505050565b60208082528101610e478184613af7565b600080600080600060a08688031215613b8957613b89600080fd5b6000613b95888861383b565b9550506020613ba68882890161383b565b94505060408601356001600160401b03811115613bc557613bc5600080fd5b613bd1888289016139a1565b9350506060613be288828901613641565b92505060808601356001600160401b03811115613c0157613c01600080fd5b613c0d888289016139a1565b9150509295509295909350565b60408101613c28828561381e565b610e47602083018461381e565b60008060408385031215613c4b57613c4b600080fd5b60006138c6858561383b565b6000610d0d826136cd565b61363681613c57565b8035610d0d81613c62565b600060208284031215613c8b57613c8b600080fd5b60006114508484613c6b565b634e487b7160e01b600052602260045260246000fd5b600281046001821680613cc157607f821691505b60208210811415613cd457613cd4613c97565b50919050565b602881526000602082017f534254456e756d657261626c653a206f776e657220696e646578206f7574206f8152676620626f756e647360c01b602082015291505b5060400190565b60208082528101610d0d81613cda565b602f81526000602082017f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636581526e103937b632b9903337b91039b2b63360891b60208201529150613d1b565b60208082528101610d0d81613d32565b601e81526000602082017f5342543a2063616c6c6572206973206e6f7420746f6b656e206f776e65720000815291505b5060200190565b60208082528101610d0d81613d8e565b6000613de361397684613931565b905082815260208101848484011115613dfe57613dfe600080fd5b61266a8482856136ef565b600082601f830112613e1d57613e1d600080fd5b8151611450848260208601613dd5565b801515613636565b8051610d0d81613e2d565b8051610d0d81613634565b60008060008060008060c08789031215613e6757613e67600080fd5b86516001600160401b03811115613e8057613e80600080fd5b613e8c89828a01613e09565b9650506020613e9d89828a01613e35565b9550506040613eae89828a01613e40565b9450506060613ebf89828a01613e40565b9350506080613ed089828a01613e40565b92505060a0613ee189828a01613e35565b9150509295509295509295565b602981526000602082017f534254456e756d657261626c653a20676c6f62616c20696e646578206f7574208152686f6620626f756e647360b81b60208201529150613d1b565b60208082528101610d0d81613eee565b634e487b7160e01b600052603260045260246000fd5b602681526000602082017f5342543a2061646472657373207a65726f206973206e6f7420612076616c69648152651037bbb732b960d11b60208201529150613d1b565b60208082528101610d0d81613f5a565b600060208284031215613fc257613fc2600080fd5b81516001600160401b03811115613fdb57613fdb600080fd5b61145084828501613e09565b60006001600160401b03821115614000576140006138d7565b5060209081020190565b600061401861397684613fe7565b8381529050602080820190840283018581111561403757614037600080fd5b835b818110156140775780516001600160401b0381111561405a5761405a600080fd5b8086016140678982613e09565b8552505060209283019201614039565b5050509392505050565b600082601f83011261409557614095600080fd5b815161145084826020860161400a565b6000602082840312156140ba576140ba600080fd5b81516001600160401b038111156140d3576140d3600080fd5b61145084828501614081565b634e487b7160e01b600052601160045260246000fd5b600082821015614107576141076140df565b500390565b634e487b7160e01b600052603160045260246000fd5b6000600019821415614136576141366140df565b5060010190565b60006020828403121561415257614152600080fd5b60006114508484613e35565b6080810161416c82876136d8565b818103602083015261417e8186613725565b905061418d604083018561381e565b81810360608301526126cc8184613725565b6000602082840312156141b4576141b4600080fd5b60006114508484613e40565b60006141ca825190565b6141d88185602086016136ef565b9290920192915050565b60006141ee82856141c0565b91506141fa82846141c0565b64173539b7b760d91b8152915060058201611450565b601581526000602082017414d0950e881a5b9d985b1a59081d1bdad95b881251605a1b81529150613dbe565b60208082528101610d0d81614210565b601f81526000602082017f5265656e7472616e637947756172643a207265656e7472616e742063616c6c0081529150613dbe565b60208082528101610d0d8161424c565b634e487b7160e01b600052601260045260246000fd5b76020b1b1b2b9b9a1b7b73a3937b61d1030b1b1b7bab73a1604d1b815260170160006142d282856141c0565b7001034b99036b4b9b9b4b733903937b6329607d1b8152601101915061145082846141c0565b6000816000190483118215151615614312576143126140df565b500290565b60008261432657614326614290565b500490565b6000821982111561433e5761433e6140df565b500190565b60408101614351828561381e565b818103602083015261145081846137c4565b600061437161397684613fe7565b8381529050602080820190840283018581111561439057614390600080fd5b835b8181101561407757806143a58882613e40565b84525060209283019201614392565b600082601f8301126143c8576143c8600080fd5b8151611450848260208601614363565b6000602082840312156143ed576143ed600080fd5b81516001600160401b0381111561440657614406600080fd5b611450848285016143b4565b6000610d0d8261377b565b601d81526000602082017f5342543a206d696e7420746f20746865207a65726f206164647265737300000081529150613dbe565b60208082528101610d0d8161441d565b601981526000602082017814d0950e881d1bdad95b88185b1c9958591e481b5a5b9d1959603a1b81529150613dbe565b60208082528101610d0d81614461565b6000816144b0576144b06140df565b506000190190565b60208082527f537472696e67733a20686578206c656e67746820696e73756666696369656e7491019081526000613dbe565b60208082528101610d0d816144b8565b6060810161450882866136d8565b61451560208301856136d8565b611450604083018461381e565b602a81526000602082017f5361666545524332303a204552433230206f7065726174696f6e20646964206e8152691bdd081cdd58d8d9595960b21b60208201529150613d1b565b60208082528101610d0d81614522565b6000610e4782846141c0565b601d81526000602082017f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000081529150613dbe565b60208082528101610d0d8161458556fe52eafc11f6f81f86878bffd31109a0d92f37506527754f00788853ff9f63b1309f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6a2646970667358221220114efb2066dcf0edcbbc770ad0954d558b8811ae581b239c7cabbdccb9d8dae264736f6c63430008080033",
}

// EthereumABI is the input ABI used to generate the binding from.
// Deprecated: Use EthereumMetaData.ABI instead.
var EthereumABI = EthereumMetaData.ABI

// EthereumBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthereumMetaData.Bin instead.
var EthereumBin = EthereumMetaData.Bin

// DeployEthereum deploys a new Ethereum contract, binding an instance of Ethereum to it.
func DeployEthereum(auth *bind.TransactOpts, backend bind.ContractBackend, admin common.Address, name string, symbol string, baseTokenURI string, paymentParams PaymentGatewayPaymentParams) (common.Address, *types.Transaction, *Ethereum, error) {
	parsed, err := EthereumMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthereumBin), backend, admin, name, symbol, baseTokenURI, paymentParams)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ethereum{EthereumCaller: EthereumCaller{contract: contract}, EthereumTransactor: EthereumTransactor{contract: contract}, EthereumFilterer: EthereumFilterer{contract: contract}}, nil
}

// Ethereum is an auto generated Go binding around an Ethereum contract.
type Ethereum struct {
	EthereumCaller     // Read-only binding to the contract
	EthereumTransactor // Write-only binding to the contract
	EthereumFilterer   // Log filterer for contract events
}

// EthereumCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthereumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthereumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthereumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthereumSession struct {
	Contract     *Ethereum         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthereumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthereumCallerSession struct {
	Contract *EthereumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// EthereumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthereumTransactorSession struct {
	Contract     *EthereumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// EthereumRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthereumRaw struct {
	Contract *Ethereum // Generic contract binding to access the raw methods on
}

// EthereumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthereumCallerRaw struct {
	Contract *EthereumCaller // Generic read-only contract binding to access the raw methods on
}

// EthereumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthereumTransactorRaw struct {
	Contract *EthereumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthereum creates a new instance of Ethereum, bound to a specific deployed contract.
func NewEthereum(address common.Address, backend bind.ContractBackend) (*Ethereum, error) {
	contract, err := bindEthereum(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ethereum{EthereumCaller: EthereumCaller{contract: contract}, EthereumTransactor: EthereumTransactor{contract: contract}, EthereumFilterer: EthereumFilterer{contract: contract}}, nil
}

// NewEthereumCaller creates a new read-only instance of Ethereum, bound to a specific deployed contract.
func NewEthereumCaller(address common.Address, caller bind.ContractCaller) (*EthereumCaller, error) {
	contract, err := bindEthereum(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumCaller{contract: contract}, nil
}

// NewEthereumTransactor creates a new write-only instance of Ethereum, bound to a specific deployed contract.
func NewEthereumTransactor(address common.Address, transactor bind.ContractTransactor) (*EthereumTransactor, error) {
	contract, err := bindEthereum(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumTransactor{contract: contract}, nil
}

// NewEthereumFilterer creates a new log filterer instance of Ethereum, bound to a specific deployed contract.
func NewEthereumFilterer(address common.Address, filterer bind.ContractFilterer) (*EthereumFilterer, error) {
	contract, err := bindEthereum(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthereumFilterer{contract: contract}, nil
}

// bindEthereum binds a generic wrapper to an already deployed contract.
func bindEthereum(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthereumMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ethereum *EthereumRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ethereum.Contract.EthereumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ethereum *EthereumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ethereum.Contract.EthereumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ethereum *EthereumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ethereum.Contract.EthereumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ethereum *EthereumCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ethereum.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ethereum *EthereumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ethereum.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ethereum *EthereumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ethereum.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Ethereum.Contract.DEFAULTADMINROLE(&_Ethereum.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Ethereum.Contract.DEFAULTADMINROLE(&_Ethereum.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Ethereum *EthereumSession) MINTERROLE() ([32]byte, error) {
	return _Ethereum.Contract.MINTERROLE(&_Ethereum.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCallerSession) MINTERROLE() ([32]byte, error) {
	return _Ethereum.Contract.MINTERROLE(&_Ethereum.CallOpts)
}

// PROJECTADMINROLE is a free data retrieval call binding the contract method 0x41c04d5e.
//
// Solidity: function PROJECT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCaller) PROJECTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "PROJECT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PROJECTADMINROLE is a free data retrieval call binding the contract method 0x41c04d5e.
//
// Solidity: function PROJECT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumSession) PROJECTADMINROLE() ([32]byte, error) {
	return _Ethereum.Contract.PROJECTADMINROLE(&_Ethereum.CallOpts)
}

// PROJECTADMINROLE is a free data retrieval call binding the contract method 0x41c04d5e.
//
// Solidity: function PROJECT_ADMIN_ROLE() view returns(bytes32)
func (_Ethereum *EthereumCallerSession) PROJECTADMINROLE() ([32]byte, error) {
	return _Ethereum.Contract.PROJECTADMINROLE(&_Ethereum.CallOpts)
}

// AddLinkPrice is a free data retrieval call binding the contract method 0x1f37c124.
//
// Solidity: function addLinkPrice() view returns(uint256)
func (_Ethereum *EthereumCaller) AddLinkPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "addLinkPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AddLinkPrice is a free data retrieval call binding the contract method 0x1f37c124.
//
// Solidity: function addLinkPrice() view returns(uint256)
func (_Ethereum *EthereumSession) AddLinkPrice() (*big.Int, error) {
	return _Ethereum.Contract.AddLinkPrice(&_Ethereum.CallOpts)
}

// AddLinkPrice is a free data retrieval call binding the contract method 0x1f37c124.
//
// Solidity: function addLinkPrice() view returns(uint256)
func (_Ethereum *EthereumCallerSession) AddLinkPrice() (*big.Int, error) {
	return _Ethereum.Contract.AddLinkPrice(&_Ethereum.CallOpts)
}

// AddLinkPriceMASA is a free data retrieval call binding the contract method 0x776d1a54.
//
// Solidity: function addLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCaller) AddLinkPriceMASA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "addLinkPriceMASA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AddLinkPriceMASA is a free data retrieval call binding the contract method 0x776d1a54.
//
// Solidity: function addLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumSession) AddLinkPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.AddLinkPriceMASA(&_Ethereum.CallOpts)
}

// AddLinkPriceMASA is a free data retrieval call binding the contract method 0x776d1a54.
//
// Solidity: function addLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCallerSession) AddLinkPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.AddLinkPriceMASA(&_Ethereum.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Ethereum *EthereumCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Ethereum *EthereumSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Ethereum.Contract.BalanceOf(&_Ethereum.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Ethereum *EthereumCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Ethereum.Contract.BalanceOf(&_Ethereum.CallOpts, owner)
}

// EnabledPaymentMethod is a free data retrieval call binding the contract method 0x7a0d1646.
//
// Solidity: function enabledPaymentMethod(address ) view returns(bool)
func (_Ethereum *EthereumCaller) EnabledPaymentMethod(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "enabledPaymentMethod", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EnabledPaymentMethod is a free data retrieval call binding the contract method 0x7a0d1646.
//
// Solidity: function enabledPaymentMethod(address ) view returns(bool)
func (_Ethereum *EthereumSession) EnabledPaymentMethod(arg0 common.Address) (bool, error) {
	return _Ethereum.Contract.EnabledPaymentMethod(&_Ethereum.CallOpts, arg0)
}

// EnabledPaymentMethod is a free data retrieval call binding the contract method 0x7a0d1646.
//
// Solidity: function enabledPaymentMethod(address ) view returns(bool)
func (_Ethereum *EthereumCallerSession) EnabledPaymentMethod(arg0 common.Address) (bool, error) {
	return _Ethereum.Contract.EnabledPaymentMethod(&_Ethereum.CallOpts, arg0)
}

// EnabledPaymentMethods is a free data retrieval call binding the contract method 0x0513c3e9.
//
// Solidity: function enabledPaymentMethods(uint256 ) view returns(address)
func (_Ethereum *EthereumCaller) EnabledPaymentMethods(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "enabledPaymentMethods", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EnabledPaymentMethods is a free data retrieval call binding the contract method 0x0513c3e9.
//
// Solidity: function enabledPaymentMethods(uint256 ) view returns(address)
func (_Ethereum *EthereumSession) EnabledPaymentMethods(arg0 *big.Int) (common.Address, error) {
	return _Ethereum.Contract.EnabledPaymentMethods(&_Ethereum.CallOpts, arg0)
}

// EnabledPaymentMethods is a free data retrieval call binding the contract method 0x0513c3e9.
//
// Solidity: function enabledPaymentMethods(uint256 ) view returns(address)
func (_Ethereum *EthereumCallerSession) EnabledPaymentMethods(arg0 *big.Int) (common.Address, error) {
	return _Ethereum.Contract.EnabledPaymentMethods(&_Ethereum.CallOpts, arg0)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_Ethereum *EthereumCaller) Exists(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "exists", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_Ethereum *EthereumSession) Exists(tokenId *big.Int) (bool, error) {
	return _Ethereum.Contract.Exists(&_Ethereum.CallOpts, tokenId)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_Ethereum *EthereumCallerSession) Exists(tokenId *big.Int) (bool, error) {
	return _Ethereum.Contract.Exists(&_Ethereum.CallOpts, tokenId)
}

// GetEnabledPaymentMethods is a free data retrieval call binding the contract method 0x10200519.
//
// Solidity: function getEnabledPaymentMethods() view returns(address[])
func (_Ethereum *EthereumCaller) GetEnabledPaymentMethods(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getEnabledPaymentMethods")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetEnabledPaymentMethods is a free data retrieval call binding the contract method 0x10200519.
//
// Solidity: function getEnabledPaymentMethods() view returns(address[])
func (_Ethereum *EthereumSession) GetEnabledPaymentMethods() ([]common.Address, error) {
	return _Ethereum.Contract.GetEnabledPaymentMethods(&_Ethereum.CallOpts)
}

// GetEnabledPaymentMethods is a free data retrieval call binding the contract method 0x10200519.
//
// Solidity: function getEnabledPaymentMethods() view returns(address[])
func (_Ethereum *EthereumCallerSession) GetEnabledPaymentMethods() ([]common.Address, error) {
	return _Ethereum.Contract.GetEnabledPaymentMethods(&_Ethereum.CallOpts)
}

// GetExtension is a free data retrieval call binding the contract method 0x776ce6a1.
//
// Solidity: function getExtension() view returns(string)
func (_Ethereum *EthereumCaller) GetExtension(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getExtension")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetExtension is a free data retrieval call binding the contract method 0x776ce6a1.
//
// Solidity: function getExtension() view returns(string)
func (_Ethereum *EthereumSession) GetExtension() (string, error) {
	return _Ethereum.Contract.GetExtension(&_Ethereum.CallOpts)
}

// GetExtension is a free data retrieval call binding the contract method 0x776ce6a1.
//
// Solidity: function getExtension() view returns(string)
func (_Ethereum *EthereumCallerSession) GetExtension() (string, error) {
	return _Ethereum.Contract.GetExtension(&_Ethereum.CallOpts)
}

// GetIdentityId is a free data retrieval call binding the contract method 0xc1177d19.
//
// Solidity: function getIdentityId(uint256 tokenId) view returns(uint256)
func (_Ethereum *EthereumCaller) GetIdentityId(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getIdentityId", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIdentityId is a free data retrieval call binding the contract method 0xc1177d19.
//
// Solidity: function getIdentityId(uint256 tokenId) view returns(uint256)
func (_Ethereum *EthereumSession) GetIdentityId(tokenId *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetIdentityId(&_Ethereum.CallOpts, tokenId)
}

// GetIdentityId is a free data retrieval call binding the contract method 0xc1177d19.
//
// Solidity: function getIdentityId(uint256 tokenId) view returns(uint256)
func (_Ethereum *EthereumCallerSession) GetIdentityId(tokenId *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetIdentityId(&_Ethereum.CallOpts, tokenId)
}

// GetMintPrice is a free data retrieval call binding the contract method 0x719d0f2b.
//
// Solidity: function getMintPrice(address paymentMethod) view returns(uint256 price)
func (_Ethereum *EthereumCaller) GetMintPrice(opts *bind.CallOpts, paymentMethod common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getMintPrice", paymentMethod)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMintPrice is a free data retrieval call binding the contract method 0x719d0f2b.
//
// Solidity: function getMintPrice(address paymentMethod) view returns(uint256 price)
func (_Ethereum *EthereumSession) GetMintPrice(paymentMethod common.Address) (*big.Int, error) {
	return _Ethereum.Contract.GetMintPrice(&_Ethereum.CallOpts, paymentMethod)
}

// GetMintPrice is a free data retrieval call binding the contract method 0x719d0f2b.
//
// Solidity: function getMintPrice(address paymentMethod) view returns(uint256 price)
func (_Ethereum *EthereumCallerSession) GetMintPrice(paymentMethod common.Address) (*big.Int, error) {
	return _Ethereum.Contract.GetMintPrice(&_Ethereum.CallOpts, paymentMethod)
}

// GetMintPriceWithProtocolFee is a free data retrieval call binding the contract method 0xeb93e855.
//
// Solidity: function getMintPriceWithProtocolFee(address paymentMethod) view returns(uint256 price, uint256 protocolFee)
func (_Ethereum *EthereumCaller) GetMintPriceWithProtocolFee(opts *bind.CallOpts, paymentMethod common.Address) (struct {
	Price       *big.Int
	ProtocolFee *big.Int
}, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getMintPriceWithProtocolFee", paymentMethod)

	outstruct := new(struct {
		Price       *big.Int
		ProtocolFee *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Price = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ProtocolFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMintPriceWithProtocolFee is a free data retrieval call binding the contract method 0xeb93e855.
//
// Solidity: function getMintPriceWithProtocolFee(address paymentMethod) view returns(uint256 price, uint256 protocolFee)
func (_Ethereum *EthereumSession) GetMintPriceWithProtocolFee(paymentMethod common.Address) (struct {
	Price       *big.Int
	ProtocolFee *big.Int
}, error) {
	return _Ethereum.Contract.GetMintPriceWithProtocolFee(&_Ethereum.CallOpts, paymentMethod)
}

// GetMintPriceWithProtocolFee is a free data retrieval call binding the contract method 0xeb93e855.
//
// Solidity: function getMintPriceWithProtocolFee(address paymentMethod) view returns(uint256 price, uint256 protocolFee)
func (_Ethereum *EthereumCallerSession) GetMintPriceWithProtocolFee(paymentMethod common.Address) (struct {
	Price       *big.Int
	ProtocolFee *big.Int
}, error) {
	return _Ethereum.Contract.GetMintPriceWithProtocolFee(&_Ethereum.CallOpts, paymentMethod)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0x217a2c7b.
//
// Solidity: function getProtocolFee(address paymentMethod, uint256 amount) view returns(uint256)
func (_Ethereum *EthereumCaller) GetProtocolFee(opts *bind.CallOpts, paymentMethod common.Address, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getProtocolFee", paymentMethod, amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProtocolFee is a free data retrieval call binding the contract method 0x217a2c7b.
//
// Solidity: function getProtocolFee(address paymentMethod, uint256 amount) view returns(uint256)
func (_Ethereum *EthereumSession) GetProtocolFee(paymentMethod common.Address, amount *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetProtocolFee(&_Ethereum.CallOpts, paymentMethod, amount)
}

// GetProtocolFee is a free data retrieval call binding the contract method 0x217a2c7b.
//
// Solidity: function getProtocolFee(address paymentMethod, uint256 amount) view returns(uint256)
func (_Ethereum *EthereumCallerSession) GetProtocolFee(paymentMethod common.Address, amount *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetProtocolFee(&_Ethereum.CallOpts, paymentMethod, amount)
}

// GetProtocolFeeSub is a free data retrieval call binding the contract method 0x126ed01c.
//
// Solidity: function getProtocolFeeSub(uint256 amount) view returns(uint256)
func (_Ethereum *EthereumCaller) GetProtocolFeeSub(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getProtocolFeeSub", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProtocolFeeSub is a free data retrieval call binding the contract method 0x126ed01c.
//
// Solidity: function getProtocolFeeSub(uint256 amount) view returns(uint256)
func (_Ethereum *EthereumSession) GetProtocolFeeSub(amount *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetProtocolFeeSub(&_Ethereum.CallOpts, amount)
}

// GetProtocolFeeSub is a free data retrieval call binding the contract method 0x126ed01c.
//
// Solidity: function getProtocolFeeSub(uint256 amount) view returns(uint256)
func (_Ethereum *EthereumCallerSession) GetProtocolFeeSub(amount *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.GetProtocolFeeSub(&_Ethereum.CallOpts, amount)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ethereum *EthereumCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ethereum *EthereumSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Ethereum.Contract.GetRoleAdmin(&_Ethereum.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Ethereum *EthereumCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Ethereum.Contract.GetRoleAdmin(&_Ethereum.CallOpts, role)
}

// GetSoulName is a free data retrieval call binding the contract method 0xb507d481.
//
// Solidity: function getSoulName() view returns(address)
func (_Ethereum *EthereumCaller) GetSoulName(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getSoulName")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSoulName is a free data retrieval call binding the contract method 0xb507d481.
//
// Solidity: function getSoulName() view returns(address)
func (_Ethereum *EthereumSession) GetSoulName() (common.Address, error) {
	return _Ethereum.Contract.GetSoulName(&_Ethereum.CallOpts)
}

// GetSoulName is a free data retrieval call binding the contract method 0xb507d481.
//
// Solidity: function getSoulName() view returns(address)
func (_Ethereum *EthereumCallerSession) GetSoulName() (common.Address, error) {
	return _Ethereum.Contract.GetSoulName(&_Ethereum.CallOpts)
}

// GetSoulNames is a free data retrieval call binding the contract method 0x7e669891.
//
// Solidity: function getSoulNames(uint256 tokenId) view returns(string[] sbtNames)
func (_Ethereum *EthereumCaller) GetSoulNames(opts *bind.CallOpts, tokenId *big.Int) ([]string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getSoulNames", tokenId)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSoulNames is a free data retrieval call binding the contract method 0x7e669891.
//
// Solidity: function getSoulNames(uint256 tokenId) view returns(string[] sbtNames)
func (_Ethereum *EthereumSession) GetSoulNames(tokenId *big.Int) ([]string, error) {
	return _Ethereum.Contract.GetSoulNames(&_Ethereum.CallOpts, tokenId)
}

// GetSoulNames is a free data retrieval call binding the contract method 0x7e669891.
//
// Solidity: function getSoulNames(uint256 tokenId) view returns(string[] sbtNames)
func (_Ethereum *EthereumCallerSession) GetSoulNames(tokenId *big.Int) ([]string, error) {
	return _Ethereum.Contract.GetSoulNames(&_Ethereum.CallOpts, tokenId)
}

// GetSoulNames0 is a free data retrieval call binding the contract method 0xb79636b6.
//
// Solidity: function getSoulNames(address owner) view returns(string[] sbtNames)
func (_Ethereum *EthereumCaller) GetSoulNames0(opts *bind.CallOpts, owner common.Address) ([]string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getSoulNames0", owner)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSoulNames0 is a free data retrieval call binding the contract method 0xb79636b6.
//
// Solidity: function getSoulNames(address owner) view returns(string[] sbtNames)
func (_Ethereum *EthereumSession) GetSoulNames0(owner common.Address) ([]string, error) {
	return _Ethereum.Contract.GetSoulNames0(&_Ethereum.CallOpts, owner)
}

// GetSoulNames0 is a free data retrieval call binding the contract method 0xb79636b6.
//
// Solidity: function getSoulNames(address owner) view returns(string[] sbtNames)
func (_Ethereum *EthereumCallerSession) GetSoulNames0(owner common.Address) ([]string, error) {
	return _Ethereum.Contract.GetSoulNames0(&_Ethereum.CallOpts, owner)
}

// GetTokenData is a free data retrieval call binding the contract method 0x46b2b087.
//
// Solidity: function getTokenData(string name) view returns(string sbtName, bool linked, uint256 identityId, uint256 tokenId, uint256 expirationDate, bool active)
func (_Ethereum *EthereumCaller) GetTokenData(opts *bind.CallOpts, name string) (struct {
	SbtName        string
	Linked         bool
	IdentityId     *big.Int
	TokenId        *big.Int
	ExpirationDate *big.Int
	Active         bool
}, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "getTokenData", name)

	outstruct := new(struct {
		SbtName        string
		Linked         bool
		IdentityId     *big.Int
		TokenId        *big.Int
		ExpirationDate *big.Int
		Active         bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SbtName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Linked = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.IdentityId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TokenId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.ExpirationDate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// GetTokenData is a free data retrieval call binding the contract method 0x46b2b087.
//
// Solidity: function getTokenData(string name) view returns(string sbtName, bool linked, uint256 identityId, uint256 tokenId, uint256 expirationDate, bool active)
func (_Ethereum *EthereumSession) GetTokenData(name string) (struct {
	SbtName        string
	Linked         bool
	IdentityId     *big.Int
	TokenId        *big.Int
	ExpirationDate *big.Int
	Active         bool
}, error) {
	return _Ethereum.Contract.GetTokenData(&_Ethereum.CallOpts, name)
}

// GetTokenData is a free data retrieval call binding the contract method 0x46b2b087.
//
// Solidity: function getTokenData(string name) view returns(string sbtName, bool linked, uint256 identityId, uint256 tokenId, uint256 expirationDate, bool active)
func (_Ethereum *EthereumCallerSession) GetTokenData(name string) (struct {
	SbtName        string
	Linked         bool
	IdentityId     *big.Int
	TokenId        *big.Int
	ExpirationDate *big.Int
	Active         bool
}, error) {
	return _Ethereum.Contract.GetTokenData(&_Ethereum.CallOpts, name)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ethereum *EthereumCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ethereum *EthereumSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Ethereum.Contract.HasRole(&_Ethereum.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Ethereum *EthereumCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Ethereum.Contract.HasRole(&_Ethereum.CallOpts, role, account)
}

// IsAvailable is a free data retrieval call binding the contract method 0x965306aa.
//
// Solidity: function isAvailable(string name) view returns(bool available)
func (_Ethereum *EthereumCaller) IsAvailable(opts *bind.CallOpts, name string) (bool, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "isAvailable", name)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAvailable is a free data retrieval call binding the contract method 0x965306aa.
//
// Solidity: function isAvailable(string name) view returns(bool available)
func (_Ethereum *EthereumSession) IsAvailable(name string) (bool, error) {
	return _Ethereum.Contract.IsAvailable(&_Ethereum.CallOpts, name)
}

// IsAvailable is a free data retrieval call binding the contract method 0x965306aa.
//
// Solidity: function isAvailable(string name) view returns(bool available)
func (_Ethereum *EthereumCallerSession) IsAvailable(name string) (bool, error) {
	return _Ethereum.Contract.IsAvailable(&_Ethereum.CallOpts, name)
}

// MasaToken is a free data retrieval call binding the contract method 0xebda4396.
//
// Solidity: function masaToken() view returns(address)
func (_Ethereum *EthereumCaller) MasaToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "masaToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MasaToken is a free data retrieval call binding the contract method 0xebda4396.
//
// Solidity: function masaToken() view returns(address)
func (_Ethereum *EthereumSession) MasaToken() (common.Address, error) {
	return _Ethereum.Contract.MasaToken(&_Ethereum.CallOpts)
}

// MasaToken is a free data retrieval call binding the contract method 0xebda4396.
//
// Solidity: function masaToken() view returns(address)
func (_Ethereum *EthereumCallerSession) MasaToken() (common.Address, error) {
	return _Ethereum.Contract.MasaToken(&_Ethereum.CallOpts)
}

// MintPrice is a free data retrieval call binding the contract method 0x6817c76c.
//
// Solidity: function mintPrice() view returns(uint256)
func (_Ethereum *EthereumCaller) MintPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "mintPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintPrice is a free data retrieval call binding the contract method 0x6817c76c.
//
// Solidity: function mintPrice() view returns(uint256)
func (_Ethereum *EthereumSession) MintPrice() (*big.Int, error) {
	return _Ethereum.Contract.MintPrice(&_Ethereum.CallOpts)
}

// MintPrice is a free data retrieval call binding the contract method 0x6817c76c.
//
// Solidity: function mintPrice() view returns(uint256)
func (_Ethereum *EthereumCallerSession) MintPrice() (*big.Int, error) {
	return _Ethereum.Contract.MintPrice(&_Ethereum.CallOpts)
}

// MintPriceMASA is a free data retrieval call binding the contract method 0x1830e881.
//
// Solidity: function mintPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCaller) MintPriceMASA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "mintPriceMASA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintPriceMASA is a free data retrieval call binding the contract method 0x1830e881.
//
// Solidity: function mintPriceMASA() view returns(uint256)
func (_Ethereum *EthereumSession) MintPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.MintPriceMASA(&_Ethereum.CallOpts)
}

// MintPriceMASA is a free data retrieval call binding the contract method 0x1830e881.
//
// Solidity: function mintPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCallerSession) MintPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.MintPriceMASA(&_Ethereum.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ethereum *EthereumCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ethereum *EthereumSession) Name() (string, error) {
	return _Ethereum.Contract.Name(&_Ethereum.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ethereum *EthereumCallerSession) Name() (string, error) {
	return _Ethereum.Contract.Name(&_Ethereum.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Ethereum *EthereumCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Ethereum *EthereumSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Ethereum.Contract.OwnerOf(&_Ethereum.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Ethereum *EthereumCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Ethereum.Contract.OwnerOf(&_Ethereum.CallOpts, tokenId)
}

// OwnerOf0 is a free data retrieval call binding the contract method 0x920ffa26.
//
// Solidity: function ownerOf(string name) view returns(address)
func (_Ethereum *EthereumCaller) OwnerOf0(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "ownerOf0", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf0 is a free data retrieval call binding the contract method 0x920ffa26.
//
// Solidity: function ownerOf(string name) view returns(address)
func (_Ethereum *EthereumSession) OwnerOf0(name string) (common.Address, error) {
	return _Ethereum.Contract.OwnerOf0(&_Ethereum.CallOpts, name)
}

// OwnerOf0 is a free data retrieval call binding the contract method 0x920ffa26.
//
// Solidity: function ownerOf(string name) view returns(address)
func (_Ethereum *EthereumCallerSession) OwnerOf0(name string) (common.Address, error) {
	return _Ethereum.Contract.OwnerOf0(&_Ethereum.CallOpts, name)
}

// ProjectFeeReceiver is a free data retrieval call binding the contract method 0x99b589cb.
//
// Solidity: function projectFeeReceiver() view returns(address)
func (_Ethereum *EthereumCaller) ProjectFeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "projectFeeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProjectFeeReceiver is a free data retrieval call binding the contract method 0x99b589cb.
//
// Solidity: function projectFeeReceiver() view returns(address)
func (_Ethereum *EthereumSession) ProjectFeeReceiver() (common.Address, error) {
	return _Ethereum.Contract.ProjectFeeReceiver(&_Ethereum.CallOpts)
}

// ProjectFeeReceiver is a free data retrieval call binding the contract method 0x99b589cb.
//
// Solidity: function projectFeeReceiver() view returns(address)
func (_Ethereum *EthereumCallerSession) ProjectFeeReceiver() (common.Address, error) {
	return _Ethereum.Contract.ProjectFeeReceiver(&_Ethereum.CallOpts)
}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Ethereum *EthereumCaller) ProtocolFeeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "protocolFeeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Ethereum *EthereumSession) ProtocolFeeAmount() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeeAmount(&_Ethereum.CallOpts)
}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Ethereum *EthereumCallerSession) ProtocolFeeAmount() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeeAmount(&_Ethereum.CallOpts)
}

// ProtocolFeePercent is a free data retrieval call binding the contract method 0xd6e6eb9f.
//
// Solidity: function protocolFeePercent() view returns(uint256)
func (_Ethereum *EthereumCaller) ProtocolFeePercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "protocolFeePercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeePercent is a free data retrieval call binding the contract method 0xd6e6eb9f.
//
// Solidity: function protocolFeePercent() view returns(uint256)
func (_Ethereum *EthereumSession) ProtocolFeePercent() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeePercent(&_Ethereum.CallOpts)
}

// ProtocolFeePercent is a free data retrieval call binding the contract method 0xd6e6eb9f.
//
// Solidity: function protocolFeePercent() view returns(uint256)
func (_Ethereum *EthereumCallerSession) ProtocolFeePercent() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeePercent(&_Ethereum.CallOpts)
}

// ProtocolFeePercentSub is a free data retrieval call binding the contract method 0x135f470c.
//
// Solidity: function protocolFeePercentSub() view returns(uint256)
func (_Ethereum *EthereumCaller) ProtocolFeePercentSub(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "protocolFeePercentSub")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeePercentSub is a free data retrieval call binding the contract method 0x135f470c.
//
// Solidity: function protocolFeePercentSub() view returns(uint256)
func (_Ethereum *EthereumSession) ProtocolFeePercentSub() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeePercentSub(&_Ethereum.CallOpts)
}

// ProtocolFeePercentSub is a free data retrieval call binding the contract method 0x135f470c.
//
// Solidity: function protocolFeePercentSub() view returns(uint256)
func (_Ethereum *EthereumCallerSession) ProtocolFeePercentSub() (*big.Int, error) {
	return _Ethereum.Contract.ProtocolFeePercentSub(&_Ethereum.CallOpts)
}

// ProtocolFeeReceiver is a free data retrieval call binding the contract method 0x39a51be5.
//
// Solidity: function protocolFeeReceiver() view returns(address)
func (_Ethereum *EthereumCaller) ProtocolFeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "protocolFeeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProtocolFeeReceiver is a free data retrieval call binding the contract method 0x39a51be5.
//
// Solidity: function protocolFeeReceiver() view returns(address)
func (_Ethereum *EthereumSession) ProtocolFeeReceiver() (common.Address, error) {
	return _Ethereum.Contract.ProtocolFeeReceiver(&_Ethereum.CallOpts)
}

// ProtocolFeeReceiver is a free data retrieval call binding the contract method 0x39a51be5.
//
// Solidity: function protocolFeeReceiver() view returns(address)
func (_Ethereum *EthereumCallerSession) ProtocolFeeReceiver() (common.Address, error) {
	return _Ethereum.Contract.ProtocolFeeReceiver(&_Ethereum.CallOpts)
}

// QueryLinkPrice is a free data retrieval call binding the contract method 0xb97d6b23.
//
// Solidity: function queryLinkPrice() view returns(uint256)
func (_Ethereum *EthereumCaller) QueryLinkPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "queryLinkPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QueryLinkPrice is a free data retrieval call binding the contract method 0xb97d6b23.
//
// Solidity: function queryLinkPrice() view returns(uint256)
func (_Ethereum *EthereumSession) QueryLinkPrice() (*big.Int, error) {
	return _Ethereum.Contract.QueryLinkPrice(&_Ethereum.CallOpts)
}

// QueryLinkPrice is a free data retrieval call binding the contract method 0xb97d6b23.
//
// Solidity: function queryLinkPrice() view returns(uint256)
func (_Ethereum *EthereumCallerSession) QueryLinkPrice() (*big.Int, error) {
	return _Ethereum.Contract.QueryLinkPrice(&_Ethereum.CallOpts)
}

// QueryLinkPriceMASA is a free data retrieval call binding the contract method 0x13150b48.
//
// Solidity: function queryLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCaller) QueryLinkPriceMASA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "queryLinkPriceMASA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QueryLinkPriceMASA is a free data retrieval call binding the contract method 0x13150b48.
//
// Solidity: function queryLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumSession) QueryLinkPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.QueryLinkPriceMASA(&_Ethereum.CallOpts)
}

// QueryLinkPriceMASA is a free data retrieval call binding the contract method 0x13150b48.
//
// Solidity: function queryLinkPriceMASA() view returns(uint256)
func (_Ethereum *EthereumCallerSession) QueryLinkPriceMASA() (*big.Int, error) {
	return _Ethereum.Contract.QueryLinkPriceMASA(&_Ethereum.CallOpts)
}

// SoulName is a free data retrieval call binding the contract method 0x0f2e68af.
//
// Solidity: function soulName() view returns(address)
func (_Ethereum *EthereumCaller) SoulName(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "soulName")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SoulName is a free data retrieval call binding the contract method 0x0f2e68af.
//
// Solidity: function soulName() view returns(address)
func (_Ethereum *EthereumSession) SoulName() (common.Address, error) {
	return _Ethereum.Contract.SoulName(&_Ethereum.CallOpts)
}

// SoulName is a free data retrieval call binding the contract method 0x0f2e68af.
//
// Solidity: function soulName() view returns(address)
func (_Ethereum *EthereumCallerSession) SoulName() (common.Address, error) {
	return _Ethereum.Contract.SoulName(&_Ethereum.CallOpts)
}

// SoulboundIdentity is a free data retrieval call binding the contract method 0x77bed5ed.
//
// Solidity: function soulboundIdentity() view returns(address)
func (_Ethereum *EthereumCaller) SoulboundIdentity(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "soulboundIdentity")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SoulboundIdentity is a free data retrieval call binding the contract method 0x77bed5ed.
//
// Solidity: function soulboundIdentity() view returns(address)
func (_Ethereum *EthereumSession) SoulboundIdentity() (common.Address, error) {
	return _Ethereum.Contract.SoulboundIdentity(&_Ethereum.CallOpts)
}

// SoulboundIdentity is a free data retrieval call binding the contract method 0x77bed5ed.
//
// Solidity: function soulboundIdentity() view returns(address)
func (_Ethereum *EthereumCallerSession) SoulboundIdentity() (common.Address, error) {
	return _Ethereum.Contract.SoulboundIdentity(&_Ethereum.CallOpts)
}

// StableCoin is a free data retrieval call binding the contract method 0x992642e5.
//
// Solidity: function stableCoin() view returns(address)
func (_Ethereum *EthereumCaller) StableCoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "stableCoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StableCoin is a free data retrieval call binding the contract method 0x992642e5.
//
// Solidity: function stableCoin() view returns(address)
func (_Ethereum *EthereumSession) StableCoin() (common.Address, error) {
	return _Ethereum.Contract.StableCoin(&_Ethereum.CallOpts)
}

// StableCoin is a free data retrieval call binding the contract method 0x992642e5.
//
// Solidity: function stableCoin() view returns(address)
func (_Ethereum *EthereumCallerSession) StableCoin() (common.Address, error) {
	return _Ethereum.Contract.StableCoin(&_Ethereum.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ethereum *EthereumCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ethereum *EthereumSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Ethereum.Contract.SupportsInterface(&_Ethereum.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Ethereum *EthereumCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Ethereum.Contract.SupportsInterface(&_Ethereum.CallOpts, interfaceId)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Ethereum *EthereumCaller) SwapRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "swapRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Ethereum *EthereumSession) SwapRouter() (common.Address, error) {
	return _Ethereum.Contract.SwapRouter(&_Ethereum.CallOpts)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Ethereum *EthereumCallerSession) SwapRouter() (common.Address, error) {
	return _Ethereum.Contract.SwapRouter(&_Ethereum.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ethereum *EthereumCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ethereum *EthereumSession) Symbol() (string, error) {
	return _Ethereum.Contract.Symbol(&_Ethereum.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ethereum *EthereumCallerSession) Symbol() (string, error) {
	return _Ethereum.Contract.Symbol(&_Ethereum.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Ethereum *EthereumCaller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Ethereum *EthereumSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.TokenByIndex(&_Ethereum.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Ethereum *EthereumCallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.TokenByIndex(&_Ethereum.CallOpts, index)
}

// TokenOfOwner is a free data retrieval call binding the contract method 0x294cdf0d.
//
// Solidity: function tokenOfOwner(address owner) view returns(uint256)
func (_Ethereum *EthereumCaller) TokenOfOwner(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenOfOwner", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwner is a free data retrieval call binding the contract method 0x294cdf0d.
//
// Solidity: function tokenOfOwner(address owner) view returns(uint256)
func (_Ethereum *EthereumSession) TokenOfOwner(owner common.Address) (*big.Int, error) {
	return _Ethereum.Contract.TokenOfOwner(&_Ethereum.CallOpts, owner)
}

// TokenOfOwner is a free data retrieval call binding the contract method 0x294cdf0d.
//
// Solidity: function tokenOfOwner(address owner) view returns(uint256)
func (_Ethereum *EthereumCallerSession) TokenOfOwner(owner common.Address) (*big.Int, error) {
	return _Ethereum.Contract.TokenOfOwner(&_Ethereum.CallOpts, owner)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Ethereum *EthereumCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Ethereum *EthereumSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.TokenOfOwnerByIndex(&_Ethereum.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Ethereum *EthereumCallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Ethereum.Contract.TokenOfOwnerByIndex(&_Ethereum.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0x4cf12d26.
//
// Solidity: function tokenURI(string name) view returns(string)
func (_Ethereum *EthereumCaller) TokenURI(opts *bind.CallOpts, name string) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenURI", name)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0x4cf12d26.
//
// Solidity: function tokenURI(string name) view returns(string)
func (_Ethereum *EthereumSession) TokenURI(name string) (string, error) {
	return _Ethereum.Contract.TokenURI(&_Ethereum.CallOpts, name)
}

// TokenURI is a free data retrieval call binding the contract method 0x4cf12d26.
//
// Solidity: function tokenURI(string name) view returns(string)
func (_Ethereum *EthereumCallerSession) TokenURI(name string) (string, error) {
	return _Ethereum.Contract.TokenURI(&_Ethereum.CallOpts, name)
}

// TokenURI0 is a free data retrieval call binding the contract method 0x93702f33.
//
// Solidity: function tokenURI(address owner) view returns(string)
func (_Ethereum *EthereumCaller) TokenURI0(opts *bind.CallOpts, owner common.Address) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenURI0", owner)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI0 is a free data retrieval call binding the contract method 0x93702f33.
//
// Solidity: function tokenURI(address owner) view returns(string)
func (_Ethereum *EthereumSession) TokenURI0(owner common.Address) (string, error) {
	return _Ethereum.Contract.TokenURI0(&_Ethereum.CallOpts, owner)
}

// TokenURI0 is a free data retrieval call binding the contract method 0x93702f33.
//
// Solidity: function tokenURI(address owner) view returns(string)
func (_Ethereum *EthereumCallerSession) TokenURI0(owner common.Address) (string, error) {
	return _Ethereum.Contract.TokenURI0(&_Ethereum.CallOpts, owner)
}

// TokenURI1 is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ethereum *EthereumCaller) TokenURI1(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "tokenURI1", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI1 is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ethereum *EthereumSession) TokenURI1(tokenId *big.Int) (string, error) {
	return _Ethereum.Contract.TokenURI1(&_Ethereum.CallOpts, tokenId)
}

// TokenURI1 is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ethereum *EthereumCallerSession) TokenURI1(tokenId *big.Int) (string, error) {
	return _Ethereum.Contract.TokenURI1(&_Ethereum.CallOpts, tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Ethereum *EthereumCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Ethereum *EthereumSession) TotalSupply() (*big.Int, error) {
	return _Ethereum.Contract.TotalSupply(&_Ethereum.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Ethereum *EthereumCallerSession) TotalSupply() (*big.Int, error) {
	return _Ethereum.Contract.TotalSupply(&_Ethereum.CallOpts)
}

// WrappedNativeToken is a free data retrieval call binding the contract method 0x17fcb39b.
//
// Solidity: function wrappedNativeToken() view returns(address)
func (_Ethereum *EthereumCaller) WrappedNativeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ethereum.contract.Call(opts, &out, "wrappedNativeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedNativeToken is a free data retrieval call binding the contract method 0x17fcb39b.
//
// Solidity: function wrappedNativeToken() view returns(address)
func (_Ethereum *EthereumSession) WrappedNativeToken() (common.Address, error) {
	return _Ethereum.Contract.WrappedNativeToken(&_Ethereum.CallOpts)
}

// WrappedNativeToken is a free data retrieval call binding the contract method 0x17fcb39b.
//
// Solidity: function wrappedNativeToken() view returns(address)
func (_Ethereum *EthereumCallerSession) WrappedNativeToken() (common.Address, error) {
	return _Ethereum.Contract.WrappedNativeToken(&_Ethereum.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Ethereum *EthereumTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Ethereum *EthereumSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.Burn(&_Ethereum.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Ethereum *EthereumTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.Burn(&_Ethereum.TransactOpts, tokenId)
}

// DisablePaymentMethod is a paid mutator transaction binding the contract method 0x94a665e9.
//
// Solidity: function disablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumTransactor) DisablePaymentMethod(opts *bind.TransactOpts, _paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "disablePaymentMethod", _paymentMethod)
}

// DisablePaymentMethod is a paid mutator transaction binding the contract method 0x94a665e9.
//
// Solidity: function disablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumSession) DisablePaymentMethod(_paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.DisablePaymentMethod(&_Ethereum.TransactOpts, _paymentMethod)
}

// DisablePaymentMethod is a paid mutator transaction binding the contract method 0x94a665e9.
//
// Solidity: function disablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumTransactorSession) DisablePaymentMethod(_paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.DisablePaymentMethod(&_Ethereum.TransactOpts, _paymentMethod)
}

// EnablePaymentMethod is a paid mutator transaction binding the contract method 0xc86aadb6.
//
// Solidity: function enablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumTransactor) EnablePaymentMethod(opts *bind.TransactOpts, _paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "enablePaymentMethod", _paymentMethod)
}

// EnablePaymentMethod is a paid mutator transaction binding the contract method 0xc86aadb6.
//
// Solidity: function enablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumSession) EnablePaymentMethod(_paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.EnablePaymentMethod(&_Ethereum.TransactOpts, _paymentMethod)
}

// EnablePaymentMethod is a paid mutator transaction binding the contract method 0xc86aadb6.
//
// Solidity: function enablePaymentMethod(address _paymentMethod) returns()
func (_Ethereum *EthereumTransactorSession) EnablePaymentMethod(_paymentMethod common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.EnablePaymentMethod(&_Ethereum.TransactOpts, _paymentMethod)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.GrantRole(&_Ethereum.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.GrantRole(&_Ethereum.TransactOpts, role, account)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) payable returns(uint256)
func (_Ethereum *EthereumTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) payable returns(uint256)
func (_Ethereum *EthereumSession) Mint(to common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.Mint(&_Ethereum.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) payable returns(uint256)
func (_Ethereum *EthereumTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.Mint(&_Ethereum.TransactOpts, to)
}

// Mint0 is a paid mutator transaction binding the contract method 0xee1fe2ad.
//
// Solidity: function mint(address paymentMethod, address to) payable returns(uint256)
func (_Ethereum *EthereumTransactor) Mint0(opts *bind.TransactOpts, paymentMethod common.Address, to common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "mint0", paymentMethod, to)
}

// Mint0 is a paid mutator transaction binding the contract method 0xee1fe2ad.
//
// Solidity: function mint(address paymentMethod, address to) payable returns(uint256)
func (_Ethereum *EthereumSession) Mint0(paymentMethod common.Address, to common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.Mint0(&_Ethereum.TransactOpts, paymentMethod, to)
}

// Mint0 is a paid mutator transaction binding the contract method 0xee1fe2ad.
//
// Solidity: function mint(address paymentMethod, address to) payable returns(uint256)
func (_Ethereum *EthereumTransactorSession) Mint0(paymentMethod common.Address, to common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.Mint0(&_Ethereum.TransactOpts, paymentMethod, to)
}

// MintIdentityWithName is a paid mutator transaction binding the contract method 0x5141453e.
//
// Solidity: function mintIdentityWithName(address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumTransactor) MintIdentityWithName(opts *bind.TransactOpts, to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "mintIdentityWithName", to, name, yearsPeriod, _tokenURI)
}

// MintIdentityWithName is a paid mutator transaction binding the contract method 0x5141453e.
//
// Solidity: function mintIdentityWithName(address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumSession) MintIdentityWithName(to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.Contract.MintIdentityWithName(&_Ethereum.TransactOpts, to, name, yearsPeriod, _tokenURI)
}

// MintIdentityWithName is a paid mutator transaction binding the contract method 0x5141453e.
//
// Solidity: function mintIdentityWithName(address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumTransactorSession) MintIdentityWithName(to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.Contract.MintIdentityWithName(&_Ethereum.TransactOpts, to, name, yearsPeriod, _tokenURI)
}

// MintIdentityWithName0 is a paid mutator transaction binding the contract method 0x98acb9a9.
//
// Solidity: function mintIdentityWithName(address paymentMethod, address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumTransactor) MintIdentityWithName0(opts *bind.TransactOpts, paymentMethod common.Address, to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "mintIdentityWithName0", paymentMethod, to, name, yearsPeriod, _tokenURI)
}

// MintIdentityWithName0 is a paid mutator transaction binding the contract method 0x98acb9a9.
//
// Solidity: function mintIdentityWithName(address paymentMethod, address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumSession) MintIdentityWithName0(paymentMethod common.Address, to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.Contract.MintIdentityWithName0(&_Ethereum.TransactOpts, paymentMethod, to, name, yearsPeriod, _tokenURI)
}

// MintIdentityWithName0 is a paid mutator transaction binding the contract method 0x98acb9a9.
//
// Solidity: function mintIdentityWithName(address paymentMethod, address to, string name, uint256 yearsPeriod, string _tokenURI) payable returns(uint256)
func (_Ethereum *EthereumTransactorSession) MintIdentityWithName0(paymentMethod common.Address, to common.Address, name string, yearsPeriod *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Ethereum.Contract.MintIdentityWithName0(&_Ethereum.TransactOpts, paymentMethod, to, name, yearsPeriod, _tokenURI)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.RenounceRole(&_Ethereum.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.RenounceRole(&_Ethereum.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.RevokeRole(&_Ethereum.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Ethereum *EthereumTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.RevokeRole(&_Ethereum.TransactOpts, role, account)
}

// SetAddLinkPrice is a paid mutator transaction binding the contract method 0x289c686b.
//
// Solidity: function setAddLinkPrice(uint256 _addLinkPrice) returns()
func (_Ethereum *EthereumTransactor) SetAddLinkPrice(opts *bind.TransactOpts, _addLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setAddLinkPrice", _addLinkPrice)
}

// SetAddLinkPrice is a paid mutator transaction binding the contract method 0x289c686b.
//
// Solidity: function setAddLinkPrice(uint256 _addLinkPrice) returns()
func (_Ethereum *EthereumSession) SetAddLinkPrice(_addLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetAddLinkPrice(&_Ethereum.TransactOpts, _addLinkPrice)
}

// SetAddLinkPrice is a paid mutator transaction binding the contract method 0x289c686b.
//
// Solidity: function setAddLinkPrice(uint256 _addLinkPrice) returns()
func (_Ethereum *EthereumTransactorSession) SetAddLinkPrice(_addLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetAddLinkPrice(&_Ethereum.TransactOpts, _addLinkPrice)
}

// SetAddLinkPriceMASA is a paid mutator transaction binding the contract method 0x3c72ae70.
//
// Solidity: function setAddLinkPriceMASA(uint256 _addLinkPriceMASA) returns()
func (_Ethereum *EthereumTransactor) SetAddLinkPriceMASA(opts *bind.TransactOpts, _addLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setAddLinkPriceMASA", _addLinkPriceMASA)
}

// SetAddLinkPriceMASA is a paid mutator transaction binding the contract method 0x3c72ae70.
//
// Solidity: function setAddLinkPriceMASA(uint256 _addLinkPriceMASA) returns()
func (_Ethereum *EthereumSession) SetAddLinkPriceMASA(_addLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetAddLinkPriceMASA(&_Ethereum.TransactOpts, _addLinkPriceMASA)
}

// SetAddLinkPriceMASA is a paid mutator transaction binding the contract method 0x3c72ae70.
//
// Solidity: function setAddLinkPriceMASA(uint256 _addLinkPriceMASA) returns()
func (_Ethereum *EthereumTransactorSession) SetAddLinkPriceMASA(_addLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetAddLinkPriceMASA(&_Ethereum.TransactOpts, _addLinkPriceMASA)
}

// SetMasaToken is a paid mutator transaction binding the contract method 0x76ad1997.
//
// Solidity: function setMasaToken(address _masaToken) returns()
func (_Ethereum *EthereumTransactor) SetMasaToken(opts *bind.TransactOpts, _masaToken common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setMasaToken", _masaToken)
}

// SetMasaToken is a paid mutator transaction binding the contract method 0x76ad1997.
//
// Solidity: function setMasaToken(address _masaToken) returns()
func (_Ethereum *EthereumSession) SetMasaToken(_masaToken common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMasaToken(&_Ethereum.TransactOpts, _masaToken)
}

// SetMasaToken is a paid mutator transaction binding the contract method 0x76ad1997.
//
// Solidity: function setMasaToken(address _masaToken) returns()
func (_Ethereum *EthereumTransactorSession) SetMasaToken(_masaToken common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMasaToken(&_Ethereum.TransactOpts, _masaToken)
}

// SetMintPrice is a paid mutator transaction binding the contract method 0xf4a0a528.
//
// Solidity: function setMintPrice(uint256 _mintPrice) returns()
func (_Ethereum *EthereumTransactor) SetMintPrice(opts *bind.TransactOpts, _mintPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setMintPrice", _mintPrice)
}

// SetMintPrice is a paid mutator transaction binding the contract method 0xf4a0a528.
//
// Solidity: function setMintPrice(uint256 _mintPrice) returns()
func (_Ethereum *EthereumSession) SetMintPrice(_mintPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMintPrice(&_Ethereum.TransactOpts, _mintPrice)
}

// SetMintPrice is a paid mutator transaction binding the contract method 0xf4a0a528.
//
// Solidity: function setMintPrice(uint256 _mintPrice) returns()
func (_Ethereum *EthereumTransactorSession) SetMintPrice(_mintPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMintPrice(&_Ethereum.TransactOpts, _mintPrice)
}

// SetMintPriceMASA is a paid mutator transaction binding the contract method 0x4962a158.
//
// Solidity: function setMintPriceMASA(uint256 _mintPriceMASA) returns()
func (_Ethereum *EthereumTransactor) SetMintPriceMASA(opts *bind.TransactOpts, _mintPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setMintPriceMASA", _mintPriceMASA)
}

// SetMintPriceMASA is a paid mutator transaction binding the contract method 0x4962a158.
//
// Solidity: function setMintPriceMASA(uint256 _mintPriceMASA) returns()
func (_Ethereum *EthereumSession) SetMintPriceMASA(_mintPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMintPriceMASA(&_Ethereum.TransactOpts, _mintPriceMASA)
}

// SetMintPriceMASA is a paid mutator transaction binding the contract method 0x4962a158.
//
// Solidity: function setMintPriceMASA(uint256 _mintPriceMASA) returns()
func (_Ethereum *EthereumTransactorSession) SetMintPriceMASA(_mintPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetMintPriceMASA(&_Ethereum.TransactOpts, _mintPriceMASA)
}

// SetProjectFeeReceiver is a paid mutator transaction binding the contract method 0x8d018461.
//
// Solidity: function setProjectFeeReceiver(address _projectFeeReceiver) returns()
func (_Ethereum *EthereumTransactor) SetProjectFeeReceiver(opts *bind.TransactOpts, _projectFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setProjectFeeReceiver", _projectFeeReceiver)
}

// SetProjectFeeReceiver is a paid mutator transaction binding the contract method 0x8d018461.
//
// Solidity: function setProjectFeeReceiver(address _projectFeeReceiver) returns()
func (_Ethereum *EthereumSession) SetProjectFeeReceiver(_projectFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProjectFeeReceiver(&_Ethereum.TransactOpts, _projectFeeReceiver)
}

// SetProjectFeeReceiver is a paid mutator transaction binding the contract method 0x8d018461.
//
// Solidity: function setProjectFeeReceiver(address _projectFeeReceiver) returns()
func (_Ethereum *EthereumTransactorSession) SetProjectFeeReceiver(_projectFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProjectFeeReceiver(&_Ethereum.TransactOpts, _projectFeeReceiver)
}

// SetProtocolFeeAmount is a paid mutator transaction binding the contract method 0x00bdfde5.
//
// Solidity: function setProtocolFeeAmount(uint256 _protocolFeeAmount) returns()
func (_Ethereum *EthereumTransactor) SetProtocolFeeAmount(opts *bind.TransactOpts, _protocolFeeAmount *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setProtocolFeeAmount", _protocolFeeAmount)
}

// SetProtocolFeeAmount is a paid mutator transaction binding the contract method 0x00bdfde5.
//
// Solidity: function setProtocolFeeAmount(uint256 _protocolFeeAmount) returns()
func (_Ethereum *EthereumSession) SetProtocolFeeAmount(_protocolFeeAmount *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeeAmount(&_Ethereum.TransactOpts, _protocolFeeAmount)
}

// SetProtocolFeeAmount is a paid mutator transaction binding the contract method 0x00bdfde5.
//
// Solidity: function setProtocolFeeAmount(uint256 _protocolFeeAmount) returns()
func (_Ethereum *EthereumTransactorSession) SetProtocolFeeAmount(_protocolFeeAmount *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeeAmount(&_Ethereum.TransactOpts, _protocolFeeAmount)
}

// SetProtocolFeePercent is a paid mutator transaction binding the contract method 0xa4983421.
//
// Solidity: function setProtocolFeePercent(uint256 _protocolFeePercent) returns()
func (_Ethereum *EthereumTransactor) SetProtocolFeePercent(opts *bind.TransactOpts, _protocolFeePercent *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setProtocolFeePercent", _protocolFeePercent)
}

// SetProtocolFeePercent is a paid mutator transaction binding the contract method 0xa4983421.
//
// Solidity: function setProtocolFeePercent(uint256 _protocolFeePercent) returns()
func (_Ethereum *EthereumSession) SetProtocolFeePercent(_protocolFeePercent *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeePercent(&_Ethereum.TransactOpts, _protocolFeePercent)
}

// SetProtocolFeePercent is a paid mutator transaction binding the contract method 0xa4983421.
//
// Solidity: function setProtocolFeePercent(uint256 _protocolFeePercent) returns()
func (_Ethereum *EthereumTransactorSession) SetProtocolFeePercent(_protocolFeePercent *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeePercent(&_Ethereum.TransactOpts, _protocolFeePercent)
}

// SetProtocolFeePercentSub is a paid mutator transaction binding the contract method 0x6bfd499f.
//
// Solidity: function setProtocolFeePercentSub(uint256 _protocolFeePercentSub) returns()
func (_Ethereum *EthereumTransactor) SetProtocolFeePercentSub(opts *bind.TransactOpts, _protocolFeePercentSub *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setProtocolFeePercentSub", _protocolFeePercentSub)
}

// SetProtocolFeePercentSub is a paid mutator transaction binding the contract method 0x6bfd499f.
//
// Solidity: function setProtocolFeePercentSub(uint256 _protocolFeePercentSub) returns()
func (_Ethereum *EthereumSession) SetProtocolFeePercentSub(_protocolFeePercentSub *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeePercentSub(&_Ethereum.TransactOpts, _protocolFeePercentSub)
}

// SetProtocolFeePercentSub is a paid mutator transaction binding the contract method 0x6bfd499f.
//
// Solidity: function setProtocolFeePercentSub(uint256 _protocolFeePercentSub) returns()
func (_Ethereum *EthereumTransactorSession) SetProtocolFeePercentSub(_protocolFeePercentSub *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeePercentSub(&_Ethereum.TransactOpts, _protocolFeePercentSub)
}

// SetProtocolFeeReceiver is a paid mutator transaction binding the contract method 0x46877b1a.
//
// Solidity: function setProtocolFeeReceiver(address _protocolFeeReceiver) returns()
func (_Ethereum *EthereumTransactor) SetProtocolFeeReceiver(opts *bind.TransactOpts, _protocolFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setProtocolFeeReceiver", _protocolFeeReceiver)
}

// SetProtocolFeeReceiver is a paid mutator transaction binding the contract method 0x46877b1a.
//
// Solidity: function setProtocolFeeReceiver(address _protocolFeeReceiver) returns()
func (_Ethereum *EthereumSession) SetProtocolFeeReceiver(_protocolFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeeReceiver(&_Ethereum.TransactOpts, _protocolFeeReceiver)
}

// SetProtocolFeeReceiver is a paid mutator transaction binding the contract method 0x46877b1a.
//
// Solidity: function setProtocolFeeReceiver(address _protocolFeeReceiver) returns()
func (_Ethereum *EthereumTransactorSession) SetProtocolFeeReceiver(_protocolFeeReceiver common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetProtocolFeeReceiver(&_Ethereum.TransactOpts, _protocolFeeReceiver)
}

// SetQueryLinkPrice is a paid mutator transaction binding the contract method 0xfd48ac83.
//
// Solidity: function setQueryLinkPrice(uint256 _queryLinkPrice) returns()
func (_Ethereum *EthereumTransactor) SetQueryLinkPrice(opts *bind.TransactOpts, _queryLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setQueryLinkPrice", _queryLinkPrice)
}

// SetQueryLinkPrice is a paid mutator transaction binding the contract method 0xfd48ac83.
//
// Solidity: function setQueryLinkPrice(uint256 _queryLinkPrice) returns()
func (_Ethereum *EthereumSession) SetQueryLinkPrice(_queryLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetQueryLinkPrice(&_Ethereum.TransactOpts, _queryLinkPrice)
}

// SetQueryLinkPrice is a paid mutator transaction binding the contract method 0xfd48ac83.
//
// Solidity: function setQueryLinkPrice(uint256 _queryLinkPrice) returns()
func (_Ethereum *EthereumTransactorSession) SetQueryLinkPrice(_queryLinkPrice *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetQueryLinkPrice(&_Ethereum.TransactOpts, _queryLinkPrice)
}

// SetQueryLinkPriceMASA is a paid mutator transaction binding the contract method 0x7db8cb68.
//
// Solidity: function setQueryLinkPriceMASA(uint256 _queryLinkPriceMASA) returns()
func (_Ethereum *EthereumTransactor) SetQueryLinkPriceMASA(opts *bind.TransactOpts, _queryLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setQueryLinkPriceMASA", _queryLinkPriceMASA)
}

// SetQueryLinkPriceMASA is a paid mutator transaction binding the contract method 0x7db8cb68.
//
// Solidity: function setQueryLinkPriceMASA(uint256 _queryLinkPriceMASA) returns()
func (_Ethereum *EthereumSession) SetQueryLinkPriceMASA(_queryLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetQueryLinkPriceMASA(&_Ethereum.TransactOpts, _queryLinkPriceMASA)
}

// SetQueryLinkPriceMASA is a paid mutator transaction binding the contract method 0x7db8cb68.
//
// Solidity: function setQueryLinkPriceMASA(uint256 _queryLinkPriceMASA) returns()
func (_Ethereum *EthereumTransactorSession) SetQueryLinkPriceMASA(_queryLinkPriceMASA *big.Int) (*types.Transaction, error) {
	return _Ethereum.Contract.SetQueryLinkPriceMASA(&_Ethereum.TransactOpts, _queryLinkPriceMASA)
}

// SetSoulName is a paid mutator transaction binding the contract method 0xee7a9ec5.
//
// Solidity: function setSoulName(address _soulName) returns()
func (_Ethereum *EthereumTransactor) SetSoulName(opts *bind.TransactOpts, _soulName common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setSoulName", _soulName)
}

// SetSoulName is a paid mutator transaction binding the contract method 0xee7a9ec5.
//
// Solidity: function setSoulName(address _soulName) returns()
func (_Ethereum *EthereumSession) SetSoulName(_soulName common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSoulName(&_Ethereum.TransactOpts, _soulName)
}

// SetSoulName is a paid mutator transaction binding the contract method 0xee7a9ec5.
//
// Solidity: function setSoulName(address _soulName) returns()
func (_Ethereum *EthereumTransactorSession) SetSoulName(_soulName common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSoulName(&_Ethereum.TransactOpts, _soulName)
}

// SetSoulboundIdentity is a paid mutator transaction binding the contract method 0x3ad3033e.
//
// Solidity: function setSoulboundIdentity(address _soulboundIdentity) returns()
func (_Ethereum *EthereumTransactor) SetSoulboundIdentity(opts *bind.TransactOpts, _soulboundIdentity common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setSoulboundIdentity", _soulboundIdentity)
}

// SetSoulboundIdentity is a paid mutator transaction binding the contract method 0x3ad3033e.
//
// Solidity: function setSoulboundIdentity(address _soulboundIdentity) returns()
func (_Ethereum *EthereumSession) SetSoulboundIdentity(_soulboundIdentity common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSoulboundIdentity(&_Ethereum.TransactOpts, _soulboundIdentity)
}

// SetSoulboundIdentity is a paid mutator transaction binding the contract method 0x3ad3033e.
//
// Solidity: function setSoulboundIdentity(address _soulboundIdentity) returns()
func (_Ethereum *EthereumTransactorSession) SetSoulboundIdentity(_soulboundIdentity common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSoulboundIdentity(&_Ethereum.TransactOpts, _soulboundIdentity)
}

// SetStableCoin is a paid mutator transaction binding the contract method 0x23af4e17.
//
// Solidity: function setStableCoin(address _stableCoin) returns()
func (_Ethereum *EthereumTransactor) SetStableCoin(opts *bind.TransactOpts, _stableCoin common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setStableCoin", _stableCoin)
}

// SetStableCoin is a paid mutator transaction binding the contract method 0x23af4e17.
//
// Solidity: function setStableCoin(address _stableCoin) returns()
func (_Ethereum *EthereumSession) SetStableCoin(_stableCoin common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetStableCoin(&_Ethereum.TransactOpts, _stableCoin)
}

// SetStableCoin is a paid mutator transaction binding the contract method 0x23af4e17.
//
// Solidity: function setStableCoin(address _stableCoin) returns()
func (_Ethereum *EthereumTransactorSession) SetStableCoin(_stableCoin common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetStableCoin(&_Ethereum.TransactOpts, _stableCoin)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Ethereum *EthereumTransactor) SetSwapRouter(opts *bind.TransactOpts, _swapRouter common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setSwapRouter", _swapRouter)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Ethereum *EthereumSession) SetSwapRouter(_swapRouter common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSwapRouter(&_Ethereum.TransactOpts, _swapRouter)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Ethereum *EthereumTransactorSession) SetSwapRouter(_swapRouter common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetSwapRouter(&_Ethereum.TransactOpts, _swapRouter)
}

// SetWrappedNativeToken is a paid mutator transaction binding the contract method 0xda058ae3.
//
// Solidity: function setWrappedNativeToken(address _wrappedNativeToken) returns()
func (_Ethereum *EthereumTransactor) SetWrappedNativeToken(opts *bind.TransactOpts, _wrappedNativeToken common.Address) (*types.Transaction, error) {
	return _Ethereum.contract.Transact(opts, "setWrappedNativeToken", _wrappedNativeToken)
}

// SetWrappedNativeToken is a paid mutator transaction binding the contract method 0xda058ae3.
//
// Solidity: function setWrappedNativeToken(address _wrappedNativeToken) returns()
func (_Ethereum *EthereumSession) SetWrappedNativeToken(_wrappedNativeToken common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetWrappedNativeToken(&_Ethereum.TransactOpts, _wrappedNativeToken)
}

// SetWrappedNativeToken is a paid mutator transaction binding the contract method 0xda058ae3.
//
// Solidity: function setWrappedNativeToken(address _wrappedNativeToken) returns()
func (_Ethereum *EthereumTransactorSession) SetWrappedNativeToken(_wrappedNativeToken common.Address) (*types.Transaction, error) {
	return _Ethereum.Contract.SetWrappedNativeToken(&_Ethereum.TransactOpts, _wrappedNativeToken)
}

// EthereumBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the Ethereum contract.
type EthereumBurnIterator struct {
	Event *EthereumBurn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthereumBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumBurn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthereumBurn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthereumBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumBurn represents a Burn event raised by the Ethereum contract.
type EthereumBurn struct {
	Owner   common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) FilterBurn(opts *bind.FilterOpts, _owner []common.Address, _tokenId []*big.Int) (*EthereumBurnIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Ethereum.contract.FilterLogs(opts, "Burn", _ownerRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EthereumBurnIterator{contract: _Ethereum.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *EthereumBurn, _owner []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Ethereum.contract.WatchLogs(opts, "Burn", _ownerRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumBurn)
				if err := _Ethereum.contract.UnpackLog(event, "Burn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBurn is a log parse operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) ParseBurn(log types.Log) (*EthereumBurn, error) {
	event := new(EthereumBurn)
	if err := _Ethereum.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the Ethereum contract.
type EthereumMintIterator struct {
	Event *EthereumMint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthereumMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumMint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthereumMint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthereumMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumMint represents a Mint event raised by the Ethereum contract.
type EthereumMint struct {
	Owner   common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) FilterMint(opts *bind.FilterOpts, _owner []common.Address, _tokenId []*big.Int) (*EthereumMintIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Ethereum.contract.FilterLogs(opts, "Mint", _ownerRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EthereumMintIterator{contract: _Ethereum.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *EthereumMint, _owner []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Ethereum.contract.WatchLogs(opts, "Mint", _ownerRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumMint)
				if err := _Ethereum.contract.UnpackLog(event, "Mint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _owner, uint256 indexed _tokenId)
func (_Ethereum *EthereumFilterer) ParseMint(log types.Log) (*EthereumMint, error) {
	event := new(EthereumMint)
	if err := _Ethereum.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Ethereum contract.
type EthereumRoleAdminChangedIterator struct {
	Event *EthereumRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthereumRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthereumRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthereumRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumRoleAdminChanged represents a RoleAdminChanged event raised by the Ethereum contract.
type EthereumRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ethereum *EthereumFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*EthereumRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Ethereum.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &EthereumRoleAdminChangedIterator{contract: _Ethereum.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ethereum *EthereumFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *EthereumRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Ethereum.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumRoleAdminChanged)
				if err := _Ethereum.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Ethereum *EthereumFilterer) ParseRoleAdminChanged(log types.Log) (*EthereumRoleAdminChanged, error) {
	event := new(EthereumRoleAdminChanged)
	if err := _Ethereum.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Ethereum contract.
type EthereumRoleGrantedIterator struct {
	Event *EthereumRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthereumRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthereumRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthereumRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumRoleGranted represents a RoleGranted event raised by the Ethereum contract.
type EthereumRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthereumRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ethereum.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthereumRoleGrantedIterator{contract: _Ethereum.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *EthereumRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ethereum.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumRoleGranted)
				if err := _Ethereum.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) ParseRoleGranted(log types.Log) (*EthereumRoleGranted, error) {
	event := new(EthereumRoleGranted)
	if err := _Ethereum.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Ethereum contract.
type EthereumRoleRevokedIterator struct {
	Event *EthereumRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthereumRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthereumRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthereumRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumRoleRevoked represents a RoleRevoked event raised by the Ethereum contract.
type EthereumRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthereumRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ethereum.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthereumRoleRevokedIterator{contract: _Ethereum.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *EthereumRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Ethereum.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumRoleRevoked)
				if err := _Ethereum.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Ethereum *EthereumFilterer) ParseRoleRevoked(log types.Log) (*EthereumRoleRevoked, error) {
	event := new(EthereumRoleRevoked)
	if err := _Ethereum.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
