package ovsdb

import (
	"bufio"
	//"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	//log "github.com/golang/glog"
	//"github.com/cenkalti/rpc2"
	//"io"
	//"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var counter uint64 = 0

type message struct {
	Seq    uint64
	Method string
	Error  string
}

type rpcjJsonClient struct {
	mutex        sync.Mutex // protects pending, seq, request
	pending      map[string]chan json.RawMessage
	Address      string
	Port         string
	conn         *net.TCPConn
	handlers     []NotificationHandler
	disconnected bool
}

type tableUpdate struct {
	Table map[string]map[string]map[string]interface{}
}

// clientRequest represents a JSON-RPC request sent by a client.
type clientRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id string `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error",omitempty`
	Id     int64           `json:"id"`
}

type clientRR struct {
	Id string `json:"id"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`
	// response result
	Result json.RawMessage `json:"result"`
	// response error
	Error interface{} `json:"error"`
}

func (client *rpcjJsonClient) AddNotificationHandler(notifier NotificationHandler) {
	if client.handlers == nil {
		client.handlers = []NotificationHandler{}
	}
	client.handlers = append(client.handlers, notifier)

}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func (client *rpcjJsonClient) encodeClientRequest(method string, args interface{}, responseChannel chan json.RawMessage) ([]byte, error) {
	id := strconv.Itoa(int(rand.Int31()))
	request := &clientRequest{
		Method: method,
		Params: args, //interface{}{args},
		Id:     id,
	}

	if responseChannel != nil {
		//fmt.Println("Adding channel to request ID:", id)
		//client.mutex.Lock()
		client.pending[id] = responseChannel
		//client.mutex.Unlock()
	}

	return json.Marshal(request)
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func (client *rpcjJsonClient) decodeClientResponse(response clientRR) {

	//fmt.Println("Response ID:", response.Id)

	//client.mutex.Lock()
	responseChannel, ok := client.pending[response.Id]
	//client.mutex.Unlock()

	// // means there was a request associated to the response
	if !ok {
		//fmt.Println("NO REQUEST ID ASSOCIATED WITH RESPONSE !!")
		//fmt.Println("Not interest in response")
		return
	}

	// Close the channel once we are done here
	defer close(responseChannel)

	//fmt.Printf("%s\n", response.Result)

	if response.Result != nil {
		responseChannel <- response.Result
	}

	delete(client.pending, response.Id)
}

//===========================================

//=============================================================

func NewRPCJsonClient(address, port string) rpcjJsonClient {
	return rpcjJsonClient{
		Address: address,
		Port:    port,
		pending: make(map[string]chan json.RawMessage),
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

		go c.readLoop()

		return nil
	}

}

func (c *rpcjJsonClient) Call(method string, args interface{}, interested bool) (chan json.RawMessage, error) {

	var responseChannel chan json.RawMessage
	if interested {
		responseChannel = make(chan json.RawMessage)
	}

	data, err := c.encodeClientRequest(method, args, responseChannel)

	if err != nil {
		return responseChannel, err
	}

	// fmt.Println("=======================================================")
	// fmt.Printf("REQUEST: %s\n\n", data)

	// write data
	_, err = c.conn.Write(data)
	return responseChannel, err

}

func (c *rpcjJsonClient) Close() error {
	c.disconnected = true
	return c.conn.Close()
}

func (client *rpcjJsonClient) readLoop() {

	//max := 0
	// create a copy
	reader := bufio.NewReader(client.conn)
	//var buffer bytes.Buffer

	for {

		//max++
		//fmt.Println(max)

		if client.disconnected {
			fmt.Println("Client closed the connection")
			return //errors.New("Client is disconnected")
		}

		var rr clientRR

		//err = json.Unmarshal(tcpData, &rr)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&rr); err != nil {
			fmt.Println("Error decoding RequestResponse !", err)
			continue
		}

		//fmt.Printf("RESULT: %s\n\n", rr.Result)

		// means it's a server request !
		if rr.Method != "" {
			switch {
			case rr.Method == "echo":
				for _, handler := range client.handlers {
					if _, ok := rr.Params.([]interface{}); ok {
						var resp []interface{}
						handler.Echo(resp)
						client.Call("echo", []string{"ping"}, false)
					} else {
						fmt.Println("ECHO conversion error")
					}
				}
				continue
			case rr.Method == "update":
				//fmt.Println("UPDATE", data)
				if params, ok := rr.Params.([]interface{}); ok {

					if len(params) < 2 {
						continue //errors.New("Invalid Update message")
						//continue
					}
					// Ignore params[0] as we dont use the <json-value> currently for comparison

					raw, ok := params[1].(map[string]interface{})
					if !ok {
						continue //errors.New("Invalid Update message - 2")
					}
					var rowUpdates map[string]map[string]RowUpdate

					b, err := json.Marshal(raw)
					if err != nil {
						continue //err
					}
					err = json.Unmarshal(b, &rowUpdates)
					if err != nil {
						continue //err
					}

					// Update the local DB cache with the tableUpdates
					tableUpdates := getTableUpdatesFromRawUnmarshal(rowUpdates)
					for _, handler := range client.handlers {
						handler.Update(params, tableUpdates)
					}

				} else {
					fmt.Println("Could not cast to TableUpdates")
				}

				continue
			case rr.Method == "steal":
				continue
			case rr.Method == "lock":
				for _, handler := range client.handlers {
					if params, ok := rr.Params.([]interface{}); ok {
						handler.Locked(params)
					}
				}
				continue
			case rr.Method == "unlock":
				for _, handler := range client.handlers {
					if params, ok := rr.Params.([]interface{}); ok {
						handler.Locked(params)
					}
				}
				continue
			}
		} else {
			client.decodeClientResponse(rr)
		}

	}

}
