package cmd

import (
	"errors"
	"os"
	"regexp"
	"sort"
)

const flagRE = `^-[a-zA-Z]$` // single letter only

func ParseArgs() (Config, []error) {
	args := os.Args[1:]
	flagsPos := FindAllFlags(args)

	return MakeConfig(args, flagsPos)
}

func FindAllFlags(args []string) (flagsPos map[int]string) {
	flagsPos = make(map[int]string)
	flagPattern := regexp.MustCompile(flagRE)
	
	for i, arg := range args {
		match := flagPattern.FindString(arg)
		if len(match) > 0 {
			flagsPos[i] = match
		}
	}

	return
}

func MakeConfig(args []string, flagsPos map[int]string) (Config, []error) {
	// ensure flags order
	positions := make([]int, 0, len(flagsPos))
	for pos := range flagsPos {
		positions = append(positions, pos)
	}
	sort.Ints(positions)

	// make config
	config := Config{}
	var err error
	errs := make([]error, 0)
	for i, pos := range positions {
		err = nil
		firstArg := pos + 1
		var lastArg int
		if i == len(positions) - 1 {
			lastArg = len(args)
		} else {
			lastArg = positions[i + 1]
		}
		argSlice := args[firstArg:lastArg]
		switch flag := flagsPos[pos]; flag {
		case "-i":
			config.addInputFiles(argSlice)
		case "-g":
			err = config.setGrid(argSlice)
		case "-r":
			err = config.setOutputResolution(argSlice)
		case "-f":
			err = config.addFolders(argSlice)
		case "-m":
			config.setModes(argSlice)
		default:
			err = errors.New("Unknown flag: " + flag + " (skipped)")
		}
		if err != nil { errs = append(errs, err) }
	}

	return config, errs
}

