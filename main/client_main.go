package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	_ "net/http/pprof"
	log "github.com/sirupsen/logrus"
	//"github.com/villcore/net-proxy-go/client"
	"../client"
	"net/http"
	"runtime"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	go func() {
		http.ListenAndServe("localhost:6001", nil)
	}()

	go func() {
		http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		});
		http.ListenAndServe("localhost:6002", nil)
	}()

	listenPort := 50081

	remoteAddr := "127.0.0.1"
	remotePort := "60081"

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

		go client.AcceptConn(conn, remoteAddr, remotePort)
	}
}
