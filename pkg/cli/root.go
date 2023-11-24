package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bootnodes = "/ip4/137.66.11.250/udp/4001/quic-v1/p2p/16Uiu2HAmJtYy4A8pzChDQQLrPsu1SQU5apzCftCVaAjFk539CLc9,/ip4/168.220.95.86/udp/4001/quic-v1/p2p/16Uiu2HAmHgacSbefXenfMnoDNqv9sDjzGCuCgZa1zoJCC7CBXNka,/ip4/35.224.231.145/udp/4001/quic-v1/p2p/16Uiu2HAm47nBiewWLLzCREtY8vwPQtr5jTqyrEoUo6WnngwhsQuR,/ip4/104.198.43.138/udp/4001/quic-v1/p2p/16Uiu2HAkxiP8jjdHQWeCxTr7pD6BvoPkS8Z1skjCy9vdSRMACDcc"

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Masa node",
	Long:  `Start the Masa node with predefined bootnodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting the Masa node with bootnodes: %s\n", bootnodes)
		// Add your code to start the node here
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
