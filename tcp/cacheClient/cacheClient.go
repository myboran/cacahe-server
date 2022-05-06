package cacheClient

import (
	"fmt"
	"net"
	"strconv"
)

type client struct {
	network string
	address string
}

type Cmd struct {
	Op    string
	Key   string
	Value string
	Error error
}

func New(network string, address string) client {
	return client{network: network, address: address}
}

func (c *client) Run(cmd *Cmd) {
	conn, err := net.Dial(c.network, c.address)
	if err != nil {
		fmt.Println("net Dial error: ", err)
		return
	}
	defer conn.Close()
	op := cmd.Op
	bytes := make([]byte, 1024)
	switch op {
	case "get":
		klen := len(cmd.Key)
		bytes = []byte("G" + strconv.Itoa(klen) + " " + cmd.Key)
	case "set":
		klen := len(cmd.Key)
		vlen := len(cmd.Value)
		bytes = []byte("S" + strconv.Itoa(klen) + " " + strconv.Itoa(vlen) + " " + cmd.Key + cmd.Value)
	case "del":
		klen := len(cmd.Key)
		bytes = []byte("D" + strconv.Itoa(klen) + " " + cmd.Key)
	default:
		fmt.Printf("command is %v, could be get/set/del: ", op)
		return
	}
	fmt.Println("send: ", bytes)
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println("conn write error: ", err)
		return
	}
	bytes2 := make([]byte, 1024)
	_, err = conn.Read(bytes2)
	if err != nil {
		fmt.Println("conn read error: ", err)
		return
	}
	fmt.Println("success: res->", string(bytes2))
}
