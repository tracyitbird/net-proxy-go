package main

import (
	"../encrypt"
	"fmt"
	"../conf"
)

func sliceModify(slice *[]int) {
	// slice[0] = 88
	*slice = append(*slice, 6)
}

func main() {
	slice := []int{1, 2, 3, 4, 5}
	sliceModify(&slice)
	fmt.Println(slice)

	testStr := "CONNECT www.baidu.com:443 HTTP/1.1\r\nHost: www.baidu.com:443\r\nUser-Agent: curl/7.53.1\r\nProxy-Connection: Keep-Alive\r\n\r\n"

	cipher, err := encrypt.NewCipher("villcore")
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

	//json
	clientConfig, _:= conf.ReadClientConf("client.conf")
	fmt.Println(clientConfig)

	serverConfig, _ := conf.ReadServerConf("server.conf")
	fmt.Println(len(serverConfig.PortPair))
}
