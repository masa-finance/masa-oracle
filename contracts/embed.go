//go:generate yarn install
package contracts

import "embed"

//go:embed *
var EmbeddedContracts embed.FS
