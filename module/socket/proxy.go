package socket

import (
	"errors"
	"fmt"
	"github.com/txthinking/socks5"
	"net"
)

type proxy struct {
	Host         string `json:"host"`
	Port         uint16 `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Type         string `json:"type"`
	socks5Client *socks5.Client
}

func (proxy) New(host string, port uint16, Type, user, pass string) (*proxy, error) {
	p := proxy{
		Host:     host,
		Port:     port,
		Type:     Type,
		Username: user,
		Password: pass,
	}
	var err error
	switch Type {
	case "socks5":
		err = p.initSocks5()
	}
	return &p, err
}
func (p *proxy) initSocks5() error {
	if p.Type != "socks5" {
		return errors.New("type is not socks5")
	}
	c, err := socks5.NewClient(fmt.Sprintf("%s:%d", p.Username, p.Port), p.Username, p.Password, 0, 60)
	if err != nil {
		return err
	}
	p.socks5Client = c
	return err
}
func (p *proxy) proxyBySocks5(netWork, add string) (net.Conn, error) {
	if p.Type != "socks5" {
		return nil, errors.New("type is not socks5")
	}
	if p.socks5Client == nil {
		return nil, errors.New("socks5 conn err")
	}
	return p.socks5Client.Dial(netWork, add)
}
func (p *proxy) close() {
	_ = p.socks5Client.Close()
}
