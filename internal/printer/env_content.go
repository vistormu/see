package printer

import (
	"fmt"

	"see/internal/builder"

	"github.com/vistormu/go-dsa/ansi"
)

func printEnvVariable(envVar *builder.EnvVariable, args builder.Args) error {
	name := fmt.Sprintf("%s$%s%s", ansi.Bold+ansi.Green, envVar.Name, ansi.Reset)
	value := envVar.Value
	if args.Filter != "" {
		value = filterLines(value, args.Filter)
	}

	fmt.Printf("%s\n\n%s\n\n", name, value)

	if args.Copy {
		if err := copyFn(value); err != nil {
			return err
		}
	}

	return nil
}
