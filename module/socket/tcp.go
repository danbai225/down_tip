package socket

import (
	"fmt"
	"net"
)

type tcpClient struct {
	client
	conn  net.Conn
	proxy *proxy
}

func (tcpClient) New(client *client) (*tcpClient, error) {
	t := tcpClient{
		client: *client,
	}
	//t.Host = client.Host
	//t.Port = client.Port
	//t.ws = client.ws
	//t.Password = client.Password
	//t.ProxyType = client.ProxyType
	//t.ProxyHost = client.ProxyHost
	//t.ProxyPort = client.ProxyPort
	if t.ProxyType != "" {
		p, err := proxy{}.New(t.ProxyHost, t.ProxyPort, t.ProxyType, t.UserName, t.Password)
		t.proxy = p
		if err != nil {
			return &t, err
		}
	}
	err := t.initConn()
	return &t, err
}
func (tcp *tcpClient) initConn() error {
	if tcp.conn == nil {
		var err error
		add := fmt.Sprintf("%s:%d", tcp.Host, tcp.Port)
		if tcp.proxy == nil {
			tcp.conn, err = net.Dial("tcp", add)
		} else {
			tcp.conn, err = tcp.proxy.proxyBySocks5("tcp", add)
		}
		if err != nil {
			return err
		}
		go tcp.handle()
		return err
	}
	return nil
}
func (tcp *tcpClient) write(data []byte) error {
	_, err := tcp.conn.Write(data)
	return err
}
func (tcp *tcpClient) handle() {
	defer tcp.close()
	bytes := make([]byte, 1024)
	for {
		l, err := tcp.conn.Read(bytes)
		if err != nil {
			_ = tcp.ws.WriteJSON(Msg{
				Tag:  "tcp",
				Type: Err,
				Data: err.Error(),
			})
			return
		}
		err = tcp.ws.WriteJSON(Msg{
			Tag:  "tcp",
			Type: TcpDataRead,
			Data: string(bytes[:l]),
		})
		if err != nil {
			return
		}
	}
}
func (tcp *tcpClient) close() {
	if tcp.conn != nil {
		_ = tcp.conn.Close()
	}
}

type tcpServer struct {
	server
	listener net.Listener
	connMap  map[string]net.Conn
}

func (tcpServer) New(server *server) (*tcpServer, error) {
	t := tcpServer{
		server:   *server,
		listener: nil,
		connMap:  make(map[string]net.Conn),
	}
	err := t.initListener()
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (l *tcpServer) initListener() error {
	add := fmt.Sprintf("%s:%d", l.Host, l.Port)
	listen, err := net.Listen("tcp", add)
	l.listener = listen
	return err
}

//func (l *tcpServer) write(data []byte) error {
//	_, err := l.listener.Write(data)
//	return err
//}
func (l *tcpServer) handle() {
	defer l.close()
	for {
		accept, err := l.listener.Accept()
		if err != nil {
			_ = l.ws.WriteJSON(Msg{
				Tag:  "tcp",
				Type: TCPServerErr,
				Data: err.Error(),
			})
			return
		}
		l.connMap[accept.RemoteAddr().String()] = accept
		go l.handleClient(accept)
		err = l.ws.WriteJSON(Msg{
			Tag:  "tcp",
			Type: TCPServerNewConn,
			Data: accept.RemoteAddr().String(),
		})
		if err != nil {
			return
		}
	}
}
func (l *tcpServer) handleClient(c net.Conn) {
	defer func() {
		_ = c.Close()
	}()
	bytes := make([]byte, 1024)
	for {
		read, err := c.Read(bytes)
		if err != nil {
			_ = l.ws.WriteJSON(Msg{
				Tag:  "tcp",
				Type: TCPServerErr,
				Data: err.Error(),
			})
			return
		}
		err = l.ws.WriteJSON(Msg{
			Mete: c.RemoteAddr().String(),
			Tag:  "tcp",
			Type: TCPServerNewMsg,
			Data: string(bytes[:read]),
		})
		if err != nil {
			return
		}
	}
}
func (l *tcpServer) close() {
	if l.listener != nil {
		_ = l.listener.Close()
	}
	for _, conn := range l.connMap {
		if conn != nil {
			_ = conn.Close()
		}
	}
}
