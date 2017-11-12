package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	listenPort := 20080
	fmt.Print("server start...\n")
	log.SetOutput(os.Stdout)

	listener, err := net.Listen("tcp", ":20080")
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

		//go acceptConn(localConn)
	}
}
