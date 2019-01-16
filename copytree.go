// Package goshutil provides similar functionality to
// https://docs.python.org/3.7/library/shutil.html#shutil.copytree
package goshutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File implements file object interface.
type File struct{}

// New returns a File object.
func New() *File {
	return &File{}
}

// copyFile copes contents of the file named src to a file named dst.
func (f *File) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			err = cerr
		}
		if cerr := in.Close(); cerr != nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}

// CopyTree recursively creates and populates directory tree.
// similar to python shutil.copytree
func (f *File) CopyTree(src, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.Mkdir(dst, os.ModePerm); err != nil {
		return err
	}
	for _, fh := range files {
		srcName := filepath.Join(src, fh.Name())
		dstName := filepath.Join(dst, fh.Name())
		fileInfo, err := os.Stat(srcName)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if err := f.CopyTree(srcName, dstName); err != nil {
				return err
			}
		} else {
			if err := f.copyFile(srcName, dstName); err != nil {
				return err
			}
		}
	}
	return nil
}
