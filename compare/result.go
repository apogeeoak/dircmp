package compare

import "fmt"

type Result struct {
	Output string
	Path   string
	Stat   StatType
	Error  error
}

func Error(err error) Result {
	return Result{Error: err}
}

func Output(output, path string) Result {
	return Result{Output: output, Path: path}
}

func Stat(stat StatType) Result {
	return Result{Stat: stat}
}

func (r Result) String() string {
	if r.Error != nil {
		return "ERROR: " + r.Error.Error()
	} else if r.Stat != None {
		return r.Stat.String()
	} else if r.Output != "" {
		return fmt.Sprintf("%-30s | %s", r.Output, r.Path)
	}
	return "Empty Result."
}
