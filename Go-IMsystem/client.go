package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	serverIP   string
	serverPort int
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	selectFlag int
}

func (c *Client) Menu() bool {
	var selectFlag int
	fmt.Println("1.Public chat\n" +
		"2.Private chat\n" +
		"3.Rename\n" +
		"0.Exit")

	_, err := fmt.Scanln(&selectFlag)
	if err != nil {
		fmt.Println("fmt.Scanln error:", err)
	}

	if selectFlag >= 0 && selectFlag <= 3 {
		c.selectFlag = selectFlag
		return true
	}
	fmt.Println("Please enter a valid number")
	return false
}

func (c *Client) Run() {
	for c.selectFlag != 0 {
		for !c.Menu() {
		}

		switch c.selectFlag {
		case 1:
			fmt.Println("=====>Public chat<=====")
		case 2:
			fmt.Println("=====>Private chat<=====")
		case 3:
			fmt.Println("=====>Rename<=====")
		}
	}
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		selectFlag: 999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set server ip(default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8080, "set server port(default: 8080)")
}

func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("=====>NewClient error<=====")
		return
	}
	fmt.Println("=====>client has been created<=====")

	client.Run()
}
