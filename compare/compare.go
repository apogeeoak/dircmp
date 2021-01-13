package compare

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apogeeoak/dircmp/lib/collection"
)

const (
	samples = 4
	size    = 4000
)

func Compare(original string, compared string) (*Stat, error) {
	stat := &Stat{}

	// Ensure starting directories exist.
	if err := directoriesExists(original, compared); err != nil {
		return nil, err
	}

	directories := collection.Stack{}
	directories.Push("")

	for !directories.IsEmpty() {
		// Read next directory on stack.
		dir, err := directories.Pop()
		if err != nil {
			return nil, err
		}

		// Read contents of directories
		orig, err := ioutil.ReadDir(filepath.Join(original, dir))
		if err != nil {
			return nil, err
		}
		comp, err := ioutil.ReadDir(filepath.Join(compared, dir))
		if err != nil {
			return nil, err
		}

		oIndex := 0
		for _, cInfo := range comp {
			path := filepath.Join(dir, cInfo.Name())

			// Search for original FileInfo that matches compared FileInfo.
			var oInfo os.FileInfo
			for ; oIndex < len(orig) && orig[oIndex].Name() <= cInfo.Name(); oIndex++ {
				oInfo = orig[oIndex]
			}

			// Branch on directory or file.
			if cInfo.IsDir() {
				// Comparison failed on non-empty string.
				if cmp := compareDirectories(oInfo, cInfo); cmp != "" {
					stat.DiffDir()
					printOutput(path, cmp)
				} else {
					directories.Push(path)
				}
			} else {
				stat.AddSearched()
				// Comparison failed on non-empty string.
				cmp, err := compareFiles(original, compared, oInfo, cInfo, path)
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

func compareDirectories(oInfo os.FileInfo, cInfo os.FileInfo) string {
	// Comparison failed: Directory only in compared.
	if oInfo == nil || oInfo.Name() != cInfo.Name() || !oInfo.IsDir() {
		return "Directory only in compared"
	}
	return ""
}

func compareFiles(original string, compared string, oInfo os.FileInfo, cInfo os.FileInfo, path string) (string, error) {
	oPath := filepath.Join(original, path)
	cPath := filepath.Join(compared, path)

	// Comparison failed: File only in compared.
	if oInfo == nil || oInfo.Name() != cInfo.Name() || oInfo.IsDir() {
		return "File only in compared.", nil
	}

	// Comparison failed: File sizes differs.
	if oInfo.Size() != cInfo.Size() {
		return "File sizes differs.", nil
	}

	// Open files.
	oFile, err := os.Open(oPath)
	if err != nil {
		return "", err
	}
	defer oFile.Close()
	cFile, err := os.Open(cPath)
	if err != nil {
		return "", err
	}
	defer cFile.Close()

	// Offset amount between reads to obtain equidistant samples. The last sample may not include the final few bytes.
	offset := (oInfo.Size() - samples*size) / (samples - 1)

	// Read files.
	return compareFilesRead(oFile, cFile, size, offset)
}

func compareFilesRead(oFile *os.File, cFile *os.File, size int, offset int64) (string, error) {
	oBytes := make([]byte, size)
	cBytes := make([]byte, size)
	// Ensure offset is positive.
	offset = max(offset, 0)

	for {
		_, oErr := oFile.Read(oBytes)
		_, cErr := cFile.Read(cBytes)

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
		oFile.Seek(offset, 1)
		cFile.Seek(offset, 1)
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
