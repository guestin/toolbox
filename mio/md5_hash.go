package mio

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func Md5Hex(str string) string {
	return Md5Hex2([]byte(str))
}

func Md5Hex2(byteSlice []byte) string {
	md5Bytes := md5.Sum(byteSlice)
	return hex.EncodeToString(md5Bytes[:])
}

func Md5HexFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	return Md5HexFile2(file)
}

func Md5HexFile2(file *os.File) (string, error) {
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return "", err
	}
	result := md5Hash.Sum(nil)
	return hex.EncodeToString(result), nil
}

func Md5HexStream(reader io.Reader) (string, error) {
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, reader); err != nil {
		return "", err
	}
	result := md5Hash.Sum(nil)
	return hex.EncodeToString(result), nil
}

func Md5HexStream2(readSeeker io.ReadSeeker) (string, error) {
	restore, err := SaveSeekerPos(readSeeker)
	if err != nil {
		return "", err
	}
	defer restore()
	return Md5HexStream(readSeeker)
}
