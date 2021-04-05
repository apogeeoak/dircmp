package compare

import "fmt"

type Stats struct {
	filesSearched        int
	differentFiles       int
	differentDirectories int
	errors               int
}

func (s *Stats) Add(statType StatType) {
	switch statType {
	case StatSearchedFile:
		s.filesSearched++
	case StatDifferentFile:
		s.differentFiles++
	case StatDifferentDirectory:
		s.differentDirectories++
	case StatError:
		s.errors++
	default:
		// No-op.
	}
}

func (s *Stats) String() string {
	total := s.differentFiles + s.differentDirectories
	return fmt.Sprintf("Searched %d file(s), %d file(s) different, %d director(ies) different, %d total entr(ies) different. %d error(s).", s.filesSearched, s.differentFiles, s.differentDirectories, total, s.errors)
}

type StatType int

const (
	StatNone StatType = iota
	StatSearchedFile
	StatDifferentFile
	StatDifferentDirectory
	StatError
)

func (s StatType) String() string {
	switch s {
	case StatNone:
		return "None"
	case StatSearchedFile:
		return "Searched file"
	case StatDifferentFile:
		return "Different file"
	case StatDifferentDirectory:
		return "Different directory"
	case StatError:
		return "Error"
	default:
		panic("Undefined enum string value.")
	}
}
