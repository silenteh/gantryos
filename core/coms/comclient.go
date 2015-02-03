package coms

import (
	//protobuf "github.com/gogo/protobuf/proto"
	//log "github.com/golang/glog"
	//"github.com/silenteh/gantryos/core/proto"

	//"bufio"
	"errors"
	//"io/ioutil"
	"net"
)

type gantryTCPClient struct {
	RemoteAddr string
	RemotePort string
}

type gantryUDPClient struct {
	RemoteAddr string
	RemotePort string
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

func (client *gantryTCPClient) Connect() (*gantryTCPConn, error) {

	addr, err := net.ResolveTCPAddr("tcp4", client.RemoteAddr+":"+client.RemotePort)

	if err != nil {
		return nil, err
	}

	if tcpConn, err := net.DialTCP("tcp4", nil, addr); err != nil { //.DialUDP("udp", nil, addr); err != nil {
		return nil, err
	} else {
		return &gantryTCPConn{tcpConn}, nil
	}
}

func (client *gantryTCPConn) WriteMessage(data []byte) error {

	//client.conn.SetWriteBuffer(512)
	client.conn.SetWriteBuffer(1024)
	client.conn.SetNoDelay(true)

	if _, err := client.conn.Write(data); err != nil {
		return err
	}

	return nil

}

func (client *gantryTCPConn) Close() error {
	return client.conn.Close()
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
		return &gantryUDPConn{udpConn}, nil
	}

}

func (client *gantryUDPConn) WriteMessage(data []byte) error {
	if len(data) > 512 {
		return errors.New("UDP packet too big. Max allowed: 512 bytes")
	}

	client.conn.SetWriteBuffer(512)

	if _, err := client.conn.Write(data); err != nil {
		return err
	}

	return nil

}

func (client *gantryUDPConn) Close() error {
	return client.conn.Close()
}
