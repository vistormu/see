package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vistormu/go-dsa/ansi"

	"see/internal/builder"
	"see/internal/printer"
)

const (
	Version = "0.0.3"
	seeName = ansi.Bold + ansi.Italic + ansi.Magenta + "see" + ansi.Reset
)

var (
	sortBy  string
	version bool
	help    bool
	filter  string
	depth   int
	nerd    bool
)

var rootCmd = &cobra.Command{
	Use:   seeName,
	Short: "a better way to visualize your file system",
	Long:  seeName + ` is the replacement of ls, tree, and cat commands with a more user-friendly output, with a focus on git repositories`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && version {
			fmt.Printf("%s version %s%s%s\n",
				seeName,
				ansi.Green+ansi.Bold,
				Version,
				ansi.Reset,
			)
			return nil
		}

		element := "."
		if len(args) != 0 {
			element = args[0]
		}

		// parse args
		parsedArgs := builder.Args{
			Element: element,
			Sort:    sortBy,
			Filter:  filter,
			Depth:   depth,
			Nerd:    nerd,
		}

		command, err := builder.BuildCommand(parsedArgs)
		if err != nil {
			return err
		}

		err = printer.Print(command, parsedArgs)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	// --sort, -s
	rootCmd.Flags().StringVarP(
		&sortBy,
		"sort",
		"s",
		"name",
		`sort order: "name", "kind", "size", "git-status", "date"`,
	)

	// --version, -v
	rootCmd.Flags().BoolVarP(
		&version,
		"version",
		"v",
		false,
		"print the version of see",
	)

	// --help, -h
	rootCmd.Flags().BoolVarP(
		&help,
		"help",
		"h",
		false,
		"print the help message",
	)

	// --filter, -f
	rootCmd.Flags().StringVarP(
		&filter,
		"filter",
		"f",
		"",
		"filter files by name, supports glob patterns",
	)

	// --depth, -d
	rootCmd.Flags().IntVarP(
		&depth,
		"depth",
		"d",
		1,
		"set the maximum depth of directories to traverse, 0 means no limit",
	)

	// --nerd, -n
	rootCmd.Flags().BoolVarP(
		&nerd,
		"nerd",
		"n",
		false,
		"show all information of the printed items",
	)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
