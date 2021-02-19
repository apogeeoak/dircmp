package compare

import (
	"flag"
	"fmt"
	"os"
)

const (
	versionNumber = "0.2"
)

type offsetFunc func(size int64) int64

type Config struct {
	Original string
	Compared string

	SampleSize int
	Offset     offsetFunc
}

func ParseConfig() *Config {
	return ParseConfigArgs(os.Args[0], os.Args[1:])
}

func ParseConfigArgs(name string, args []string) *Config {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	flags.Usage = usage

	// Flags.
	entire := false
	entire_usage := "Read entire file for comparison. More accurate but slower."
	flags.BoolVar(&entire, "e", entire, entire_usage)
	flags.BoolVar(&entire, "entire", entire, entire_usage)

	version := false
	version_usage := "Output version information and exit."
	flags.BoolVar(&version, "V", version, version_usage)
	flags.BoolVar(&version, "version", version, version_usage)

	samples := 4
	samples_usage := "Number of samples to take."
	flags.IntVar(&samples, "samples", samples, samples_usage)

	sampleSize := 4000
	size_usage := "Size of samples."
	flags.IntVar(&sampleSize, "size", sampleSize, size_usage)

	// Parse flags.
	flags.Parse(args)

	// Output version and exit if defined.
	versionOutput(version)

	// Arguments.
	original := flags.Arg(0)
	compared := flags.Arg(1)
	if original == "" {
		failArgumentNotDefined("ORIGINAL")
	}
	if compared == "" {
		failArgumentNotDefined("COMPARED")
	}
	if flags.NArg() > 2 {
		failTooManyArguments(flags.Args())
	}

	offset := offset(entire, samples, sampleSize)

	return &Config{original, compared, sampleSize, offset}
}

func failf(format string, a ...interface{}) {
	output := flag.CommandLine.Output()
	err := fmt.Errorf(format, a...)
	fmt.Fprintf(output, "dircmp: %s\n", err)
	fmt.Fprintf(output, "dircmp: Try 'dircmp --help' for more information.\n")
	os.Exit(2)
}

func offset(entire bool, samples int, sampleSize int) offsetFunc {
	// No offset if entire flag is set to read entire file.
	if entire {
		return func(int64) int64 { return 0 }
	}

	// Offset amount between reads to obtain equidistant samples. The last sample may not include the final few bytes.
	return func(size int64) int64 {
		return (size - int64(samples*sampleSize)) / int64(samples-1)
	}
}

func usage() {
	output := flag.CommandLine.Output()
	fmt.Fprintf(output, "Usage: dircmp [OPTION]... ORIGINAL COMPARED\n")
	fmt.Fprintf(output, "Compare directory COMPARED to directory ORIGINAL.\n")
	fmt.Fprintf(output, "Options:\n")
	flag.PrintDefaults()
}

func versionOutput(version bool) {
	if version {
		fmt.Println("dircmp version", versionNumber)
		os.Exit(0)
	}
}

func failArgumentNotDefined(arg string) {
	failf("argument %s missing.", arg)
}

func failTooManyArguments(arg []string) {
	failf("too many arguments %s.", arg)
}
