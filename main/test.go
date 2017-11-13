package main

import (
	"fmt"
)

func main() {
	firstLine := "CONNECT s0.ssl.qhres.com:443 HTTP/1.1"
	fmt.Println(firstLine)

	var firstBlank int
	var lastBlank int
	firstBlank = -1
	lastBlank = -1

	for index, value := range firstLine {
		//fmt.Println(value)
		if int(value) == 32 {
			if firstBlank < 0 {
				firstBlank = index
			}

			if firstBlank > 0 {
				lastBlank = index
			}
		}
	}

	fmt.Println(firstBlank, lastBlank)
	//addressAndPort := []byte(firstLine)[firstBlank + 1:lastBlank]
}
