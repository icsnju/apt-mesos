package utils

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// Tar package files to a tar.gz file
// src source file
// dstTar destination tar file
// failIfExist if tar file exists, whether to overwrite it
func Tar(src string, dstTar string, failIfExist bool) (err error) {
	src = path.Clean(src)

	// check if source file not exists
	if !Exists(src) {
		return errors.New("file or directory not exist：" + src)
	}

	// check if tar file not exists
	if FileExists(dstTar) {
		if failIfExist {
			return errors.New("Tar file exists: " + dstTar)
		} else {
			if err := os.Remove(dstTar); err != nil {
				return err
			}
		}
	}

	// Create empty tar
	fw, er := os.Create(dstTar)
	if er != nil {
		return er
	}
	defer fw.Close()

	tw := tar.NewWriter(fw)
	defer func() {
		if er := tw.Close(); er != nil {
			err = er
		}
	}()

	fi, er := os.Stat(src)
	if er != nil {
		return er
	}

	if fi.IsDir() {
		tarDir(src, "", tw, fi)
	} else {
		tarFile(src, "", tw, fi)
	}

	return nil
}

func tarDir(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	srcFull := srcBase + srcRelative
	fis, er := ioutil.ReadDir(srcFull)
	if er != nil {
		return er
	}

	for _, fi := range fis {
		if fi.IsDir() {
			tarDir(srcBase, srcRelative+"/"+fi.Name(), tw, fi)
		} else {
			tarFile(srcBase, srcRelative+"/"+fi.Name(), tw, fi)
		}
	}

	if len(srcRelative) > 0 {
		fmt.Println(123)
		hdr, er := tar.FileInfoHeader(fi, "")
		if er != nil {
			return er
		}
		hdr.Name = srcRelative
		if er = tw.WriteHeader(hdr); er != nil {
			return er
		}
	}

	return nil
}

func tarFile(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	srcFull := srcBase + srcRelative

	hdr, er := tar.FileInfoHeader(fi, "")
	if er != nil {
		return er
	}
	hdr.Name = srcRelative
	if er = tw.WriteHeader(hdr); er != nil {
		return er
	}

	fr, er := os.Open(srcFull)
	if er != nil {
		return er
	}
	defer fr.Close()

	if _, er = io.Copy(tw, fr); er != nil {
		return er
	}
	return nil
}

// UnTar a tar file
func UnTar(srcTar string, dstDir string) (err error) {
	// 清理路径字符串
	dstDir = path.Clean(dstDir) + string(os.PathSeparator)

	fr, er := os.Open(srcTar)
	if er != nil {
		return er
	}
	defer fr.Close()

	tr := tar.NewReader(fr)

	for hdr, er := tr.Next(); er != io.EOF; hdr, er = tr.Next() {
		if er != nil {
			return er
		}

		fi := hdr.FileInfo()

		dstFullPath := dstDir + hdr.Name

		if hdr.Typeflag == tar.TypeDir {
			os.MkdirAll(dstFullPath, fi.Mode().Perm())
			os.Chmod(dstFullPath, fi.Mode().Perm())
		} else {
			os.MkdirAll(path.Dir(dstFullPath), os.ModePerm)
			if er := unTarFile(dstFullPath, tr); er != nil {
				return er
			}
			os.Chmod(dstFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

func unTarFile(dstFile string, tr *tar.Reader) error {
	fw, er := os.Create(dstFile)
	if er != nil {
		return er
	}
	defer fw.Close()

	_, er = io.Copy(fw, tr)
	if er != nil {
		return er
	}

	return nil
}

// Exists check whether a file is exists
func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// FileExists check whether a file is exists and is a file
func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

// FileExists check whether a dir is exists
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
