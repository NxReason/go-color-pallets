package cmd

import (
	"fmt"
	"os"
	"regexp"
)

var Ifa *InputFileArgument
var ofa *OutputFileArgument
var flags map[string]Argument

func init() {
	flags = make(map[string]Argument)

	Ifa = &InputFileArgument{
		ArgumentData {
			flag: "-i",
			description: "File(s) to create pallets from",
			required: true,
		},
		make([]string, 0),
	}
	flags[Ifa.flag] = Ifa

	ofa = &OutputFileArgument {
		ArgumentData {
			flag: "-o",
			description: "Created pallet's file name",
			required: false,
		},
		make([]string, 0),
	}
	flags[ofa.flag] = ofa
}

func ParseArgs() []error {
	args := os.Args[1:]
	flagPattern := regexp.MustCompile(`^-[a-zA-Z]$`) // single letter only
	
	var currentArg Argument
	badFlag := false
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// check if its a flag
		match := flagPattern.FindAllString(arg, -1)
		if len(match) != 0 {
			switch (arg) {
			case Ifa.flag:
				badFlag = false
				currentArg = Ifa
			case ofa.flag:
				badFlag = false
				currentArg = ofa
			default:
				badFlag = true
				fmt.Println("unknow flag:", arg, "- values will be ignored until next flag found")
			}
			continue
		}

		// if its not a flag check if command is in valid state
		if (badFlag) { continue }
		currentArg.AddValue(arg)
	}

	// collect validation errors
	errs := make([]error, 0)
	for _, arg := range flags {
		err := arg.Validate()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

