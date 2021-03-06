package main

import (
	"net"
	"os"

	//"github.com/villcore/net-proxy-go/server"
	"../server"
	"../conf"
	"fmt"
	"sync"
	"log"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	serverConfig, err := conf.ReadServerConf("server.conf")
	if err != nil {
		fmt.Println("can not load conf file ...")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(serverConfig.PortPair))

	for _, portAndPair := range serverConfig.PortPair {
		go startListen(portAndPair)
	}
	wg.Wait()
}

func startListen(portAndPassword conf.PortAndPassword) {
	fmt.Println(portAndPassword)
	listenPort := portAndPassword.ListenPort
	listenAddrAndPort := ":" + listenPort
	password := portAndPassword.Password
	fmt.Println("start listen port ", listenPort)
	log.Println("server start...")
	listener, err := net.Listen("tcp", listenAddrAndPort)
	if err != nil {
		log.Println("erver starting listen failed at port [%v] ...", listenPort)
	}

	addr := listener.Addr()
	log.Println("server staring listen address : [%v] ...", addr.String())

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Println("accept conn [%v] failed ...\n", localConn.LocalAddr())
		}
		log.Println("accept conn [%v] success ...\n", localConn.RemoteAddr())

		go server.AcceptConn(localConn, password)
	}
}
