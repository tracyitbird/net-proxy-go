package common

import (
	"../encrypt"
	"fmt"
)

type PackageHandler interface {
	Handle(pkg *Package) (newPkg Package)
}

//encrypt handler
type EncryptHandler struct {
	encrypt encrypt.Cipher
	iv []byte
	init bool
}

func (encryptHandler *EncryptHandler) Handle(pkg *Package) (newPkg Package) {
	header := pkg.GetHeader()
	body := pkg.GetBody()

	var nPkg Package
	var encryptHeader []byte
	encryptBody := make([]byte, len(body))

	encryptHandler.encrypt.Encrypt(encryptBody, body)

	if !encryptHandler.init {
		//init
		iv, err := encryptHandler.encrypt.InitEncrypt()
		if err != nil {
			fmt.Println("init encrypt handler error ...")
		} else {
			iv = make([]byte, 16)
		}
		encryptHandler.iv = iv
		encryptHeader = make([]byte, len(header) + 4 + len(iv))

		tmp := make([]byte, len(header))
		encryptHandler.encrypt.Encrypt(tmp, header)

		copy(encryptHeader[:4], IntToBytes(len(iv))[:])
		copy(encryptHeader[4:len(iv)], iv[:])
		copy(encryptHeader[4 + len(iv):], tmp[:])
		encryptHandler.init = true
	} else {
		encryptHeader = make([]byte, len(header))
		encryptHandler.encrypt.Encrypt(encryptHeader, header)
	}

	nPkg.ValueOf(encryptHeader, encryptBody)
	return nPkg
}

//decrypt handler
type DecryptHandler struct {
	decrypt encrypt.Cipher
	iv []byte
	init bool
}

func (decryptHandler *DecryptHandler) Handle(pkg *Package) (newPkg Package) {
	header := pkg.GetHeader()
	body := pkg.GetBody()

	var nPkg Package
	var decryptHeader []byte
	decryptBody := make([]byte, len(body))

	if !decryptHandler.init {
		lvLen := BytesToInt(header[:4])
		iv := header[4:4+lvLen]
		decryptHandler.decrypt.InitDecrypt(iv)

		decryptHeader = make([]byte, len(header) - 4 - len(iv))
		decryptHandler.decrypt.Decrypt(decryptHeader, header[4+len(iv):])
		decryptHandler.decrypt.Decrypt(decryptBody, body)

		nPkg = *NewPackage()
		nPkg.ValueOf(decryptHeader, decryptBody)
		decryptHandler.decrypt.Decrypt(decryptBody, body)

		decryptHandler.init = true
	} else {
		decryptHeader = make([]byte, len(header))

		decryptHandler.decrypt.Decrypt(decryptHeader, header[:])
		decryptHandler.decrypt.Decrypt(decryptBody, body)
	}

	nPkg.ValueOf(decryptHeader, decryptBody)
	return nPkg
}

type CompressHandler struct {

}

type DecompressHandler struct {
	
}
