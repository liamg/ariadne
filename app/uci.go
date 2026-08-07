package app

import (
	"github.com/liamg/ariadne/uci"
	"github.com/spf13/cobra"
)

var uciCmd = &cobra.Command{
	Use:   "uci",
	Short: "Start a UCI protocol handler",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		u := uci.New(cmd.InOrStdin(), cmd.OutOrStdout())
		return u.Run(cmd.Context())
	},
}
