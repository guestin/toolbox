package mod

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5Stream(fStream io.Reader) (md5Str string, err error) {
	hash := md5.New()
	if _, err := io.Copy(hash, fStream); err != nil {
		return "", err
	}
	md5Bytes := hash.Sum(nil)
	return hex.EncodeToString(md5Bytes[:]), nil
}
