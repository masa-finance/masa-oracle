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

// PackageNameMetaData contains all meta data concerning the PackageName contract.
var PackageNameMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractrVOTEToken\",\"name\":\"_rVoteToken\",\"type\":\"address\"},{\"internalType\":\"contracttMASAToken\",\"name\":\"_tMasaToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakingContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPDATER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rVoteToken\",\"outputs\":[{\"internalType\":\"contractrVOTEToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakingContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tMasaToken\",\"outputs\":[{\"internalType\":\"contracttMASAToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"users\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"reputation_score\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reputation_votes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"vote_count\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"reputation_score\",\"type\":\"string\"}],\"name\":\"addUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"reputation_score\",\"type\":\"string\"}],\"name\":\"updateUserScore\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"voteAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reputationVote\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"}],\"name\":\"getUserInfo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200162838038062001628833981016040819052620000349162000164565b600180546001600160a01b038087166001600160a01b0319928316179092556002805486841690831617905560048054928516929091169190911790556200007e600082620000b4565b620000aa7f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab82620000b4565b50505050620001e5565b620000c08282620000c4565b5050565b6000828152602081815260408083206001600160a01b038516845290915290205460ff16620000c0576000828152602081815260408083206001600160a01b03851684529091529020805460ff19166001179055620001203390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b600080600080608085870312156200017b57600080fd5b84516200018881620001cc565b60208601519094506200019b81620001cc565b6040860151909350620001ae81620001cc565b6060860151909250620001c181620001cc565b939692955090935050565b6001600160a01b0381168114620001e257600080fd5b50565b61143380620001f56000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063689e1c031161009757806391d148541161006657806391d148541461025a578063a217fddf1461026d578063bf246ea214610275578063d547741f1461028857600080fd5b8063689e1c03146101fe5780637c9b7fdd1461022157806386a08a251461023457806386f524661461024757600080fd5b80632f2ff15d116100d35780632f2ff15d146101865780633535f48b1461019957806336568abe146101c457806347e63380146101d757600080fd5b806301ffc9a714610105578063079eaf341461012d5780630891d35814610142578063248a9ca314610155575b600080fd5b610118610113366004611092565b61029b565b60405190151581526020015b60405180910390f35b61014061013b3660046110f9565b6102d2565b005b61014061015036600461115d565b61038e565b61017861016336600461103d565b60009081526020819052604090206001015490565b604051908152602001610124565b610140610194366004611056565b6107f1565b6004546101ac906001600160a01b031681565b6040516001600160a01b039091168152602001610124565b6101406101d2366004611056565b61081b565b6101787f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab81565b61021161020c3660046110bc565b610899565b6040516101249493929190611294565b61021161022f3660046110bc565b6109dc565b6002546101ac906001600160a01b031681565b6101406102553660046110f9565b610b3e565b610118610268366004611056565b610ba3565b610178600081565b6001546101ac906001600160a01b031681565b610140610296366004611056565b610bcc565b60006001600160e01b03198216637965db0b60e01b14806102cc57506301ffc9a760e01b6001600160e01b03198316145b92915050565b7f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab6102fc81610bf1565b604051806080016040528084815260200183815260200160008152602001600081525060038460405161032f91906111f0565b90815260200160405180910390206000820151816000019080519060200190610359929190610ef5565b5060208281015180516103729260018501920190610ef5565b5060408201516002820155606090910151600390910155505050565b81670de0b6b3a7640000146103ea5760405162461bcd60e51b815260206004820152601b60248201527f566f746520616d6f756e74206d757374206265203120746f6b656e000000000060448201526064015b60405180910390fd5b670de0b6b3a7640000811015801561040a5750678ac7230489e800008111155b61046e5760405162461bcd60e51b815260206004820152602f60248201527f52657075746174696f6e20766f7465206d757374206265206265747765656e2060448201526e3120616e6420313020746f6b656e7360881b60648201526084016103e1565b6002546040516370a0823160e01b815230600482015283916001600160a01b0316906370a082319060240160206040518083038186803b1580156104b157600080fd5b505afa1580156104c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e991906111ab565b10156105475760405162461bcd60e51b815260206004820152602760248201527f4e6f7420656e6f75676820744d41534120746f6b656e7320696e2074686520636044820152661bdb9d1c9858dd60ca1b60648201526084016103e1565b6001546040516323b872dd60e01b8152336004820152306024820152604481018490526001600160a01b03909116906323b872dd90606401602060405180830381600087803b15801561059957600080fd5b505af11580156105ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d1919061101b565b5060006105df6002846112e5565b905060006105ed8285611326565b600154604051632770a7eb60e21b8152306004820152602481018590529192506001600160a01b031690639dc29fac90604401600060405180830381600087803b15801561063a57600080fd5b505af115801561064e573d6000803e3d6000fd5b50506001546004805460405163a9059cbb60e01b81526001600160a01b039182169281019290925260248201869052909116925063a9059cbb9150604401602060405180830381600087803b1580156106a657600080fd5b505af11580156106ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106de919061101b565b506002546001600160a01b031663a9059cbb336106fc87600a611307565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401602060405180830381600087803b15801561074257600080fd5b505af1158015610756573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061077a919061101b565b508260038660405161078c91906111f0565b908152602001604051809103902060020160008282546107ac91906112cd565b9250508190555060016003866040516107c591906111f0565b908152602001604051809103902060030160008282546107e591906112cd565b90915550505050505050565b60008281526020819052604090206001015461080c81610bf1565b6108168383610bfe565b505050565b6001600160a01b038116331461088b5760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b60648201526084016103e1565b6108958282610c82565b5050565b80516020818301810180516003825292820191909301209152805481906108bf90611380565b80601f01602080910402602001604051908101604052809291908181526020018280546108eb90611380565b80156109385780601f1061090d57610100808354040283529160200191610938565b820191906000526020600020905b81548152906001019060200180831161091b57829003601f168201915b50505050509080600101805461094d90611380565b80601f016020809104026020016040519081016040528092919081815260200182805461097990611380565b80156109c65780601f1061099b576101008083540402835291602001916109c6565b820191906000526020600020905b8154815290600101906020018083116109a957829003601f168201915b5050505050908060020154908060030154905084565b60608060008060006003866040516109f491906111f0565b90815260200160405180910390209050806000018160010182600201548360030154838054610a2290611380565b80601f0160208091040260200160405190810160405280929190818152602001828054610a4e90611380565b8015610a9b5780601f10610a7057610100808354040283529160200191610a9b565b820191906000526020600020905b815481529060010190602001808311610a7e57829003601f168201915b50505050509350828054610aae90611380565b80601f0160208091040260200160405190810160405280929190818152602001828054610ada90611380565b8015610b275780601f10610afc57610100808354040283529160200191610b27565b820191906000526020600020905b815481529060010190602001808311610b0a57829003601f168201915b505050505092509450945094509450509193509193565b7f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab610b6881610bf1565b81600384604051610b7991906111f0565b90815260200160405180910390206001019080519060200190610b9d929190610ef5565b50505050565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b600082815260208190526040902060010154610be781610bf1565b6108168383610c82565b610bfb8133610ce7565b50565b610c088282610ba3565b610895576000828152602081815260408083206001600160a01b03851684529091529020805460ff19166001179055610c3e3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b610c8c8282610ba3565b15610895576000828152602081815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b610cf18282610ba3565b61089557610cfe81610d40565b610d09836020610d52565b604051602001610d1a92919061120c565b60408051601f198184030181529082905262461bcd60e51b82526103e191600401611281565b60606102cc6001600160a01b03831660145b60606000610d61836002611307565b610d6c9060026112cd565b67ffffffffffffffff811115610d8457610d846113e7565b6040519080825280601f01601f191660200182016040528015610dae576020820181803683370190505b509050600360fc1b81600081518110610dc957610dc96113d1565b60200101906001600160f81b031916908160001a905350600f60fb1b81600181518110610df857610df86113d1565b60200101906001600160f81b031916908160001a9053506000610e1c846002611307565b610e279060016112cd565b90505b6001811115610e9f576f181899199a1a9b1b9c1cb0b131b232b360811b85600f1660108110610e5b57610e5b6113d1565b1a60f81b828281518110610e7157610e716113d1565b60200101906001600160f81b031916908160001a90535060049490941c93610e9881611369565b9050610e2a565b508315610eee5760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016103e1565b9392505050565b828054610f0190611380565b90600052602060002090601f016020900481019282610f235760008555610f69565b82601f10610f3c57805160ff1916838001178555610f69565b82800160010185558215610f69579182015b82811115610f69578251825591602001919060010190610f4e565b50610f75929150610f79565b5090565b5b80821115610f755760008155600101610f7a565b600082601f830112610f9f57600080fd5b813567ffffffffffffffff80821115610fba57610fba6113e7565b604051601f8301601f19908116603f01168101908282118183101715610fe257610fe26113e7565b81604052838152866020858801011115610ffb57600080fd5b836020870160208301376000602085830101528094505050505092915050565b60006020828403121561102d57600080fd5b81518015158114610eee57600080fd5b60006020828403121561104f57600080fd5b5035919050565b6000806040838503121561106957600080fd5b8235915060208301356001600160a01b038116811461108757600080fd5b809150509250929050565b6000602082840312156110a457600080fd5b81356001600160e01b031981168114610eee57600080fd5b6000602082840312156110ce57600080fd5b813567ffffffffffffffff8111156110e557600080fd5b6110f184828501610f8e565b949350505050565b6000806040838503121561110c57600080fd5b823567ffffffffffffffff8082111561112457600080fd5b61113086838701610f8e565b9350602085013591508082111561114657600080fd5b5061115385828601610f8e565b9150509250929050565b60008060006060848603121561117257600080fd5b833567ffffffffffffffff81111561118957600080fd5b61119586828701610f8e565b9660208601359650604090950135949350505050565b6000602082840312156111bd57600080fd5b5051919050565b600081518084526111dc81602086016020860161133d565b601f01601f19169290920160200192915050565b6000825161120281846020870161133d565b9190910192915050565b7f416363657373436f6e74726f6c3a206163636f756e742000000000000000000081526000835161124481601785016020880161133d565b7001034b99036b4b9b9b4b733903937b6329607d1b601791840191820152835161127581602884016020880161133d565b01602801949350505050565b602081526000610eee60208301846111c4565b6080815260006112a760808301876111c4565b82810360208401526112b981876111c4565b604084019590955250506060015292915050565b600082198211156112e0576112e06113bb565b500190565b60008261130257634e487b7160e01b600052601260045260246000fd5b500490565b6000816000190483118215151615611321576113216113bb565b500290565b600082821015611338576113386113bb565b500390565b60005b83811015611358578181015183820152602001611340565b83811115610b9d5750506000910152565b600081611378576113786113bb565b506000190190565b600181811c9082168061139457607f821691505b602082108114156113b557634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea264697066735822122099173f74c898035ba9060c2c986da6715fe0fba762fe000f7d67b4bcc3049c7864736f6c63430008070033",
}

// PackageNameABI is the input ABI used to generate the binding from.
// Deprecated: Use PackageNameMetaData.ABI instead.
var PackageNameABI = PackageNameMetaData.ABI

// PackageNameBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PackageNameMetaData.Bin instead.
var PackageNameBin = PackageNameMetaData.Bin

// DeployPackageName deploys a new Ethereum contract, binding an instance of PackageName to it.
func DeployPackageName(auth *bind.TransactOpts, backend bind.ContractBackend, _rVoteToken common.Address, _tMasaToken common.Address, _stakingContractAddress common.Address, admin common.Address) (common.Address, *types.Transaction, *PackageName, error) {
	parsed, err := PackageNameMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PackageNameBin), backend, _rVoteToken, _tMasaToken, _stakingContractAddress, admin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PackageName{PackageNameCaller: PackageNameCaller{contract: contract}, PackageNameTransactor: PackageNameTransactor{contract: contract}, PackageNameFilterer: PackageNameFilterer{contract: contract}}, nil
}

// PackageName is an auto generated Go binding around an Ethereum contract.
type PackageName struct {
	PackageNameCaller     // Read-only binding to the contract
	PackageNameTransactor // Write-only binding to the contract
	PackageNameFilterer   // Log filterer for contract events
}

// PackageNameCaller is an auto generated read-only Go binding around an Ethereum contract.
type PackageNameCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PackageNameTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PackageNameTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PackageNameFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PackageNameFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PackageNameSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PackageNameSession struct {
	Contract     *PackageName      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PackageNameCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PackageNameCallerSession struct {
	Contract *PackageNameCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PackageNameTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PackageNameTransactorSession struct {
	Contract     *PackageNameTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PackageNameRaw is an auto generated low-level Go binding around an Ethereum contract.
type PackageNameRaw struct {
	Contract *PackageName // Generic contract binding to access the raw methods on
}

// PackageNameCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PackageNameCallerRaw struct {
	Contract *PackageNameCaller // Generic read-only contract binding to access the raw methods on
}

// PackageNameTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PackageNameTransactorRaw struct {
	Contract *PackageNameTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPackageName creates a new instance of PackageName, bound to a specific deployed contract.
func NewPackageName(address common.Address, backend bind.ContractBackend) (*PackageName, error) {
	contract, err := bindPackageName(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PackageName{PackageNameCaller: PackageNameCaller{contract: contract}, PackageNameTransactor: PackageNameTransactor{contract: contract}, PackageNameFilterer: PackageNameFilterer{contract: contract}}, nil
}

// NewPackageNameCaller creates a new read-only instance of PackageName, bound to a specific deployed contract.
func NewPackageNameCaller(address common.Address, caller bind.ContractCaller) (*PackageNameCaller, error) {
	contract, err := bindPackageName(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PackageNameCaller{contract: contract}, nil
}

// NewPackageNameTransactor creates a new write-only instance of PackageName, bound to a specific deployed contract.
func NewPackageNameTransactor(address common.Address, transactor bind.ContractTransactor) (*PackageNameTransactor, error) {
	contract, err := bindPackageName(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PackageNameTransactor{contract: contract}, nil
}

// NewPackageNameFilterer creates a new log filterer instance of PackageName, bound to a specific deployed contract.
func NewPackageNameFilterer(address common.Address, filterer bind.ContractFilterer) (*PackageNameFilterer, error) {
	contract, err := bindPackageName(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PackageNameFilterer{contract: contract}, nil
}

// bindPackageName binds a generic wrapper to an already deployed contract.
func bindPackageName(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PackageNameMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PackageName *PackageNameRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PackageName.Contract.PackageNameCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PackageName *PackageNameRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PackageName.Contract.PackageNameTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PackageName *PackageNameRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PackageName.Contract.PackageNameTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PackageName *PackageNameCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PackageName.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PackageName *PackageNameTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PackageName.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PackageName *PackageNameTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PackageName.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PackageName *PackageNameCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PackageName *PackageNameSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _PackageName.Contract.DEFAULTADMINROLE(&_PackageName.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PackageName *PackageNameCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _PackageName.Contract.DEFAULTADMINROLE(&_PackageName.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_PackageName *PackageNameCaller) UPDATERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "UPDATER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_PackageName *PackageNameSession) UPDATERROLE() ([32]byte, error) {
	return _PackageName.Contract.UPDATERROLE(&_PackageName.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_PackageName *PackageNameCallerSession) UPDATERROLE() ([32]byte, error) {
	return _PackageName.Contract.UPDATERROLE(&_PackageName.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PackageName *PackageNameCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PackageName *PackageNameSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _PackageName.Contract.GetRoleAdmin(&_PackageName.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PackageName *PackageNameCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _PackageName.Contract.GetRoleAdmin(&_PackageName.CallOpts, role)
}

// GetUserInfo is a free data retrieval call binding the contract method 0x7c9b7fdd.
//
// Solidity: function getUserInfo(string id) view returns(string, string, uint256, uint256)
func (_PackageName *PackageNameCaller) GetUserInfo(opts *bind.CallOpts, id string) (string, string, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "getUserInfo", id)

	if err != nil {
		return *new(string), *new(string), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, err

}

// GetUserInfo is a free data retrieval call binding the contract method 0x7c9b7fdd.
//
// Solidity: function getUserInfo(string id) view returns(string, string, uint256, uint256)
func (_PackageName *PackageNameSession) GetUserInfo(id string) (string, string, *big.Int, *big.Int, error) {
	return _PackageName.Contract.GetUserInfo(&_PackageName.CallOpts, id)
}

// GetUserInfo is a free data retrieval call binding the contract method 0x7c9b7fdd.
//
// Solidity: function getUserInfo(string id) view returns(string, string, uint256, uint256)
func (_PackageName *PackageNameCallerSession) GetUserInfo(id string) (string, string, *big.Int, *big.Int, error) {
	return _PackageName.Contract.GetUserInfo(&_PackageName.CallOpts, id)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PackageName *PackageNameCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PackageName *PackageNameSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _PackageName.Contract.HasRole(&_PackageName.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PackageName *PackageNameCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _PackageName.Contract.HasRole(&_PackageName.CallOpts, role, account)
}

// RVoteToken is a free data retrieval call binding the contract method 0xbf246ea2.
//
// Solidity: function rVoteToken() view returns(address)
func (_PackageName *PackageNameCaller) RVoteToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "rVoteToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RVoteToken is a free data retrieval call binding the contract method 0xbf246ea2.
//
// Solidity: function rVoteToken() view returns(address)
func (_PackageName *PackageNameSession) RVoteToken() (common.Address, error) {
	return _PackageName.Contract.RVoteToken(&_PackageName.CallOpts)
}

// RVoteToken is a free data retrieval call binding the contract method 0xbf246ea2.
//
// Solidity: function rVoteToken() view returns(address)
func (_PackageName *PackageNameCallerSession) RVoteToken() (common.Address, error) {
	return _PackageName.Contract.RVoteToken(&_PackageName.CallOpts)
}

// StakingContractAddress is a free data retrieval call binding the contract method 0x3535f48b.
//
// Solidity: function stakingContractAddress() view returns(address)
func (_PackageName *PackageNameCaller) StakingContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "stakingContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingContractAddress is a free data retrieval call binding the contract method 0x3535f48b.
//
// Solidity: function stakingContractAddress() view returns(address)
func (_PackageName *PackageNameSession) StakingContractAddress() (common.Address, error) {
	return _PackageName.Contract.StakingContractAddress(&_PackageName.CallOpts)
}

// StakingContractAddress is a free data retrieval call binding the contract method 0x3535f48b.
//
// Solidity: function stakingContractAddress() view returns(address)
func (_PackageName *PackageNameCallerSession) StakingContractAddress() (common.Address, error) {
	return _PackageName.Contract.StakingContractAddress(&_PackageName.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PackageName *PackageNameCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PackageName *PackageNameSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PackageName.Contract.SupportsInterface(&_PackageName.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PackageName *PackageNameCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PackageName.Contract.SupportsInterface(&_PackageName.CallOpts, interfaceId)
}

// TMasaToken is a free data retrieval call binding the contract method 0x86a08a25.
//
// Solidity: function tMasaToken() view returns(address)
func (_PackageName *PackageNameCaller) TMasaToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "tMasaToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TMasaToken is a free data retrieval call binding the contract method 0x86a08a25.
//
// Solidity: function tMasaToken() view returns(address)
func (_PackageName *PackageNameSession) TMasaToken() (common.Address, error) {
	return _PackageName.Contract.TMasaToken(&_PackageName.CallOpts)
}

// TMasaToken is a free data retrieval call binding the contract method 0x86a08a25.
//
// Solidity: function tMasaToken() view returns(address)
func (_PackageName *PackageNameCallerSession) TMasaToken() (common.Address, error) {
	return _PackageName.Contract.TMasaToken(&_PackageName.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0x689e1c03.
//
// Solidity: function users(string ) view returns(string id, string reputation_score, uint256 reputation_votes, uint256 vote_count)
func (_PackageName *PackageNameCaller) Users(opts *bind.CallOpts, arg0 string) (struct {
	Id              string
	ReputationScore string
	ReputationVotes *big.Int
	VoteCount       *big.Int
}, error) {
	var out []interface{}
	err := _PackageName.contract.Call(opts, &out, "users", arg0)

	outstruct := new(struct {
		Id              string
		ReputationScore string
		ReputationVotes *big.Int
		VoteCount       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.ReputationScore = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.ReputationVotes = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.VoteCount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Users is a free data retrieval call binding the contract method 0x689e1c03.
//
// Solidity: function users(string ) view returns(string id, string reputation_score, uint256 reputation_votes, uint256 vote_count)
func (_PackageName *PackageNameSession) Users(arg0 string) (struct {
	Id              string
	ReputationScore string
	ReputationVotes *big.Int
	VoteCount       *big.Int
}, error) {
	return _PackageName.Contract.Users(&_PackageName.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0x689e1c03.
//
// Solidity: function users(string ) view returns(string id, string reputation_score, uint256 reputation_votes, uint256 vote_count)
func (_PackageName *PackageNameCallerSession) Users(arg0 string) (struct {
	Id              string
	ReputationScore string
	ReputationVotes *big.Int
	VoteCount       *big.Int
}, error) {
	return _PackageName.Contract.Users(&_PackageName.CallOpts, arg0)
}

// AddUser is a paid mutator transaction binding the contract method 0x079eaf34.
//
// Solidity: function addUser(string id, string reputation_score) returns()
func (_PackageName *PackageNameTransactor) AddUser(opts *bind.TransactOpts, id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "addUser", id, reputation_score)
}

// AddUser is a paid mutator transaction binding the contract method 0x079eaf34.
//
// Solidity: function addUser(string id, string reputation_score) returns()
func (_PackageName *PackageNameSession) AddUser(id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.Contract.AddUser(&_PackageName.TransactOpts, id, reputation_score)
}

// AddUser is a paid mutator transaction binding the contract method 0x079eaf34.
//
// Solidity: function addUser(string id, string reputation_score) returns()
func (_PackageName *PackageNameTransactorSession) AddUser(id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.Contract.AddUser(&_PackageName.TransactOpts, id, reputation_score)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.GrantRole(&_PackageName.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.GrantRole(&_PackageName.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.RenounceRole(&_PackageName.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.RenounceRole(&_PackageName.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.RevokeRole(&_PackageName.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PackageName *PackageNameTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PackageName.Contract.RevokeRole(&_PackageName.TransactOpts, role, account)
}

// UpdateUserScore is a paid mutator transaction binding the contract method 0x86f52466.
//
// Solidity: function updateUserScore(string id, string reputation_score) returns()
func (_PackageName *PackageNameTransactor) UpdateUserScore(opts *bind.TransactOpts, id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "updateUserScore", id, reputation_score)
}

// UpdateUserScore is a paid mutator transaction binding the contract method 0x86f52466.
//
// Solidity: function updateUserScore(string id, string reputation_score) returns()
func (_PackageName *PackageNameSession) UpdateUserScore(id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.Contract.UpdateUserScore(&_PackageName.TransactOpts, id, reputation_score)
}

// UpdateUserScore is a paid mutator transaction binding the contract method 0x86f52466.
//
// Solidity: function updateUserScore(string id, string reputation_score) returns()
func (_PackageName *PackageNameTransactorSession) UpdateUserScore(id string, reputation_score string) (*types.Transaction, error) {
	return _PackageName.Contract.UpdateUserScore(&_PackageName.TransactOpts, id, reputation_score)
}

// Vote is a paid mutator transaction binding the contract method 0x0891d358.
//
// Solidity: function vote(string id, uint256 voteAmount, uint256 reputationVote) returns()
func (_PackageName *PackageNameTransactor) Vote(opts *bind.TransactOpts, id string, voteAmount *big.Int, reputationVote *big.Int) (*types.Transaction, error) {
	return _PackageName.contract.Transact(opts, "vote", id, voteAmount, reputationVote)
}

// Vote is a paid mutator transaction binding the contract method 0x0891d358.
//
// Solidity: function vote(string id, uint256 voteAmount, uint256 reputationVote) returns()
func (_PackageName *PackageNameSession) Vote(id string, voteAmount *big.Int, reputationVote *big.Int) (*types.Transaction, error) {
	return _PackageName.Contract.Vote(&_PackageName.TransactOpts, id, voteAmount, reputationVote)
}

// Vote is a paid mutator transaction binding the contract method 0x0891d358.
//
// Solidity: function vote(string id, uint256 voteAmount, uint256 reputationVote) returns()
func (_PackageName *PackageNameTransactorSession) Vote(id string, voteAmount *big.Int, reputationVote *big.Int) (*types.Transaction, error) {
	return _PackageName.Contract.Vote(&_PackageName.TransactOpts, id, voteAmount, reputationVote)
}

// PackageNameRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the PackageName contract.
type PackageNameRoleAdminChangedIterator struct {
	Event *PackageNameRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *PackageNameRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PackageNameRoleAdminChanged)
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
		it.Event = new(PackageNameRoleAdminChanged)
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
func (it *PackageNameRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PackageNameRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PackageNameRoleAdminChanged represents a RoleAdminChanged event raised by the PackageName contract.
type PackageNameRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_PackageName *PackageNameFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*PackageNameRoleAdminChangedIterator, error) {

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

	logs, sub, err := _PackageName.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &PackageNameRoleAdminChangedIterator{contract: _PackageName.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_PackageName *PackageNameFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *PackageNameRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _PackageName.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PackageNameRoleAdminChanged)
				if err := _PackageName.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_PackageName *PackageNameFilterer) ParseRoleAdminChanged(log types.Log) (*PackageNameRoleAdminChanged, error) {
	event := new(PackageNameRoleAdminChanged)
	if err := _PackageName.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PackageNameRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the PackageName contract.
type PackageNameRoleGrantedIterator struct {
	Event *PackageNameRoleGranted // Event containing the contract specifics and raw log

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
func (it *PackageNameRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PackageNameRoleGranted)
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
		it.Event = new(PackageNameRoleGranted)
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
func (it *PackageNameRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PackageNameRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PackageNameRoleGranted represents a RoleGranted event raised by the PackageName contract.
type PackageNameRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_PackageName *PackageNameFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PackageNameRoleGrantedIterator, error) {

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

	logs, sub, err := _PackageName.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PackageNameRoleGrantedIterator{contract: _PackageName.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_PackageName *PackageNameFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *PackageNameRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _PackageName.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PackageNameRoleGranted)
				if err := _PackageName.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_PackageName *PackageNameFilterer) ParseRoleGranted(log types.Log) (*PackageNameRoleGranted, error) {
	event := new(PackageNameRoleGranted)
	if err := _PackageName.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PackageNameRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the PackageName contract.
type PackageNameRoleRevokedIterator struct {
	Event *PackageNameRoleRevoked // Event containing the contract specifics and raw log

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
func (it *PackageNameRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PackageNameRoleRevoked)
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
		it.Event = new(PackageNameRoleRevoked)
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
func (it *PackageNameRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PackageNameRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PackageNameRoleRevoked represents a RoleRevoked event raised by the PackageName contract.
type PackageNameRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_PackageName *PackageNameFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PackageNameRoleRevokedIterator, error) {

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

	logs, sub, err := _PackageName.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PackageNameRoleRevokedIterator{contract: _PackageName.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_PackageName *PackageNameFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *PackageNameRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _PackageName.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PackageNameRoleRevoked)
				if err := _PackageName.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_PackageName *PackageNameFilterer) ParseRoleRevoked(log types.Log) (*PackageNameRoleRevoked, error) {
	event := new(PackageNameRoleRevoked)
	if err := _PackageName.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
