package file

import (
	"fmt"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"io"
	"os"
)

type File struct {
	Path string
}

func New() (*File, error) {
	config, err := conf.GetConf()
	if err != nil {
		return nil, err
	}

	return &File{Path:config.FileStorage.FilePath}, nil
}

func (f *File) Store(src io.ReadCloser, path, dstFileName string) (string, error) {
	pathRoot := fmt.Sprintf("%s%s", f.Path, path)
	if !pathExists(pathRoot) {
		err := os.Mkdir(pathRoot, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	file := fmt.Sprintf("%s/%s", pathRoot, dstFileName)

	newFile, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = newFile.Close()
	}()

	_, err = io.Copy(newFile, src)
	if err != nil {
		return "", err
	}

	return file, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}