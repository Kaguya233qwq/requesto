package requesto

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// File represents a file to be uploaded, containing its name and content as an io.Reader.
type File struct {
	Name    string
	Content io.Reader
}

// FileFromPath creates a File object from a given file path.
// It opens the file and uses its base name as the file name.
func FileFromPath(filePath string) (File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:    filepath.Base(filePath),
		Content: file,
	}, nil
}

// FileFromBytes creates a File object from a byte slice.
// It uses the provided file name and wraps the byte slice in an io.Reader.
func FileFromBytes(fileName string, data []byte) File {
	return File{
		Name:    fileName,
		Content: bytes.NewReader(data),
	}
}
