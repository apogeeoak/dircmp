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
		orig, err := os.ReadDir(filepath.Join(config.Original, dir))
		if err != nil {
			process(Error(err), stats)
			continue
		}
		comp, err := os.ReadDir(filepath.Join(config.Compared, dir))
		if err != nil {
			process(Error(err), stats)
			continue
		}

		oIndex := 0
		for _, cEntry := range comp {
			path := filepath.Join(dir, cEntry.Name())

			// Search for original entry that matches compared entry.
			var oEntry fs.DirEntry
			for ; oIndex < len(orig) && orig[oIndex].Name() <= cEntry.Name(); oIndex++ {
				oEntry = orig[oIndex]
			}

			// Branch on directory or file.
			if cEntry.IsDir() {
				compareDirectoriesSerial(config, oEntry, cEntry, path, directories, stats)
			} else {
				compareFilesSerial(config, oEntry, cEntry, path, stats)
			}
		}
	}
	return stats, nil
}

func compareDirectoriesSerial(config *Config, orig fs.DirEntry, comp fs.DirEntry, path string, directories *collection.Stack, stats *Stats) {
	// Comparison failed on non-empty string.
	if cmp := compareDirectories(orig, comp); cmp != "" {
		process(Output(cmp, path, StatDifferentDirectory), stats)
	} else {
		directories.Push(path)
	}
}

func compareFilesSerial(config *Config, orig fs.DirEntry, comp fs.DirEntry, path string, stats *Stats) {
	stats.Add(StatSearchedFile)
	// Comparison failed on non-empty string.
	cmp, err := compareFiles(config, orig, comp, path)
	if err != nil {
		process(Error(err), stats)
	} else if cmp != "" {
		process(Output(cmp, path, StatDifferentFile), stats)
	}
}
