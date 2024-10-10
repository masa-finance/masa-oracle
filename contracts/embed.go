//go:generate yarn install
package contracts

import "embed"

//go:embed node_modules/@masa-finance/masa-contracts-oracle/addresses.json
//go:embed node_modules/@masa-finance/masa-token/addresses.json
var EmbeddedContracts embed.FS
