package main

import (
	"../encrypt"
	"fmt"
)

func main() {
	testStr := "hello this is a test str ..."

	cipher, err := encrypt.NewCipher("villcore")

	if err != nil {
		fmt.Println("init cipher error ...")
	}

	iv, err := cipher.InitEncrypt()
	if err != nil {
		fmt.Println("init cipher error ...")
	}

	cipher.InitDecrypt(iv)
	if err != nil {
		fmt.Println("init cipher error ...")
	}

	eBytes := make([]byte, len(testStr))
	dBytes := make([]byte, len(testStr))

	cipher.Encrypt(eBytes, []byte(testStr))
	fmt.Println(string(eBytes))

	cipher.Decrypt(dBytes, eBytes)

	fmt.Println(string(dBytes))

}
