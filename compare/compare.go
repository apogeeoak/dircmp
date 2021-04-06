package compare

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/apogeeoak/dircmp/lib/collection"
)

func Compare(config *Config) (*Stats, error) {
	// Ensure starting directories exist.
	if err := directoriesExists(config.Original, config.Compared); err != nil {
		return nil, err
	}

	// Initialize.
	stats := &Stats{}
	entries := make(chan Entry)
	results := make(chan Result)
	wait := &sync.WaitGroup{}

	// Start reading directories.
	go readDirectories(config, entries, results)

	// Start file comparison goroutines.
	limit := int(config.Limit)
	wait.Add(limit)
	for i := 0; i < limit; i++ {
		go func() {
			defer wait.Done()
			compareFilesParallel(config, entries, results)
		}()
	}

	// Close results when comparsion is done.
	go func() {
		wait.Wait()
		close(results)
	}()

	// Listen for results.
	for result := range results {
		process(result, stats)
	}
	return stats, nil
}

func readDirectories(config *Config, entries chan<- Entry, results chan<- Result) {
	defer close(entries)

	directories := &collection.Stack{}
	directories.Push("")

	// Iterate through directories in stack.
	for dir, err := directories.Pop(); err == nil; dir, err = directories.Pop() {
		// Read contents of directories.
		oEntries, err := os.ReadDir(filepath.Join(config.Original, dir))
		if err != nil {
			results <- Error(err)
			continue
		}
		cEntries, err := os.ReadDir(filepath.Join(config.Compared, dir))
		if err != nil {
			results <- Error(err)
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
				compareDirectoriesParallel(config, entry, directories, results)
			} else {
				entries <- entry
			}
		}
	}
}

func compareDirectoriesParallel(config *Config, entry Entry, directories *collection.Stack, results chan<- Result) {
	// Comparison failed if result is not empty.
	if result := compareDirectories(entry); !result.IsEmpty() {
		results <- result
	} else {
		directories.Push(entry.Path)
	}
}

func compareFilesParallel(config *Config, entries <-chan Entry, results chan<- Result) {
	for entry := range entries {
		results <- Stat(StatSearchedFile)

		// Comparison failed if result is not empty.
		if result := compareFiles(config, entry); !result.IsEmpty() {
			results <- result
		}
	}
}

func compareDirectories(e Entry) Result {
	// Comparison failed: Directory only in compared.
	if e.Original == nil || e.Original.Name() != e.Compared.Name() || !e.Original.IsDir() {
		return Output("Directory only in compared.", e.Path, StatDifferentDirectory)
	}
	return Empty()
}

func compareFiles(config *Config, e Entry) Result {
	// Comparison failed: File only in compared.
	if e.Original == nil || e.Original.Name() != e.Compared.Name() || e.Original.IsDir() {
		return Output("File only in compared.", e.Path, StatDifferentFile)
	}

	// Determine file size from FileInfo.
	oInfo, err := e.Original.Info()
	if err != nil {
		return Error(err)
	}
	cInfo, err := e.Compared.Info()
	if err != nil {
		return Error(err)
	}

	// Comparison failed: File sizes differ.
	if oInfo.Size() != cInfo.Size() {
		return Output("File size differs.", e.Path, StatDifferentFile)
	}

	// Ensure offset is positive.
	offset := max(config.Offset(oInfo.Size()), 0)

	// Open files.
	orig, err := os.Open(filepath.Join(config.Original, e.Path))
	if err != nil {
		return Error(err)
	}
	defer orig.Close()
	comp, err := os.Open(filepath.Join(config.Compared, e.Path))
	if err != nil {
		return Error(err)
	}
	defer comp.Close()

	// Read files.
	return compareFilesRead(config, orig, comp, e.Path, offset)
}

func compareFilesRead(config *Config, orig io.ReadSeeker, comp io.ReadSeeker, path string, offset int64) Result {
	oBytes := make([]byte, config.SampleSize)
	cBytes := make([]byte, config.SampleSize)

	for {
		_, oErr := orig.Read(oBytes)
		_, cErr := comp.Read(cBytes)

		// Error conditions.
		if oErr != nil || cErr != nil {
			// Comparison succeeded.
			if oErr == io.EOF && cErr == io.EOF {
				return Empty()
			}
			// Comparison failed: One file ended before the other.
			if oErr == io.EOF || cErr == io.EOF {
				return Output("One file ended before the other.", path, StatDifferentFile)
			}
			// Error out.
			return Error(fmt.Errorf("unable to read files: %v; %v", oErr, cErr))
		}

		// Comparison failed: File contents differ.
		if !bytes.Equal(oBytes, cBytes) {
			return Output("File content differs.", path, StatDifferentFile)
		}

		// Offset both files relative to current position.
		orig.Seek(offset, 1)
		comp.Seek(offset, 1)
	}
}

func process(result Result, stats *Stats) {
	stats.Add(result.Stat)
	if result.Error != nil {
		fmt.Fprintln(os.Stderr, result)
	} else if result.Output != "" {
		fmt.Println(result)
	}
}

func directoriesExists(directories ...string) error {
	for _, dir := range directories {
		stat, err := os.Stat(dir)
		if os.IsNotExist(err) || !stat.IsDir() {
			return fmt.Errorf("%s: no such directory", dir)
		}
	}
	return nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
