package networking

import (
	"bufio"
	//"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"time"
)

var counter uint64 = 0

type rpcjJsonClient struct {
	Address string
	Port    string
	conn    *net.TCPConn
}

// clientRequest represents a JSON-RPC request sent by a client.
type clientRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Method: method,
		Params: args, //interface{}{args},
		Id:     uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func decodeClientResponse(r io.Reader, reply interface{}) error {
	var c clientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	if c.Error != nil {
		return &Error{Data: c.Error}
	}
	if c.Result != nil {
		return json.Unmarshal(*c.Result, reply)
	}
	return nil
}

//===========================================

func NewRPCJsonClient(address, port string) *rpcjJsonClient {
	return &rpcjJsonClient{
		Address: address,
		Port:    port,
	}
}

func (c *rpcjJsonClient) Connect() error {

	addr, err := net.ResolveTCPAddr("tcp4", c.Address+":"+c.Port)

	if err != nil {
		return err
	}

	if tcpConn, err := net.DialTCP("tcp4", nil, addr); err != nil { //.DialUDP("udp", nil, addr); err != nil {
		return err
	} else {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(60 * time.Second)
		tcpConn.SetLinger(5)

		c.conn = tcpConn
		return nil
	}

}

func (c *rpcjJsonClient) Call(method string, args interface{}, response interface{}) error {

	data, err := encodeClientRequest(method, args)
	if err != nil {
		return err
	}

	// write data
	c.conn.Write(data)

	// read the response
	reader := bufio.NewReader(c.conn)

	// var buf bytes.Buffer
	// _, err = io.ReadFull(r, buf)
	// if err != nil && err != io.EOF {
	// 	//log.Errorln(err)
	// 	//continue
	// 	return err
	// }

	return decodeClientResponse(reader, response)

}

func (c *rpcjJsonClient) Close() error {
	return c.conn.Close()
}
