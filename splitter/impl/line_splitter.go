package impl

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// LineSplitter splits input by file
type LineSplitter struct {
}

// NewLineSplitter return a new file splitter
func NewLineSplitter() *LineSplitter {
	return &LineSplitter{}
}

// Split implements of LineSplitter
// {Input}: a file
// return every line of this file
func (ls *LineSplitter) Split(input string) ([]string, error) {
	var args []string
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		args = append(args, line)
	}
	return args, nil
}
