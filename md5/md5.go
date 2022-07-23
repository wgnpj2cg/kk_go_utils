package md5

import (
	"crypto/md5"
	"encoding/hex"
)

//
// GetMD5Encode
//  @Description: 获取md5
//  @param data
//  @return string
//
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//
// GetMD5EncodeWithSalt
//  @Description: 加盐md5算法
//  @param data: 待加密信息
//  @return string
//
func GetMD5EncodeWithSalt(data string) string {
	return GetMD5Encode("randomstrings" + data)
}
