package socket

import (
	"fmt"
	"net"
)

type tcpClient struct {
	client
	conn net.Conn
}

func (tcpClient) New(client *client) *tcpClient {
	t := tcpClient{}
	t.Host = client.Host
	t.Port = client.Port
	t.ws = client.ws
	return &t
}
func (tcp *tcpClient) Conn() error {
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", tcp.Host, tcp.Port))
	if err != nil {
		return err
	}
	tcp.conn = dial
	go tcp.handle()
	return err
}
func (tcp *tcpClient) write(data []byte) error {
	_, err := tcp.conn.Write(data)
	return err
}
func (tcp *tcpClient) handle() {
	bytes := make([]byte, 1024)
	for {
		l, err := tcp.conn.Read(bytes)
		if err != nil {
			tcp.ws.WriteJSON(Msg{
				Tag:  "tcp",
				Type: Err,
				Data: err.Error(),
			})
			break
		}
		tcp.ws.WriteJSON(Msg{
			Tag:  "tcp",
			Type: TcpDataRead,
			Data: string(bytes[:l]),
		})
	}
}
func (tcp *tcpClient) close() {
	if tcp.conn != nil {
		tcp.conn.Close()
	}
}
