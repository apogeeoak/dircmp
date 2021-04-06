package compare

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/apogeeoak/dircmp/lib/collection"
)

func CompareSerial(config *Config) (*Stats, error) {
	// Ensure starting directories exist.
	if err := directoriesExists(config.Original, config.Compared); err != nil {
		return nil, err
	}

	// Initialize.
	stats := &Stats{}
	directories := &collection.Stack{}
	directories.Push("")

	// Iterate through directories in stack.
	for dir, err := directories.Pop(); err == nil; dir, err = directories.Pop() {
		// Read contents of directories.
		oEntries, err := os.ReadDir(filepath.Join(config.Original, dir))
		if err != nil {
			process(Error(err), stats)
			continue
		}
		cEntries, err := os.ReadDir(filepath.Join(config.Compared, dir))
		if err != nil {
			process(Error(err), stats)
			continue
		}

		oIndex := 0
		for _, comp := range cEntries {
			path := filepath.Join(dir, comp.Name())

			// Search for original entry that matches compared entry.
			var orig fs.DirEntry
			for ; oIndex < len(oEntries) && oEntries[oIndex].Name() <= comp.Name(); oIndex++ {
				orig = oEntries[oIndex]
			}
			entry := NewEntry(orig, comp, path)

			// Branch on directory or file.
			if comp.IsDir() {
				compareDirectoriesSerial(config, entry, directories, stats)
			} else {
				compareFilesSerial(config, entry, stats)
			}
		}
	}
	return stats, nil
}

func compareDirectoriesSerial(config *Config, entry Entry, directories *collection.Stack, stats *Stats) {
	// Comparison failed if result is not empty.
	if result := compareDirectories(entry); !result.IsEmpty() {
		process(result, stats)
	} else {
		directories.Push(entry.Path)
	}
}

func compareFilesSerial(config *Config, entry Entry, stats *Stats) {
	stats.Add(StatSearchedFile)
	// Comparison failed if result is not empty.
	if result := compareFiles(config, entry); !result.IsEmpty() {
		process(result, stats)
	}
}
