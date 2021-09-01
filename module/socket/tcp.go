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
	t := tcpClient{}
	t.Host = client.Host
	t.Port = client.Port
	t.ws = client.ws
	t.Password = client.Password
	t.ProxyType = client.ProxyType
	t.ProxyHost = client.ProxyHost
	t.ProxyPort = client.ProxyPort
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
