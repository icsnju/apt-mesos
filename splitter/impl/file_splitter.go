package impl

import (
	"io/ioutil"
	"path"

	"github.com/icsnju/apt-mesos/fs"
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
func (f *FileSplitter) Split(input string) ([]string, error) {
	var files []string

	dir, err := ioutil.ReadDir(fs.NormalizePath(input))
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		files = append(files, path.Join(input, fi.Name()))
	}

	return files, nil
}
