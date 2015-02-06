package coms

import (
	"errors"
	//"fmt"
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	//"io"
	"net"
	"time"
)

type gantryTCPClient struct {
	RemoteAddr string
	RemotePort string
	conn       *net.TCPConn
}

type gantryUDPClient struct {
	RemoteAddr string
	RemotePort string
	conn       *gantryUDPConn
}

type gantryUDPConn struct {
	conn *net.UDPConn
}
type gantryTCPConn struct {
	conn *net.TCPConn
}

func NewGantryTCPClient(ip, port string) *gantryTCPClient {
	return &gantryTCPClient{
		RemoteAddr: ip,
		RemotePort: port,
	}
}

func NewGantryUDPClient(ip, port string) *gantryUDPClient {
	return &gantryUDPClient{
		RemoteAddr: ip,
		RemotePort: port,
	}
}

func (client *gantryTCPClient) Connect() error {

	addr, err := net.ResolveTCPAddr("tcp4", client.RemoteAddr+":"+client.RemotePort)

	if err != nil {
		return err
	}

	if tcpConn, err := net.DialTCP("tcp4", nil, addr); err != nil { //.DialUDP("udp", nil, addr); err != nil {
		return err
	} else {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(5 * time.Second)
		tcpConn.SetLinger(5)
		tcpConn.SetWriteDeadline(time.Now().UTC().Add(60 * time.Second))

		client.conn = tcpConn
		return nil
	}
}

func (client *gantryTCPClient) Disconnect() error {
	return client.conn.Close()
}

func (client *gantryTCPClient) Write(envelope *proto.Envelope) error {

	data, err := protobuf.Marshal(envelope)
	if err != nil {
		return err
	}

	dataSize := len(data)

	data = append([]byte{byte(dataSize)}, data...)

	_, err = client.conn.Write(data)
	return err

}

// =====================================================================
// UDP
func (client *gantryUDPClient) Connect() (*gantryUDPConn, error) {

	addr, err := net.ResolveUDPAddr("udp4", client.RemoteAddr+":"+client.RemotePort)

	if err != nil {
		return nil, err
	}

	if udpConn, err := net.DialUDP("udp", nil, addr); err != nil {
		return nil, err
	} else {
		udpConn.SetWriteBuffer(512)
		return &gantryUDPConn{udpConn}, nil
	}

}

func (client *gantryUDPConn) Write(envelope *proto.Envelope) error {

	data, err := protobuf.Marshal(envelope)
	if err != nil {
		return err
	}

	if len(data) > 512 {
		return errors.New("UDP packet too big. Max allowed: 512 bytes")
	}

	_, err = client.conn.Write(data)
	return err
}

func (client *gantryUDPConn) Close() error {
	return client.conn.Close()
}
