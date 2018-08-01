package flags

import "github.com/whitman-colm/go-discord/dat"

type Type string

const (
	Dash       Type = "-"
	DoubleDash Type = "--"
)

// A Flag is used to store a single flags data.
//
// Fields:
//  - Type: "-" or "--"
//  - Name: Name of the flag.
//      Ex: --name gabe miller --> Name = name
//  - Values: Single string of values after flag.
//      Ex: --name gabe miller --> Values = gabe miller
//
type Flag struct {
	Type  Type
	Name  string
	Value string
}

// Parse parses a message for flags.
//
// Parameters:
// - args ([]string) | A message split into []string
//
// Returns:
// - ([]*Flag) | A slice of each flag type
//
func Parse(args []string) []*Flag {
	flags := []*Flag{}
	var cur *Flag
	for _, arg := range args {
		switch {
		case len(arg) > 1 && arg[:2] == "--":
			cur = &Flag{
				Type: DoubleDash,
				Name: arg[2:],
			}
			flags = append(flags, cur)
		case arg[0] == '-':
			cur = &Flag{
				Type: Dash,
				Name: arg[1:],
			}
			flags = append(flags, cur)
		case arg[0] != '-':
			if len(flags) > 0 {
				flags[len(flags)-1].Value += arg + " "
			}
		default:
			dat.Log.Println("System recived flag that was not valid: \"" + arg + "\" .")
		}
	}
	return flags
}
