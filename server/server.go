package server

import (
	"fmt"
	"log"
	"os"
	"net"
	"strings"
)



func acceptConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	//100kb
	buf := make([]byte, 1024 * 100)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("read bytes form conn %v failed...\n", conn.RemoteAddr())
	}

	log.Printf("read %v bytes form conn %v ...\n", n, conn.RemoteAddr())
	req := string(buf[:n])
	log.Print(len([]byte(req)), "\n")
	//log.Printf("read content = \n%v ...\n", req)

	proctol := parseProctol(buf, n)
	switch proctol {
	case HTTP:
		log.Println("http proctol...")
		break
	case HTTPS:
		log.Println("https proctol...")
		break
	case SOCKS_5:
		log.Println("socks_5 proctol...")
		break
	default:
		log.Println("unrecognized proctol...")
		return
	}
}

func parseProctol(req []byte, len int) (int) {
	//TODO SOCKS_5

	//HTTPS
	headerInfo := string(req[:7])
	if strings.EqualFold("CONNECT", headerInfo) {
		return HTTPS
	}

	//HTTP
	httpOpePos := -1
	for index, val := range req {
		if index >= len {
			break
		}
		//fmt.Println(val, int(val))
		if int(val) == 32 {
			httpOpePos = index
			break
		}
	}

	if httpOpePos != -1 {
		requestMethod := string(req[:httpOpePos])
		log.Println("requestMethod = ", requestMethod)
		if strings.EqualFold(requestMethod, "GET") || strings.EqualFold(requestMethod, "POST") {
			return  HTTP
		}
	}

	return -1
}
