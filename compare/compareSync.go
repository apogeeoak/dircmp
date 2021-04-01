package compare

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/apogeeoak/dircmp/lib/collection"
)

func CompareSync(config *Config) (*Stats, error) {
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
			return nil, err
		}
		comp, err := os.ReadDir(filepath.Join(config.Compared, dir))
		if err != nil {
			return nil, err
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
				compareDirectoriesSync(config, oEntry, cEntry, path, directories, stats)
			} else {
				err := compareFilesSync(config, oEntry, cEntry, path, stats)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return stats, nil
}

func compareDirectoriesSync(config *Config, orig fs.DirEntry, comp fs.DirEntry, path string, directories *collection.Stack, stats *Stats) {
	// Comparison failed on non-empty string.
	if cmp := compareDirectories(orig, comp); cmp != "" {
		stats.DiffDir()
		printOutput(cmp, path)
	} else {
		directories.Push(path)
	}
}

func compareFilesSync(config *Config, orig fs.DirEntry, comp fs.DirEntry, path string, stats *Stats) error {
	stats.AddSearched()
	// Comparison failed on non-empty string.
	cmp, err := compareFiles(config, orig, comp, path)
	if err != nil {
		return err
	}
	if cmp != "" {
		stats.DiffFile()
		printOutput(cmp, path)
	}
	return nil
}

func printOutput(output, path string) {
	fmt.Printf("%-30s | %s\n", output, path)
}
