package app

import (
	"github.com/liamg/ariadne/bench"
	"github.com/spf13/cobra"
)

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Benchmark some known positions to fixed depths",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bench.Run(cmd.Context())
	},
}
