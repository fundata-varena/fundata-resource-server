package storage

import (
	"errors"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/storage/file"
	"github.com/fundata-varena/fundata-resource-server/storage/oss"
	"io"
)

type Storage interface {
	Store(io.ReadCloser, string, string) (string, error)
}

var instance Storage

func Init() error {
	config, err := conf.GetConf()
	if err != nil {
		return err
	}

	storageDriver := config.StorageUsing

	if storageDriver == "file" {
		instance, err = file.New()
		if err != nil {
			return err
		}
	} else if storageDriver == "oss" {
		// 还未实现的
		instance = new(oss.Oss)
	} else {
		return errors.New("storage driver doesn't support")
	}

	return nil
}

func GetInstance() (Storage, error) {
	if instance == nil {
		return nil, errors.New("init Storage first please")
	}
	return instance, nil
}
