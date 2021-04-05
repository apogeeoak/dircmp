package compare

import "fmt"

type Result struct {
	Output string
	Path   string
	Stat   StatType
	Error  error
}

func Error(err error) Result {
	return Result{Error: err, Stat: StatError}
}

func Output(output, path string, stat StatType) Result {
	return Result{Output: output, Path: path, Stat: stat}
}

func Stat(stat StatType) Result {
	return Result{Stat: stat}
}

func (r Result) String() string {
	if r.Error != nil {
		return "Error: " + r.Error.Error()
	} else if r.Output != "" {
		return fmt.Sprintf("%-30s | %s", r.Output, r.Path)
	} else if r.Stat != StatNone {
		return r.Stat.String()
	}
	return "Empty Result."
}
