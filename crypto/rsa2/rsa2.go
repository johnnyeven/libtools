/*
for alipay
*/
package rsa2

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	privateKeyHead = "-----BEGIN RSA PRIVATE KEY-----"
	privateKeyTail = "-----END RSA PRIVATE KEY-----"
	publicKeyHead  = "-----BEGIN PUBLIC KEY-----"
	publicKeyTail  = "-----END PUBLIC KEY-----"
	interval       = 64
)

func FormatPrivateKey(src string) string {
	if src == "" {
		return ""
	}
	r := strings.Replace(strings.Replace(src, privateKeyHead, "", -1), privateKeyTail, "", -1)
	r = strings.TrimSpace(r)
	r = strings.Replace(r, "\n", "", -1)
	a := make([]string, 0)

	for index := 0; index < len(r); index = index + interval {
		if index+interval < len(r) {
			a = append(a, r[index:index+interval])
		} else {
			a = append(a, r[index:])
		}
	}

	m := make([]string, 0)
	m = append(m, privateKeyHead)
	m = append(m, a[:]...)
	m = append(m, privateKeyTail)
	return strings.Join(m, "\n")
}

func FormatPublicKey(src string) string {
	if src == "" {
		return ""
	}
	r := strings.Replace(strings.Replace(src, publicKeyHead, "", -1), publicKeyTail, "", -1)
	r = strings.TrimSpace(r)
	r = strings.Replace(r, "\n", "", -1)
	a := make([]string, 0)

	for index := 0; index < len(r); index = index + interval {
		if index+interval < len(r) {
			a = append(a, r[index:index+interval])
		} else {
			a = append(a, r[index:])
		}
	}

	m := make([]string, 0)
	m = append(m, publicKeyHead)
	m = append(m, a[:]...)
	m = append(m, publicKeyTail)
	return strings.Join(m, "\n")
}

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

func RSAEncrypt(plaintext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	var data = packageData(plaintext, pub.N.BitLen()/8-11)
	var cipherData []byte = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return nil, e
		}
		cipherData = append(cipherData, c...)
	}

	return cipherData, nil
}

func RSADecrypt(ciphertext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pri = prkI.(*rsa.PrivateKey)
	}

	var data = packageData(ciphertext, pri.PublicKey.N.BitLen()/8)
	var plainData []byte = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func SignPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)

	// for java
	if err != nil {
		prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pri = prkI.(*rsa.PrivateKey)
	}
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

func VerifyPKCS1v15(src, sig, key []byte, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return errors.New("publick key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}
