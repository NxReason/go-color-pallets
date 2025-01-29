package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DEFAULT_ROWS = 8
	DEFAULT_COLS = 8
)

var IMAGE_EXTENSIONS = [...]string { ".jpg", ".jpeg", ".png" }

type Mode string
const (
	GRID 		Mode = "GRID"
	PALLETE Mode = "PALLETE"
)
var Modes = map[Mode]string {
	GRID: "GRID",
	PALLETE: "PALLETE",
}

type Config struct {
	InputFiles []string

	GridRows int
	GridCols int
	gridSet  bool

	OutputWidth  int
	OutputHeight int

	Modes []string
}

func (c *Config) SetDefaults() {
	if !c.gridSet {
		c.GridRows = DEFAULT_ROWS
		c.GridCols = DEFAULT_COLS
	}

	if len(c.Modes) == 0 {
		for _, v := range Modes {
			c.Modes = append(c.Modes, v)
		}
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

	// modes
	for _, m := range c.Modes {
		if !isValidMode(m) {
			errs = append(errs, errors.New("invalid mode: " + m))
		}
	}

	return errs
}

func isValidMode(m string) bool {
	for _, v := range Modes {
		if strings.ToUpper(m) == v { return true }
	}
	return false
}

func (c *Config) addInputFiles(files []string) {
	c.InputFiles = append(c.InputFiles, files...)
}

func (c *Config) addFolders(folders []string) error {
	for _, f := range folders {
		err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			
			if !info.IsDir() && isImageFile(path) {
				c.InputFiles = append(c.InputFiles, path)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range IMAGE_EXTENSIONS {
		if ext == validExt { return true }
	}
	return false
}

func (c *Config) setModes(args []string) {
	c.Modes = make([]string, len(args))
	for i, a := range args {
		c.Modes[i] = strings.ToUpper(a)
	}
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
	rc, err := makeUniformPair(str)
	if err != nil {
		return errors.New("wrong grid format, " + err.Error())
	}

	rows, err := strconv.Atoi(rc[0])
	if err != nil { return errors.New("can't convert value: " + rc[0] + " to number of rows") }
	cols, err := strconv.Atoi(rc[1])
	if err != nil { return errors.New("can't convert value: " + rc[1] + " to number of columns") }
	c.GridRows, c.GridCols = rows, cols
	return nil
}

func (c *Config) setOutputResolution(args []string) error {
	const syntax = "acceptable syntax: [10x10] [10*10] [10 10]"
	switch len(args) {
	case 0:
		return errors.New("not enough arguments for output resolution. " + syntax)
	case 1:
		err := c.parseResolutionFromString(args[0])
		if err != nil { return err }
	case 2:
		err := c.parseResolutionFromString(args[0] + " " + args[1])
		if err != nil { return err }
	default:
		return errors.New("too many arguments for output resolution. " + syntax)
	}

	return nil
}

func (c *Config) parseResolutionFromString(str string) error {
	rc, err := makeUniformPair(str)
	if err != nil {
		return errors.New("wrong resolution format, " + err.Error())
	}

	width, err := strconv.Atoi(rc[0])
	if err != nil { return errors.New("can't convert value: " + rc[0] + " to output width") }
	height, err := strconv.Atoi(rc[1])
	if err != nil { return errors.New("can't convert value: " + rc[1] + " to output height") }
	c.OutputWidth, c.OutputHeight = width, height
	return nil
}

func makeUniformPair(str string) ([]string, error) {
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
		return nil, errors.New("when single argument provided, acceptable formats: [10x10] [10*10]")
	}
	return rc, nil
}