package compare

import "fmt"

type Stat struct {
	filesSearched        int64
	differentDirectories int64
	differentFiles       int64
}

func (s *Stat) AddSearched() {
	s.filesSearched++
}

func (s *Stat) DiffDir() {
	s.differentDirectories++
}

func (s *Stat) DiffFile() {
	s.differentFiles++
}

func (s *Stat) String() string {
	total := s.differentFiles + s.differentDirectories
	return fmt.Sprintf("Searched %d file(s), %d file(s) different, %d director(ies) different, %d total entr(ies) different", s.filesSearched, s.differentFiles, s.differentDirectories, total)
}
