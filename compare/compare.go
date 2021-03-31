package compare

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/apogeeoak/dircmp/lib/collection"
)

func CompareSync(config *Config) (*Stat, error) {
	// Ensure starting directories exist.
	if err := directoriesExists(config.Original, config.Compared); err != nil {
		return nil, err
	}

	stat := &Stat{}
	directories := collection.Stack{}
	directories.Push("")

	for !directories.IsEmpty() {
		// Read next directory on stack.
		dir, err := directories.Pop()
		if err != nil {
			return nil, err
		}

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
				// Comparison failed on non-empty string.
				if cmp := compareDirectories(oEntry, cEntry); cmp != "" {
					stat.DiffDir()
					printOutput(path, cmp)
				} else {
					directories.Push(path)
				}
			} else {
				stat.AddSearched()
				// Comparison failed on non-empty string.
				cmp, err := compareFiles(config, oEntry, cEntry, path)
				if err != nil {
					return nil, err
				}
				if cmp != "" {
					stat.DiffFile()
					printOutput(path, cmp)
				}
			}
		}
	}
	return stat, nil
}

func compareDirectories(orig fs.DirEntry, comp fs.DirEntry) string {
	// Comparison failed: Directory only in compared.
	if orig == nil || orig.Name() != comp.Name() || !orig.IsDir() {
		return "Directory only in compared"
	}
	return ""
}

func compareFiles(config *Config, orig fs.DirEntry, comp fs.DirEntry, path string) (string, error) {
	// Comparison failed: File only in compared.
	if orig == nil || orig.Name() != comp.Name() || orig.IsDir() {
		return "File only in compared.", nil
	}

	// Determine file size from FileInfo.
	oInfo, err := orig.Info()
	if err != nil {
		return "", err
	}
	cInfo, err := comp.Info()
	if err != nil {
		return "", err
	}

	// Comparison failed: File sizes differs.
	if oInfo.Size() != cInfo.Size() {
		return "File sizes differs.", nil
	}

	// Ensure offset is positive.
	offset := max(config.Offset(oInfo.Size()), 0)

	// Open files.
	oFile, err := os.Open(filepath.Join(config.Original, path))
	if err != nil {
		return "", err
	}
	defer oFile.Close()
	cFile, err := os.Open(filepath.Join(config.Compared, path))
	if err != nil {
		return "", err
	}
	defer cFile.Close()

	// Read files.
	return compareFilesRead(config, oFile, cFile, offset)
}

func compareFilesRead(config *Config, orig *os.File, comp *os.File, offset int64) (string, error) {
	oBytes := make([]byte, config.SampleSize)
	cBytes := make([]byte, config.SampleSize)

	for {
		_, oErr := orig.Read(oBytes)
		_, cErr := comp.Read(cBytes)

		// Error conditions.
		if oErr != nil || cErr != nil {
			// Comparison succeeded.
			if oErr == io.EOF && cErr == io.EOF {
				return "", nil
			}
			// Comparison failed: One file ended before the other.
			if oErr == io.EOF || cErr == io.EOF {
				return "One file ended before the other.", nil
			}
			// Error out.
			return "", fmt.Errorf("compareFiles original %s; compared %s", oErr, cErr)
		}

		// Comparison failed: File contents differ.
		if !bytes.Equal(oBytes, cBytes) {
			return "File content differs.", nil
		}

		// Offset both files relative to current position.
		orig.Seek(offset, 1)
		comp.Seek(offset, 1)
	}
}

func directoriesExists(directories ...string) error {
	for _, dir := range directories {
		stat, err := os.Stat(dir)
		if os.IsNotExist(err) || !stat.IsDir() {
			return errNotDirectory(dir)
		}
	}
	return nil
}

func errNotDirectory(dir string) error {
	return fmt.Errorf("%s: no such directory", dir)
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func printOutput(path string, message string) {
	fmt.Printf("%-30s | %s\n", path, message)
}
