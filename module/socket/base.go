package socket

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/net/ghttp"
)

//code
const (
	TcpConn          = 11001
	TcpDataWrite     = 11002
	TcpDataRead      = 11003
	TCPDisconnect    = 11004
	TCPServerNewConn = 11101
	TCPServerNewMsg  = 11102
	Ok               = 10001
	Err              = 20001
	TCPClientErr     = 20002
	TCPServerErr     = 20002
)

type Msg struct {
	Tag  string      `json:"tag"`
	Type int         `json:"type"`
	Mete interface{} `json:"mete"`
	Data string      `json:"data"`
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
func (w *wsClient) sendErrMsg(tag, data string, Type ...int) {
	t := Err
	if len(Type) > 1 {
		t = Type[0]
	}
	w.ws.WriteJSON(Msg{
		Tag:  tag,
		Type: t,
		Data: data,
	})
}
func (w *wsClient) handle() {
	for {
		_, data, err := w.ws.ReadMessage()
		if err != nil {
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
			t, err2 := tcpClient{}.New(&c)
			if err2 != nil {
				w.sendErrMsg("tcp", err2.Error())
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
	Host      string `json:"host"`
	Port      uint16 `json:"port"`
	ProxyType string `json:"proxy_type"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	ProxyHost string `json:"proxy_host"`
	ProxyPort uint16 `json:"proxy_port"`
	ws        *ghttp.WebSocket
}
type server struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	ws   *ghttp.WebSocket
}
