package compare

import "fmt"

type Stats struct {
	filesSearched        int64
	differentFiles       int64
	differentDirectories int64
}

func (s *Stats) Add(statType StatType) {
	switch statType {
	case SearchedFile:
		s.AddSearched()
	case DifferentFile:
		s.DiffFile()
	case DifferentDirectory:
		s.DiffDir()
	default:
		// No-op.
	}
}

func (s *Stats) AddSearched() {
	s.filesSearched++
}

func (s *Stats) DiffDir() {
	s.differentDirectories++
}

func (s *Stats) DiffFile() {
	s.differentFiles++
}

func (s *Stats) String() string {
	total := s.differentFiles + s.differentDirectories
	return fmt.Sprintf("Searched %d file(s), %d file(s) different, %d director(ies) different, %d total entr(ies) different.", s.filesSearched, s.differentFiles, s.differentDirectories, total)
}

type StatType int

const (
	None StatType = iota
	SearchedFile
	DifferentFile
	DifferentDirectory
)

func (s StatType) String() string {
	switch s {
	case SearchedFile:
		return "Searched file."
	case DifferentFile:
		return "Different file."
	case DifferentDirectory:
		return "Different directory."
	default:
		return "None."
	}
}
