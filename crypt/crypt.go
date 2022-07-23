package crypt

// 全局可逆加密

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

//
// Encrypt
//  @Description: 可逆 加密方法
//  @param text: 加密信息
//  @param key: 加密密码
//  @return string
//  @return error
//
func Encrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = zeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

//
// Decrypt
//  @Description: 可逆 解密方法
//  @param decrypted: 待解密字符串
//  @param key: 解密密码
//  @return string
//  @return error
//
func Decrypt(decrypted string, key []byte) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = zeroUnPadding(out)
	return string(out), nil
}

// 加密任务 密钥

// crypt internals
var (
	cryptKey = []byte{
		0x9c, 0x93, 0x5b, 0x48, 0x73, 0x0a, 0x55, 0x4d,
		0x6b, 0xfd, 0x7c, 0x63, 0xc8, 0x86, 0xa9, 0x2b,
		0xd3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x8a, 0xfb,
		0xf4, 0xde, 0x16, 0x2b, 0x8b, 0x95, 0xf6, 0x38,
	}
	cryptBlock cipher.Block
	cryptRand  = rand.Reader
)

func crypt(out, in, iv []byte) error {
	if cryptBlock == nil {
		var err error
		cryptBlock, err = aes.NewCipher(cryptKey)
		if err != nil {
			return err
		}
	}
	stream := cipher.NewCTR(cryptBlock, iv)
	stream.XORKeyStream(out, in)
	return nil
}

//
// Obscure
//  @Description: 加密任务生成配置文件中的密码
//  @param x
//  @return string
//  @return error
//
func Obscure(x string) (string, error) {
	plaintext := []byte(x)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(cryptRand, iv); err != nil {
		return "", fmt.Errorf("failed to read iv: %w", err)
	}
	if err := crypt(ciphertext[aes.BlockSize:], plaintext, iv); err != nil {
		return "", fmt.Errorf("encrypt failed: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func Reveal(x string) (string, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(x)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed when revealing password - is it obscured?: %w", err)
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("input too short when revealing password - is it obscured?")
	}
	buf := ciphertext[aes.BlockSize:]
	iv := ciphertext[:aes.BlockSize]
	if err := crypt(buf, buf, iv); err != nil {
		return "", fmt.Errorf("decrypt failed when revealing password - is it obscured?: %w", err)
	}
	return string(buf), nil
}
