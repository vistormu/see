package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"see/internal/builder"
	"see/internal/printer"
)

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "see",
		Short: "a better print of your file system",
		Long:  `see is a command line tool that provides a better way to visualize your file system, merging commands like ls, tree, du, and df into one.`,
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			command, err := builder.BuildCommand(args)
			if err != nil {
				exitWithError(err)
			}

			err = printer.Print(command)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	// disable arguments validation
	rootCmd.DisableFlagParsing = true

	err := rootCmd.Execute()
	if err != nil {
		exitWithError(err)
	}
}
