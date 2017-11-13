package server

import (
	"net"
	"strings"
	"sync"
	"errors"
	"strconv"
	"fmt"
	"net/http"
	"bufio"

	log "github.com/sirupsen/logrus"
	//"github.com/villcore/net-proxy-go/common"
	"../common"
)

const (
	HTTP = 1
	HTTPS = 2
	SOCKS_5 = 3
)

//1.接受本地连接
//2.解析包, 解析协议, 解析目的地址
//3.构建远程连接,(可用连接连接池)

//4.循环转发(接受包 -> handler处理 -> 发送)
//5.错误处理
func AcceptConn(localConn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	bytesToPackageHandlers := make([]common.PackageHandler, 0)
	packageToBytesHandlers := make([]common.PackageHandler, 0)

	protocalDetected := false
	interrupt := false

	var buf []byte
	//var read int
	var protocal int = -1

	if !interrupt || !protocalDetected {
		pkg := *common.NewPackage()
		err := pkg.ReadWithHeader(localConn)

		if err != nil {
			log.Printf("read bytes form conn %v failed...\n", localConn.RemoteAddr())
			interrupt = true
		}

		for _, handler := range packageToBytesHandlers {
			pkg = handler.Handle(&pkg)
		}

		//detect protocal

		body := pkg.GetBody()
		buf = body
		log.Printf("server recv first pkg = \n%v\n", string(body))
		protocal = parseProtocal(body, len(body))

		switch protocal {
		case HTTP:
			log.Println("http protocal...")
			break
		case HTTPS:
			log.Println("https protocal...")
			break
		case SOCKS_5:
			log.Println("socks_5 protocal...")
			break
		default:
			log.Println("unrecognized protocal...")
			wg.Done()
			wg.Done()
			return
		}

		if protocal > 0 {
			protocalDetected = true
		}
	}

	if !protocalDetected {
		wg.Done()
		wg.Done()
		localConn.Close()
		return
	}

	addr, port, err := parseAddressAndPort(buf, protocal, localConn)
	if err != nil {
		wg.Done()
		wg.Done()
		localConn.Close()
		return
	}

	log.Printf("need connect to server [%v:%v]", addr, port)
	remoteConn, error := common.NewRemoteConn(addr, strconv.Itoa(port))

	if error != nil {
		wg.Done()
		wg.Done()
		log.Printf("build conn to remote [%v:%v] failed ...", addr, port)
		return
	}

	log.Printf("build conn to remote [%v:%v] success ...", addr, port)

	if(protocal == HTTPS) {
		//TODO bytesToPackageHandlers handle
		httpsConnectResp := "HTTP/1.0 200 Connection Established\r\n\r\n";

		httpsRespPkg := *common.NewPackage()
		httpsRespPkg.ValueOf(make([]byte, 0), []byte(httpsConnectResp))

		for _, handler := range bytesToPackageHandlers {
			httpsRespPkg = handler.Handle(&httpsRespPkg)
		}

		localConn.Write(httpsRespPkg.ToBytes())
	}

	if(protocal == HTTP) {
		remoteConn.Write(buf)
	}

	if(protocal == SOCKS_5) {
	}

	//transfer
	go common.TransferPackageToBytes(localConn, remoteConn, make([]common.PackageHandler, 0), wg)

	//transfer
	go common.TransferBytesToPackage(remoteConn, localConn, make([]common.PackageHandler, 0), wg)

	defer func() {
		if localConn != nil {
			localConn.Close()
		}

		if remoteConn != nil {
			remoteConn.Close()
		}
		log.Printf("tunnel close ...")
	}()

	wg.Wait()
}

func parseProtocal(req []byte, len int) int {
	//TODO SOCKS_5

	//HTTPS
	if len >= 7 {
		headerInfo := string(req[:7])
		if strings.EqualFold("CONNECT", headerInfo) {
			return HTTPS
		}
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
			return HTTP
		}
	}

	return -1
}

func parseAddressAndPort(firstReq []byte, protocal int, localConn net.Conn) (addr string, port int, err error) {
	//log.Printf("read content = \n%v ...\n", string(firstReq))
	switch protocal {
	case HTTP:
		return parseHttpAddress(firstReq)
		break
	case HTTPS:
		return parseHttpsAddress(firstReq)
		break
	case SOCKS_5:
		break
	default:
		log.Println("unrecognized proctol ...")
	}
	return "", -1, errors.New("unrecognized proctol...")
}

func parseHttpsAddress(firstReq []byte) (addr string, port int, err error) {
	//fmt.Printf("first line = \n%v\n", string(firstReq))
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(firstReq))))

	if err != nil {
		return "", -1, errors.New("unrecognized proctol...")
	}

	addrAndPort := req.Host
	infos := strings.Split(addrAndPort, ":")

	addr = infos[0]
	port = 443

	if(len(infos) > 1) {
		port, err = strconv.Atoi(infos[1])
	}

	fmt.Printf("addr = %v, port = %v\n", addr, port)
	return addr, port, nil
}

func parseHttpAddress(firstReq []byte) (addr string, port int, err error) {
	//fmt.Printf("first line = \n%v\n", string(firstReq))
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(firstReq))))

	if err != nil {
		return "", -1, errors.New("unrecognized proctol...")
	}

	addrAndPort := req.Host
	infos := strings.Split(addrAndPort, ":")

	addr = infos[0]
	port = 80

	if(len(infos) > 1) {
		port, err = strconv.Atoi(infos[1])
	}

	fmt.Printf("addr = %v, port = %v\n", addr, port)
	return addr, port, nil
}
