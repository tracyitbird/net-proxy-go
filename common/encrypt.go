package common

import (
	"crypto/cipher"
	"crypto/md5"
	"io"
	"crypto/rand"
	"crypto/aes"
	"errors"
)

type Cipher struct {
	enc cipher.Stream
	dec cipher.Stream
	key []byte
	iv [] byte
}

func NewCipher(password string) (cipher *Cipher, err error) {
	return &Cipher{key:generateKeyUsePassword(password, 32)}, nil
}

func md5Sum(src []byte) []byte {
	h := md5.New()
	h.Write(src)
	return h.Sum(nil)
}

func generateKeyUsePassword(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5Sum([]byte(password)))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5Sum(d))
	}
	return m[:keyLen]
}

// Initializes the block cipher with CFB mode, returns IV.
func (c *Cipher) initEncrypt() (err error) {
	if c.iv == nil {
		c.iv = make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		c.iv = iv
	} else {
		iv = c.iv
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return err
	}
	c.enc = cipher.NewCTR(block, iv)
	return errors.New("init encrypt error ...")
}

func (c *Cipher) initDecrypt(iv []byte) (err error) {
	c.dec, err = c.info.newStream(c.key, iv, Decrypt)
	return
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}
