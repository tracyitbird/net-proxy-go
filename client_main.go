package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"../net-proxy-go/client"
	"strconv"
)

const(
	HTTP = 1;
	HTTPS = 2;
	SOCKS_5 = 3
)
func main() {
	listenPort := 50081
	fmt.Print("local client start...\n")
	//
	log.SetOutput(os.Stdout)

	listenAddr := ":" + strconv.Itoa(listenPort)
	log.Printf("[%v]", listenAddr)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Printf("starting listen failed at port [%v] ...\n", listenPort)
		log.Println(err)
		return
	}

	addr := listener.Addr()
	log.Printf("staring listen address : [%v] ...\n", addr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept conn [%v] failed ...\n", conn.LocalAddr())
		}
		log.Printf("accept conn [%v] success ...\n", conn.RemoteAddr())

		go client.AcceptConn(conn)
	}
}
