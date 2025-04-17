package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func GetFileMD5(fileHeader *multipart.FileHeader) (*string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 计算 MD5
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	// 转换为字符串
	var result *string
	md5String := hex.EncodeToString(hash.Sum(nil))
	result = &md5String
	return result, nil
}
