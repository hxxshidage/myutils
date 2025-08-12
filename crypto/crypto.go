package ucry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"strings"
)

var hashAlgoMap = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

type hAsh struct {
}

var HASH = hAsh{}

func (h hAsh) Encrypt(plains []byte, algo string) string {
	if hashAlgo, ok := hashAlgoMap[algo]; !ok {
		panic(fmt.Sprintf("unkonown hash algo:%s", algo))
	} else {
		ha := hashAlgo()
		ha.Write(plains)
		return hex.EncodeToString(ha.Sum(nil))
	}
}

func (h hAsh) EncryptPlus(plain string, algo string) string {
	return h.Encrypt([]byte(plain), algo)
}

type hMAC struct {
}

var HMAC = hMAC{}

func (ha hMAC) Encrypt(keys, plains []byte, algo string) string {
	if hashAlgo, ok := hashAlgoMap[algo]; !ok {
		panic(fmt.Sprintf("unkonown hmac algo:%s", algo))
	} else {
		mac := hmac.New(hashAlgo, keys)
		mac.Write(plains)
		h := mac.Sum(nil)
		return hex.EncodeToString(h)
	}
}

func (ha hMAC) EncryptPlus(key, plain, algo string) string {
	return ha.Encrypt([]byte(key), []byte(plain), algo)
}

type b64 struct {
}

var B64 = b64{}

func (b b64) Encrypt(plains []byte) []byte {
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(plains)))
	enc.Encode(buf, plains)
	return buf
}

func (b b64) EncryptPlus(plain string) string {
	return string(b.Encrypt([]byte(plain)))
}

func (b b64) Decrypt(ciphers []byte) ([]byte, error) {
	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(ciphers)))
	n, err := enc.Decode(dbuf, ciphers)
	return dbuf[:n], err
}

func (b b64) DecryptPlus(cipher string) (string, error) {
	if plains, err := b.Decrypt([]byte(cipher)); err != nil {
		return "", err
	} else {
		return string(plains), nil
	}
}

type aesCbcP7 struct {
}

var AesCbcP7 = aesCbcP7{}

func (aesCbcP7) Encrypt(keys, ivs, plains []byte) ([]byte, error) {
	block, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}
	if len(ivs) == 0 {
		blockSize := block.BlockSize()
		ivs = keys[:blockSize]
	}

	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(plains, blockSize)
	shadow := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, ivs)
	blockMode.CryptBlocks(shadow, encryptBytes)

	return B64.Encrypt(shadow), nil
}

func (ac aesCbcP7) EncryptPlus(key, iv, plain string) (string, error) {
	if err := aesPreCheck(key, iv); err != nil {
		return "", err
	}

	if len(iv) != 16 {

	}

	if ciphers, err := ac.Encrypt([]byte(key), []byte(iv), []byte(plain)); err != nil {
		return "", err
	} else {
		return string(ciphers), nil
	}
}

func (aesCbcP7) Decrypt(keys, ivs, ciphers []byte) ([]byte, error) {
	ciphersDecode, err := B64.Decrypt(ciphers)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}

	if len(ivs) == 0 {
		blockSize := block.BlockSize()
		ivs = keys[:blockSize]
	}

	blockMode := cipher.NewCBCDecrypter(block, ivs)
	sunshine := make([]byte, len(ciphersDecode))
	blockMode.CryptBlocks(sunshine, ciphersDecode)
	sunshine = pkcs7UnPadding(sunshine)

	return sunshine, nil
}

func (ac aesCbcP7) DecryptPlus(key, iv, cipher string) (string, error) {
	if err := aesPreCheck(key, iv); err != nil {
		return "", err
	}

	if plains, err := ac.Decrypt([]byte(key), []byte(iv), []byte(cipher)); err != nil {
		return "", err
	} else {
		return string(plains), nil
	}
}

type aesEcbP7 struct {
}

var AesEcbP7 = aesEcbP7{}

func (aesEcbP7) Encrypt(keys, plains []byte) ([]byte, error) {
	block, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plains = pkcs7Padding(plains, blockSize)

	ciphertext := make([]byte, len(plains))

	for bs, be := 0, blockSize; bs < len(plains); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(ciphertext[bs:be], plains[bs:be])
	}

	return B64.Encrypt(ciphertext), nil
}

func (ae aesEcbP7) EncryptPlus(key, plain string) (string, error) {
	if err := aesPreCheck(key, ""); err != nil {
		return "", err
	}

	if ciphers, err := ae.Encrypt([]byte(key), []byte(plain)); err != nil {
		return "", err
	} else {
		return string(ciphers), nil
	}
}

func (ae aesEcbP7) Decrypt(keys, ciphers []byte) ([]byte, error) {
	block, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	plains := make([]byte, len(ciphers))

	// ECB 模式：每个块独立解密
	for bs, be := 0, blockSize; bs < len(ciphers); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(ciphers[bs:be], ciphers[bs:be])
	}

	// 去除填充
	plains = pkcs7UnPadding(ciphers)

	return plains, nil
}

func (ae aesEcbP7) DecryptPlus(key, cipher string) (string, error) {
	if err := aesPreCheck(key, ""); err != nil {
		return "", err
	}

	ciphers, err := B64.Decrypt([]byte(cipher))
	if err != nil {
		return "", err
	}
	if plains, err := ae.Decrypt([]byte(key), ciphers); err != nil {
		return "", err
	} else {
		return string(plains), nil
	}
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

func aesPreCheck(key, iv string) error {
	l := len(key)

	if l != 16 && l != 24 && l != 32 {
		return errors.New(fmt.Sprintf("aes key invalid, key:%s, len:%d", key, l))
	}

	if iv != "" && len(iv) != 16 {
		return errors.New(fmt.Sprintf("aes iv invalid, iv:%s, len:%d", iv, len(iv)))
	}

	return nil
}

// 实现了1024位和2048位的rsa加解密
type rSA struct {
}

var RSA = rSA{}

func (rSA) fmtKey(key string, isPub bool) string {
	if strings.Index(key, "-----") == -1 {
		if isPub {
			return fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", key)
		} else {
			return fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", key)
		}
	}
	return key
}

func (rSA) parsePubKey(pubKeys []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKeys)
	if block == nil {
		return nil, errors.New(fmt.Sprintf("public key error, key:%s", string(pubKeys)))
	}
	pk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pk.(*rsa.PublicKey), nil
}

func (rSA) parsePriKey(priKeys []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priKeys)
	if block == nil {
		return nil, errors.New(fmt.Sprintf("private key error, key:%s", string(priKeys)))
	}

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pk.(*rsa.PrivateKey), nil
}

func (rSA) _encrypt(pubKey *rsa.PublicKey, plains []byte) ([]byte, error) {
	keySize, srcSize := pubKey.Size(), len(plains)
	offSet, blockSize := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + blockSize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plains[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	return buffer.Bytes(), nil
}

func (rSA) _decrypt(priKey *rsa.PrivateKey, ciphers []byte) ([]byte, error) {
	blockSize := priKey.Size()

	srcSize := len(ciphers)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + blockSize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, ciphers[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	return buffer.Bytes(), nil
}

func (r rSA) Encrypt(pubKeys, plains []byte) ([]byte, error) {
	pubKey, err := r.parsePubKey(pubKeys)
	if err != nil {
		return nil, err
	}

	return r._encrypt(pubKey, plains)
}

func (r rSA) EncryptPlus(pubKey, plain string) (string, error) {
	pubKey = r.fmtKey(pubKey, true)

	if ciphers, err := r.Encrypt([]byte(pubKey), []byte(plain)); err != nil {
		return "", err
	} else {
		return string(B64.Encrypt(ciphers)), nil
	}
}

func (r rSA) Decrypt(priKeys, ciphers []byte) ([]byte, error) {
	//is1024 := r.is1024(len(priKeys))

	//var priKey *rsa.PrivateKey
	//if false {
	//	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	priKey = pk
	//} else {
	//	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	priKey = pk.(*rsa.PrivateKey)
	//}

	priKey, err := r.parsePriKey(priKeys)
	if err != nil {
		return nil, err
	}

	return r._decrypt(priKey, ciphers)

}

func (r rSA) DecryptPlus(priKey, plain string) (string, error) {
	priKey = r.fmtKey(priKey, false)

	b64ciphers, err := B64.Decrypt([]byte(plain))
	if err != nil {
		return "", err
	}

	if ciphers, err := r.Decrypt([]byte(priKey), b64ciphers); err != nil {
		return "", err
	} else {
		return string(ciphers), nil
	}
}
