package socket

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/net/ghttp"
)

//code
const (
	TcpConn       = 11001
	TcpDataWrite  = 11002
	TcpDataRead   = 11003
	TCPDisconnect = 11004
	Ok            = 10001
	Err           = 20001
)

type Msg struct {
	Tag  string `json:"tag"`
	Type int    `json:"type"`
	Data string `json:"data"`
}
type connApi interface {
	write([]byte) error
	close()
}
type wsClient struct {
	ws      *ghttp.WebSocket
	clients map[int]connApi
}

func (wsClient) New(ws *ghttp.WebSocket) *wsClient {
	w := wsClient{ws: ws, clients: make(map[int]connApi)}
	go w.handle()
	return &w
}
func (w *wsClient) handle() {
	for {
		_, data, err := w.ws.ReadMessage()
		if err != nil {
			logs.Err(err)
			break
		}
		msg := Msg{}
		json.Unmarshal(data, &msg)

		switch msg.Type {
		//建立tcp连接
		case TcpConn:
			c := client{}
			json.Unmarshal(data, &c)
			c.ws = w.ws
			t := tcpClient{}.New(&c)
			err = t.Conn()
			if err != nil {
				w.ws.WriteJSON(Msg{
					Tag:  "tcp",
					Type: Err,
					Data: err.Error(),
				})
			} else {
				w.ws.WriteJSON(Msg{
					Tag:  "tcp",
					Type: Ok,
				})
			}
			w.clients[TcpDataWrite] = t
		case TCPDisconnect:
			if w.clients[msg.Type] != nil {
				w.clients[msg.Type].close()
			}
		default:
			if c, has := w.clients[msg.Type]; has {
				err := c.write([]byte(msg.Data))
				if err != nil {
					logs.Err(err)
				}
			}
		}

	}
	for _, f := range w.clients {
		f.close()
	}
}

type client struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	ws   *ghttp.WebSocket
}
