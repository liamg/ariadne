package app

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Run() error {
	rootCmd := &cobra.Command{
		Use:   "ariadne",
		Short: "A UCI-compatible chess engine",
		Args:  cobra.NoArgs,
	}

	rootCmd.AddCommand(uciCmd)
	rootCmd.AddCommand(benchCmd)

	// default to uci if no subcommand given
	cmd, _, err := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{uciCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}
