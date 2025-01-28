package cmd

import (
	"errors"
	"strconv"
	"strings"
)

const (
	DEFAULT_ROWS = 8
	DEFAULT_COLS = 8
)

type Config struct {
	InputFiles []string

	GridRows int
	GridCols int
	gridSet bool
}

func (c *Config) SetDefaults() {
	if !c.gridSet {
		c.GridRows = DEFAULT_ROWS
		c.GridCols = DEFAULT_COLS
	}
}

func (c *Config) Validate() []error {
	errs := make([]error, 0)

	// input files
	if len(c.InputFiles) < 1 {
		errs = append(errs, errors.New("not enough input files to process. syntax: -i filename.jpg [anotherfile.png]"))
	}

	// grid cols / rows
	if c.GridRows < 1 {
		errs = append(errs, errors.New("number of grid rows must be > 1. got " + strconv.Itoa(c.GridRows)))
	}
	if c.GridCols < 1 {
		errs = append(errs, errors.New("number of grid columns must be > 1. got " + strconv.Itoa(c.GridCols)))
	}

	return errs
}

func (c *Config) addInputFiles(files []string) {
	c.InputFiles = append(c.InputFiles, files...)
}

func (c *Config) setGrid(args []string) error {
	const syntax = "acceptable syntax: [10x10] [10*10] [10 10]"
	switch len(args) {
	case 0:
		return errors.New("not enough arguments for grid. " + syntax)
	case 1:
		err := c.parseGridFromString(args[0])
		if err != nil { return err }
	case 2:
		err := c.parseGridFromString(args[0] + " " + args[1])
		if err != nil { return err }
	default:
		return errors.New("too many arguments for grid. " + syntax)
	}

	c.gridSet = true
	return nil
}

func (c *Config) parseGridFromString(str string) error {
	var uniform []rune
	for _, ch := range str {
		if ch == '*' || ch == 'x' {
			uniform = append(uniform, ' ')
			continue
		}
		uniform = append(uniform, ch)
	}

	rc := strings.Split(string(uniform), " ")
	if len(rc) != 2 {
		return errors.New("wrong grid format, if single argument provided, acceptable formats: [10x10] [10*10]")
	}

	rows, err := strconv.Atoi(rc[0])
	if err != nil { return errors.New("can't convert value: " + rc[0] + " to number of rows") }
	cols, err := strconv.Atoi(rc[1])
	if err != nil { return errors.New("can't convert value: " + rc[1] + " to number of columns") }

	c.GridRows, c.GridCols = rows, cols
	return nil
}