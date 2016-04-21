package fs

import "mime/multipart"

type FileExplorer interface {
	Init() error
	ListDir(path string) ([]ListDirEntry, error)
	Move(path string, newPath string) error
	Copy(path string, newPath string) error
	Delete(path string) error
	Rename(oldPath string, newPath string) error
	// Chmod(path string, code string) error
	Mkdir(newPath string) error
	Cat(path string) (string, error)
	Write(path, content string) error
	Download(path string) ([]byte, error)
	Upload(path string, file *multipart.FileHeader) error
	Close() error
}
