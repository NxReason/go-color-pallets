package cmd

type ArgumentData struct {
	flag        string
	description string
	required    bool
}

type Argument interface {
	Validate() error
	AddValue(string)
}

// Specific arguments
// Input file(s)
type InputFileArgument struct {
	ArgumentData
	files []string
}

func (arg *InputFileArgument) AddValue(value string) {
	arg.files = append(arg.files, value)
}

func (arg *InputFileArgument) Validate() error {
	if len(arg.files) == 0 {
		return ArgumentError{"not enough files to process (after -i flag), at lease 1 should be provided"}
	}
	return nil
}

func (arg *InputFileArgument) GetFiles() []string {
	return arg.files
}

// Output file(s)
type OutputFileArgument struct {
	ArgumentData
	files []string
}

func (arg *OutputFileArgument) AddValue(value string) {
	arg.files = append(arg.files, value)
}

func (arg *OutputFileArgument) Validate() error {
	return nil
}

// Custom argument error
type ArgumentError struct {
	msg string
}

func (ae ArgumentError) Error() string {
	return "argument error: " + ae.msg
}