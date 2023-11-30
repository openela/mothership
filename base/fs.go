package base

import (
	"embed"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

// EmbedFSToOSFS re-creates the structure of the embed.FS in the OS filesystem.
// It is mostly useful for testing but other uses are possible.
// This function takes a prefix / directory and an embed.FS and creates the
// corresponding directory structure in the OS filesystem.
// Returns all created file paths.
func EmbedFSToOSFS(prefix string, fs embed.FS, fsDirectory string) ([]string, error) {
	// Make sure the prefix exists as a directory
	err := os.MkdirAll(prefix, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "could not create directory")
	}

	// Read the directory entries
	entries, err := fs.ReadDir(fsDirectory)
	if err != nil {
		return nil, errors.Wrap(err, "could not read directory")
	}

	var res []string

	// Iterate over the entries
	for _, entry := range entries {
		// Create the full path
		full := prefix + "/" + entry.Name()

		// Check if the entry is a directory
		if entry.IsDir() {
			// If it is, recurse
			res2, err := EmbedFSToOSFS(full, fs, fsDirectory+"/"+entry.Name())
			if err != nil {
				return nil, errors.Wrap(err, "could not recurse")
			}
			res = append(res, res2...)
		} else {
			// If not, create the file
			f, err := os.Create(full)
			if err != nil {
				return nil, errors.Wrap(err, "could not create file")
			}

			// Open the file in the embed.FS
			file, err := fs.Open(strings.TrimPrefix(fsDirectory+"/"+entry.Name(), "./"))
			if err != nil {
				return nil, errors.Wrap(err, "could not open file")
			}

			// Copy the file
			_, err = io.Copy(f, file)
			if err != nil {
				return nil, errors.Wrap(err, "could not copy file")
			}

			// Close the OS file
			err = f.Close()
			if err != nil {
				return nil, errors.Wrap(err, "could not close file")
			}

			// Close the embed.FS file
			err = file.Close()
			if err != nil {
				return nil, errors.Wrap(err, "could not close file")
			}

			res = append(res, full)
		}
	}

	return res, nil
}
