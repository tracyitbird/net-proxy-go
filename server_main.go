package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"../net-proxy-go/server"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	listenPort := "60081"
	listenAddrAndPort := ":" + listenPort
	fmt.Print("server start...\n")

	listener, err := net.Listen("tcp", listenAddrAndPort)
	if err != nil {
		log.Printf("server starting listen failed at port [%v] ...\n", listenPort)
	}

	addr := listener.Addr()
	log.Printf("server staring listen address : [%v] ...\n", addr.String())

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("accept conn [%v] failed ...\n", localConn.LocalAddr())
		}
		log.Printf("accept conn [%v] success ...\n", localConn.RemoteAddr())

		go server.AcceptConn(localConn)
	}
}
