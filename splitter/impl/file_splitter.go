package impl

import (
	"io/ioutil"
	"path"
)

// FileSplitter splits input by file
type FileSplitter struct {
}

// NewFileSplitter return a new file splitter
func NewFileSplitter() *FileSplitter {
	return &FileSplitter{}
}

// Split implements of FileSplitter
// {Input}: a directory
// return an array of files in the input
func (fs *FileSplitter) Split(input string) ([]string, error) {
	var files []string

	dir, err := ioutil.ReadDir(input)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		files = append(files, path.Join(input, fi.Name()))
	}

	return files, nil
}
