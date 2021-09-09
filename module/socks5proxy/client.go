package socks5proxy

import (
	logs "github.com/danbai225/go-logs"
	"net"
	"sync/atomic"
)

type TcpClient struct {
	conn   *net.TCPConn
	server *net.TCPAddr
}

var count = int64(0)

func handleProxyRequest(localClient *net.TCPConn, serverAddr *net.TCPAddr, auth socks5Auth, recvHTTPProto string) {
	atomic.AddInt64(&count, 1)
	//logs.Info("start num:",atomic.LoadInt64(&count))
	defer atomic.AddInt64(&count, -1)
	// 远程连接IO
	dstServer, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		logs.Err("远程服务器地址连接错误!!!", err)
		return
	}
	go SecureCopy(dstServer, localClient, auth.Decrypt)
	SecureCopy(localClient, dstServer, auth.Encrypt)
	//logs.Info("end num:",atomic.LoadInt64(&count))
}

var listener *net.TCPListener

func client(listenAddrString string, serverAddrString string, encrytype string, passwd string, recvHTTPProto string) error {
	//所有客户服务端的流都加密,
	var err error
	auth, err := CreateAuth(encrytype, passwd)
	if err != nil {
		return err
	}
	logs.Info("你的密码是: %s ,请保管好你的密码", passwd)

	// proxy地址
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrString)
	if err != nil {
		return err
	}
	logs.Info("连接远程服务器: %s ....", serverAddrString)

	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrString)
	if err != nil {
		return err
	}
	logs.Info("监听本地端口: %s ", listenAddrString)

	listener, err = net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}

	for {
		localClient, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go handleProxyRequest(localClient, serverAddr, auth, recvHTTPProto)
	}
}
func close() {
	listener.Close()
}
