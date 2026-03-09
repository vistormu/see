package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vistormu/go-dsa/ansi"

	"see/internal/builder"
	"see/internal/printer"
)

const (
	Version = "0.0.5"
	seeName = ansi.Bold + ansi.Italic + ansi.Magenta + "see" + ansi.Reset
)

var (
	sortBy     string
	version    bool
	help       bool
	filter     string
	depth      int
	nerd       bool
	head       int
	tail       int
	copyOutput bool
)

var rootCmd = &cobra.Command{
	Use:           seeName,
	Short:         "a better way to visualize your file system",
	Long:          seeName + ` is the replacement of ls, tree, and cat commands with a more user-friendly output, with a focus on git repositories`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
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

		if head >= 0 && tail >= 0 {
			return fmt.Errorf("flags --head and --tail cannot be used together")
		}
		if head < -1 || tail < -1 {
			return fmt.Errorf("flags --head and --tail must be positive values")
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
			Head:    head,
			Tail:    tail,
			Copy:    copyOutput,
		}

		command, err := builder.BuildCommand(parsedArgs)
		if err != nil {
			return err
		}

		if _, isDir := command.(*builder.Directory); isDir && (head >= 0 || tail >= 0) {
			return fmt.Errorf("flags --head and --tail can only be used when showing file or env content")
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

	// --head, -H
	rootCmd.Flags().IntVarP(
		&head,
		"head",
		"H",
		-1,
		"show only the first N lines when printing file content",
	)

	// --tail, -t
	rootCmd.Flags().IntVarP(
		&tail,
		"tail",
		"t",
		-1,
		"show only the last N lines when printing file content",
	)

	// --copy, -c
	rootCmd.Flags().BoolVarP(
		&copyOutput,
		"copy",
		"c",
		false,
		"copy rendered output content to your clipboard",
	)
}

func normalizeOptionalIntFlags(args []string) []string {
	normalized := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg != "--head" && arg != "-H" && arg != "--tail" && arg != "-t" {
			normalized = append(normalized, arg)
			continue
		}

		defaultValue := "10"
		if i == len(args)-1 {
			normalized = append(normalized, arg, defaultValue)
			continue
		}

		next := args[i+1]
		if strings.HasPrefix(next, "-") {
			normalized = append(normalized, arg, defaultValue)
			continue
		}

		if _, err := strconv.Atoi(next); err == nil {
			normalized = append(normalized, arg, next)
			i++
			continue
		}

		normalized = append(normalized, arg, defaultValue)
	}

	return normalized
}

func Execute() {
	rootCmd.SetArgs(normalizeOptionalIntFlags(os.Args[1:]))
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
