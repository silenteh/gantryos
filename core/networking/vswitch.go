package networking

// http://json-rpc.org/wiki/specification

import (
	"fmt"
	//"github.com/socketplane/libovsdb"
	//"net/rpc"
	"net/rpc/jsonrpc"
)

type request struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     uint64   `json:"id"`
}
type notification struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type rpcresponse struct {
	Result []string  `json:"result"`
	Error  *ovsError `json:"error"`
	Id     string    `json:"id"`
}

type ovsError struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

type Args struct {
	A, B string
}

type vSwitch struct {
	//client *libovsdb.OvsdbClient
}

func connectController() {

	c, err := jsonrpc.Dial("tcp4", "192.168.1.111:6633")

	defer c.Close()
	// By default libovsdb connects to 127.0.0.0:6400.

	//client, err := libovsdb.Connect("192.168.1.111", 6633)
	if err != nil {
		fmt.Println(err)
		return
	}
	//c := jsonrpc.NewClientCodec(conn)
	//c := jsonrpc.NewClientCodec(conn)
	//c.Call("echo", , reply)
	var resp []Args //rpcresponse
	err = c.Call("echo", "", &resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", resp)

	//client.ListDbs()

}
