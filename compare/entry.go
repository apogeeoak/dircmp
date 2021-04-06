package compare

import "io/fs"

type Entry struct {
	Original fs.DirEntry
	Compared fs.DirEntry
	Path     string
}

func NewEntry(orig, comp fs.DirEntry, path string) Entry {
	return Entry{Original: orig, Compared: comp, Path: path}
}
