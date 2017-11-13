package client

import (
	"net"
	"sync"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/villcore/net-proxy-go/common"
	//"../common"
)

func init() {
	log.SetOutput(os.Stdout)
}

//1.接受本地连接
//2.构建远程连接,(可用连接连接池)
//3.循环转发(接受包 -> handler处理 -> 发送)
//4.错误处理
func AcceptConn(localConn net.Conn, remoteAddr string, remotePort string) {
	var wg sync.WaitGroup
	wg.Add(2)

	remoteConn, error := common.GetRemoteConn(remoteAddr, remotePort)

	if error != nil {
		if localConn != nil {
			localConn.Close()
		}
		if remoteConn != nil {
			remoteConn.Close()
		}
		log.Printf("build conn to remote [%v:%v] failed ...", remoteAddr, remotePort)
		wg.Done()
		wg.Done()
	} else {
		//transfer
		go common.TransferBytesToPackage(localConn, remoteConn, make([]common.PackageHandler, 0), &wg)

		//transfer
		go common.TransferPackageToBytes(remoteConn, localConn, make([]common.PackageHandler, 0), &wg)
	}

	wg.Wait()
	log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!tunnel close ...")
	defer func() {
		if localConn != nil {
			localConn.Close()
		}
		if remoteConn != nil {
			remoteConn.Close()
		}
		log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tunnel close ...")
	}()
}
