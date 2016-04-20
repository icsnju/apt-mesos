package fs

import (
	"fmt"
	"os"
	"strings"

	"github.com/icsnju/apt-mesos/utils"
)

var (
	MfsRoot = "/mnt/mfsmount"
)

type MfsFileExplorer struct {
	FileExplorer
}

func NewMfsFileExplorer() *MfsFileExplorer {
	return &MfsFileExplorer{}
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return MfsRoot + "/" + path
	}
	return MfsRoot + path
}

func (fe *MfsFileExplorer) Mkdir(path string, name string) error {
	return os.Mkdir(normalizePath(path)+"/"+name, os.ModeDir)
}

func (fe *MfsFileExplorer) ListDir(path string) ([]ListDirEntry, error) {
	directory, err := os.Open(normalizePath(path))
	if err != nil {
		return nil, err
	}
	files, err := directory.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var list []ListDirEntry
	for _, file := range files {
		entry := ListDirEntry{
			Name:   file.Name(),
			Rights: file.Mode().String(),
			Date:   file.ModTime().String(),
			Size:   fmt.Sprintf("%d", file.Size()),
		}
		if file.IsDir() {
			entry.Type = "dir"
		} else {
			entry.Type = "file"
		}
		list = append(list, entry)
	}
	return list, nil
}

func (fe *MfsFileExplorer) Move(path string, newPath string) error {
	err := fe.Copy(path, newPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(path)
	return err
}

func (fe *MfsFileExplorer) Copy(path string, newPath string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return utils.CopyDir(path, newPath)
	}
	return utils.CopyFile(path, newPath)
}

func (fe *MfsFileExplorer) Delete(path string) error {
	return os.RemoveAll(normalizePath(path))
}

func (fe *MfsFileExplorer) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
