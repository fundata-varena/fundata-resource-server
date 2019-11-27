package oss

import "io"

type Config4Oss struct {

}

type Oss struct {

}

func (o *Oss) Store(src io.ReadCloser, path, dstFileName string) (string, error) {
	return "", nil
}
