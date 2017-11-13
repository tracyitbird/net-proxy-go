package main

import (
	"net"
	"os"

	_ "net/http/pprof"
	log "github.com/sirupsen/logrus"
	"github.com/villcore/net-proxy-go/server"
	//"../server"
	"net/http"
	"strconv"
	"runtime"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	go func() {
		http.ListenAndServe("localhost:7001", nil)
	}()

	go func() {
		http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		});
		http.ListenAndServe("localhost:7002", nil)
	}()

	listenPort := "60081"
	listenAddrAndPort := ":" + listenPort

	log.WithField("fname", "server_main").Info("server start...")

	listener, err := net.Listen("tcp", listenAddrAndPort)
	if err != nil {
		log.WithField("fname", "server_main").Info("erver starting listen failed at port [%v] ...", listenPort)
	}

	addr := listener.Addr()
	log.WithField("fname", "server_main").Info("server staring listen address : [%v] ...", addr.String())

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.WithField("fname", "server_main").Info("accept conn [%v] failed ...\n", localConn.LocalAddr())
		}
		log.WithField("fname", "server_main").Info("accept conn [%v] success ...\n", localConn.RemoteAddr())

		go server.AcceptConn(localConn)
	}
}
