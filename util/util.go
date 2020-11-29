package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (sha1Stream *Sha1Stream) Update(data []byte) {
	if sha1Stream._sha1 == nil {
		sha1Stream._sha1 = sha1.New()
	}
	sha1Stream._sha1.Write(data)
}

func (sha1Stream *Sha1Stream) Sum() string {
	return hex.EncodeToString(sha1Stream._sha1.Sum(nil))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum(nil))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Open(path)
	if err == nil {
		return true, nil
	}

	if os.ErrNotExist == err {
		return false, nil
	}
	return false, err
}

//GetFileSize 获取文件长度
func GetFileSize(fileName string) int64 {
	var length int64
	filepath.Walk(fileName, func(path string, f os.FileInfo, err error) error {
		length = f.Size()
		return nil
	})
	return length
}
