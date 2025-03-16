package config

import "slices"

const (
	OutFmtRaw OutFormat = iota
	OutFmtText

	// must be last
	numOutFormats
)

// TODO(keithpaterson): could add more output types, e.g. OutBase64, or something like that

var (
	outNames = []string{"raw", "text"}
)

type OutFormat int

// OutputConfig stores the configuration for how to output the response.
// Templates are described in more detail in the [output] package.
type OutputConfig struct {
	// Specify the output format
	//  "raw": dumps whatever is in the response as-is
	//  "text": converts the response to text and applies the template before writing it out.
	Format string `toml:"format,required"     validate:"required,oneof=raw text"`

	// Specify a file to write the response data into.
	// The special strings "stdout" and "stderr" can be used to direct output to the system streams.
	Filename string `toml:"filename,omitempty"  validate:"omitempty,gt=0"`

	// A text template that can add additional information to the output stream.
	// If Template is not specified, the default template "${response:Body}" will be used.
	Template string `toml:"template,omitempty"  validate:"omitempty,gt=0"`
}

type Options map[string]string

func newOutputConfig() OutputConfig {
	return OutputConfig{Format: "text"}
}

func (of OutFormat) String() string {
	if of < 0 || of >= numOutFormats {
		return "undefined"
	}
	return outNames[of]
}

func (o OutputConfig) OutFormat() OutFormat {
	index := slices.Index(outNames, o.Format)
	if index < 0 {
		return OutFmtText
	}
	return OutFormat(index)
}
