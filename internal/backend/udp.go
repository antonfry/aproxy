package backend

import (
	"fmt"
	"net"
	"sync/atomic"
)

type UDPBackend struct {
	Host  string
	Port  string
	Conn  *net.UDPConn
	Alive atomic.Bool
}

func New(host, port string) *UDPBackend {
	return &UDPBackend{
		Host:  host,
		Port:  port,
		Conn:  nil,
		Alive: atomic.Bool{},
	}
}

func (b *UDPBackend) String() string {
	return fmt.Sprint("Host: ", b.Host)
}

func (b *UDPBackend) GetConn() error {
	host, err := net.ResolveUDPAddr("udp", net.JoinHostPort(b.Host, b.Port))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, host)
	if err != nil {
		return err
	}
	b.Conn = conn
	return nil
}

func (b *UDPBackend) SetAlive() {
	b.Alive.Store(true)
}

func (b *UDPBackend) SetDead() {
	b.Alive.Store(false)
}

func (b *UDPBackend) IsAlive() bool {
	return b.Alive.Load()
}

func (b *UDPBackend) Close() error {
	return b.Conn.Close()
}
