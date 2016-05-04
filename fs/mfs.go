package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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

func (fe *MfsFileExplorer) Mkdir(newPath string) error {
	return os.Mkdir(normalizePath(newPath), os.ModePerm)
}

func (fe *MfsFileExplorer) Cat(path string) (string, error) {
	fi, err := os.Open(normalizePath(path))
	if err != nil {
		return "", err
	}
	defer fi.Close()

	chunks := make([]byte, 1024, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			return "", nil
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks), nil
}

func (fe *MfsFileExplorer) Write(path, content string) error {
	fo, err := os.Create(normalizePath(path))
	if err != nil {
		return err
	}
	defer fo.Close()
	fo.WriteString(content)
	return nil
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
	// at there we do not have to normalize path
	// it will be normalized in function copy
	err := fe.Copy(path, newPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(normalizePath(path))
	return err
}

func (fe *MfsFileExplorer) Copy(path string, newPath string) error {
	path = normalizePath(path)
	newPath = normalizePath(newPath)
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return utils.CopyDir(path, newPath)
	}
	return utils.CopyFile(path, newPath)
}

func (fe *MfsFileExplorer) Download(path string) ([]byte, error) {
	fi, err := os.Open(normalizePath(path))
	defer fi.Close()
	if err != nil {
		return nil, err
	}
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func (fe *MfsFileExplorer) Upload(path string, file *multipart.FileHeader) error {
	fi, err := file.Open()
	defer fi.Close()
	if err != nil {
		return err
	}

	fileNew := normalizePath(path + "/" + file.Filename)
	fo, err := os.OpenFile(fileNew, os.O_WRONLY|os.O_CREATE, 0777)
	defer fo.Close()
	if err != nil {
		return err
	}

	io.Copy(fo, fi)
	return nil
}

func (fe *MfsFileExplorer) Delete(path string) error {
	return os.RemoveAll(normalizePath(path))
}

func (fe *MfsFileExplorer) Rename(oldPath, newPath string) error {
	return os.Rename(normalizePath(oldPath), normalizePath(newPath))
}
