package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGrid_WrongArgCount(t *testing.T) {
	assert := assert.New(t)
	config := Config{}
	
	args := []string{}
	err := config.setGrid(args)
	assert.ErrorContains(err, "not enough arguments")

	args = []string{ "25", "18", "10" }
	err = config.setGrid(args)
	assert.ErrorContains(err, "too many arguments")
}

func TestSetGrid_SingleArg(t *testing.T) {
	assert := assert.New(t)
	config := Config{}
	args := []string{"10x5"}

	config.setGrid(args)

	assert.Equal(10, config.GridRows)
	assert.Equal(5, config.GridCols)
}

func TestSetGrid_TwoArgs(t *testing.T) {
	assert := assert.New(t)
	config := Config{}
	args := []string{"25", "67"}

	config.setGrid(args)

	assert.Equal(25, config.GridRows)
	assert.Equal(67, config.GridCols)
}

func TestParseGridFromString_ValidFormat(t *testing.T) {
	assert := assert.New(t)
	config := Config{}

	str := "10*5"
	err := config.parseGridFromString(str)
	want := Config { GridRows: 10, GridCols: 5 }

	assert.Equal(nil, err)
	assert.Equal(want, config)

	str = "8x7"
	err = config.parseGridFromString(str)
	want = Config { GridRows: 8, GridCols: 7 }
	
	assert.Equal(nil, err)
	assert.Equal(want, config)
}

func TestParseGridFromString_InvalidFormat(t *testing.T) {
	assert := assert.New(t)
	config := Config{}

	str := "15_2"
	err := config.parseGridFromString(str)
	
	assert.ErrorContains(err, "wrong grid format")

	str = "1*5*7"
	err = config.parseGridFromString(str)

	assert.ErrorContains(err, "wrong grid format")
}

func TestParseGridFromString_InvalidValues(t *testing.T) {
	assert := assert.New(t)
	config := Config{}

	str := "NaN*5"
	err := config.parseGridFromString(str)

	assert.ErrorContains(err, "to number of rows")
	assert.ErrorContains(err, "NaN")

	str = "10x3.5"
	err = config.parseGridFromString(str)

	assert.ErrorContains(err, "to number of columns")
	assert.ErrorContains(err, "3.5")
}

func TestSetDefaults_Grid(t *testing.T) {
	config := Config{}
	config.SetDefaults()

	assert.Equal(t, DEFAULT_ROWS, config.GridRows)
	assert.Equal(t, DEFAULT_COLS, config.GridCols)
}
func TestSetDefaults_GridSkip(t *testing.T) {
	config := Config{ GridRows: -1, GridCols: -1, gridSet: true }
	config.SetDefaults()

	assert.Equal(t, -1, config.GridRows)
	assert.Equal(t, -1, config.GridCols)
}

func TestValidate_ValidState(t *testing.T) {
	config := Config {
		InputFiles: []string{ "filename.png" },
		GridRows: 5,
		GridCols: 10,
	}
	errs := config.Validate()

	assert.Len(t, errs, 0)
}

func TestValidate_AfterSetDefaults(t *testing.T) {
	config := Config { InputFiles: []string{"filename.jpg"} }
	config.SetDefaults()

	errs := config.Validate()

	assert.Len(t, errs, 0)
}

func TestValidate_NoInputs(t *testing.T) {
	config := Config { GridRows: 10, GridCols: 10 }
	
	errs := config.Validate()

	assert.Len(t, errs, 1)
	assert.ErrorContains(t, errs[0], "not enough input files")
}
func TestValidate_InvalidState(t *testing.T) {
	config := Config {
		InputFiles: nil,
		GridRows: 0,
		GridCols: -10,
	}

	errs := config.Validate()

	assert.Len(t, errs, 3)
	assert.ErrorContains(t, errs[0], "not enough input files")
	assert.ErrorContains(t, errs[1], "number of grid rows")
	assert.ErrorContains(t, errs[2], "number of grid columns")
}