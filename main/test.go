package main

import (
	"../encrypt"
	"fmt"
)

func sliceModify(slice *[]int) {
	// slice[0] = 88
	*slice = append(*slice, 6)
}

func main() {
	slice := []int{1, 2, 3, 4, 5}
	sliceModify(&slice)
	fmt.Println(slice)

	testStr := "hello this is a test str ..."



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
