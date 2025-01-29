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
	errs := make([]error, 0)
	for i, pos := range positions {
		firstArg := pos + 1
		var lastArg int
		if i == len(positions) - 1 {
			lastArg = len(args)
		} else {
			lastArg = positions[i + 1]
		}
		switch flag := flagsPos[pos]; flag {
		case "-i":
			config.addInputFiles(args[firstArg:lastArg])
		case "-g":
			err := config.setGrid(args[firstArg:lastArg])
			if err != nil { errs = append(errs, err) }
		case "-r":
			err := config.setOutputResolution(args[firstArg:lastArg])
			if err != nil { errs = append(errs, err) }
		default:
			errs = append(errs, errors.New("Unknown flag: " + flag + " (skipped)"))
		}
	}

	return config, errs
}

