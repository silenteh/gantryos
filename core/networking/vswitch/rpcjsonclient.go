package vswitch

import (
	"bufio"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
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
	closeLock    sync.Mutex
	writeLock    sync.Mutex
	readLock     sync.Mutex
	pending      map[string]chan json.RawMessage
	Address      string
	Port         string
	conn         *net.TCPConn
	handlers     []notificationHandler
	disconnected bool
}

// type tableUpdate struct {
// 	Table map[string]map[string]map[string]interface{}
// }

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

func (client *rpcjJsonClient) addnotificationHandler(notifier notificationHandler) {
	if client.handlers == nil {
		client.handlers = []notificationHandler{}
	}
	client.handlers = append(client.handlers, notifier)

}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func (client *rpcjJsonClient) encodeClientRequest(method string, args interface{}, respInterested bool) ([]byte, chan json.RawMessage, error) {
	id := strconv.Itoa(int(rand.Int31()))
	request := &clientRequest{
		Method: method,
		Params: args, //interface{}{args},
		Id:     id,
	}

	var channel chan json.RawMessage

	if respInterested {
		channel = make(chan json.RawMessage, 1)
		//fmt.Printf("## %s encodeClientRequest Mutex Locked\n", id)
		client.mutex.Lock()
		client.pending[id] = channel
		client.mutex.Unlock()
		//fmt.Printf("## %s encodeClientRequest Mutex UnLocked\n", id)
	}

	data, err := json.Marshal(request)

	return data, channel, err
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func (client *rpcjJsonClient) decodeClientResponse(response clientRR) {
	//fmt.Printf("RESULT: %s *** %s\n\n", response.Id, response.Result)
	client.mutex.Lock()
	//fmt.Printf("## %s decodeClientResponse Mutex Locked\n", response.Id)
	responseChannel, ok := client.pending[response.Id]
	// // means there was a request associated to the response
	if ok {
		responseChannel <- response.Result
		close(responseChannel)
		delete(client.pending, response.Id)
	}
	client.mutex.Unlock()
	//fmt.Printf("## %s decodeClientResponse Mutex Unlocked\n", response.Id)

}

//===========================================

//=============================================================

func newRPCJsonClient(address, port string) rpcjJsonClient {
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

	data, responseChannel, err := c.encodeClientRequest(method, args, interested)

	if err != nil {
		log.Errorln("Error sending request to vswitch", err)
		return responseChannel, err
	}

	// fmt.Println("=======================================================")
	//fmt.Printf("REQUEST: %s\n\n", data)

	// write data
	c.writeLock.Lock()
	//fmt.Println("## Call Write Locked")
	_, err = c.conn.Write(data)
	if err != nil {
		//fmt.Println("WRITE ERROR !:", err)
		log.Errorln("Error writing to socket", err)
	}
	c.writeLock.Unlock()
	//fmt.Println("## Call Write UnLocked")
	return responseChannel, err

}

func (c *rpcjJsonClient) Close() error {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	c.disconnected = true
	return c.conn.Close()
}

func (client *rpcjJsonClient) readLoop() {

	reader := bufio.NewReader(client.conn)
	decoder := json.NewDecoder(reader)

	for {

		// lock the access to the variable
		client.closeLock.Lock()
		// read the value
		disconnected := client.disconnected
		// unlock the access
		client.closeLock.Unlock()

		if disconnected {
			log.Infoln("Client disconnected from vswitch")
			return
		}

		var rr clientRR

		// this is a CRITICAL SECTION !!!
		client.closeLock.Lock()
		client.writeLock.Lock()
		if err := decoder.Decode(&rr); err != nil {
			client.closeLock.Unlock()
			client.writeLock.Lock()
			log.Errorln("Error decoding RequestResponse", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		client.writeLock.Unlock()
		client.closeLock.Unlock()

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
					}
				}
				time.Sleep(100 * time.Millisecond)
				continue
			case rr.Method == "update":

				// //fmt.Println("UPDATE", data)
				// if params, ok := rr.Params.([]interface{}); ok {

				// 	if len(params) < 2 {
				// 		continue //errors.New("Invalid Update message")
				// 		//continue
				// 	}
				// 	// Ignore params[0] as we dont use the <json-value> currently for comparison

				// 	raw, ok := params[1].(map[string]interface{})
				// 	if !ok {
				// 		continue //errors.New("Invalid Update message - 2")
				// 	}
				// 	var rowUpdates map[string]map[string]RowUpdate

				// 	b, err := json.Marshal(raw)
				// 	if err != nil {
				// 		continue //err
				// 	}
				// 	err = json.Unmarshal(b, &rowUpdates)
				// 	if err != nil {
				// 		continue //err
				// 	}

				// 	// Update the local DB cache with the tableUpdates
				// 	tableUpdates := getTableUpdatesFromRawUnmarshal(rowUpdates)
				// 	for _, handler := range client.handlers {
				// 		handler.Update(params, tableUpdates)
				// 	}

				// } else {
				// 	fmt.Println("Could not cast to TableUpdates")
				// }

				time.Sleep(100 * time.Millisecond)

				continue
			case rr.Method == "steal":
				time.Sleep(100 * time.Millisecond)
				continue
			case rr.Method == "lock":
				for _, handler := range client.handlers {
					if params, ok := rr.Params.([]interface{}); ok {
						handler.Locked(params)
					}
				}
				time.Sleep(100 * time.Millisecond)
				continue
			case rr.Method == "unlock":

				for _, handler := range client.handlers {
					if params, ok := rr.Params.([]interface{}); ok {
						handler.Locked(params)
					}
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
		} else {
			client.closeLock.Lock()
			client.decodeClientResponse(rr)
			client.closeLock.Unlock()
		}

		time.Sleep(100 * time.Millisecond)

	}

}
