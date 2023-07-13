package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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
		"3.Search online users\n" +
		"4.Rename\n" +
		"0.Exit")

	_, err := fmt.Scanln(&selectFlag)
	if err != nil {
		fmt.Println("fmt.Scanln error:", err)
	}

	if selectFlag >= 0 && selectFlag <= 4 {
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
			c.PublicChat()
		case 2:
			fmt.Println("=====>Private chat<=====")
			c.PrivateChat()
		case 3:
			fmt.Println("=====>Search online users<=====")
			c.SearchOnlineUsers()
		case 4:
			fmt.Println("=====>Rename<=====")
			c.Rename()
		}
	}
}

func (c *Client) DealResponse() {
	_, err := io.Copy(os.Stdout, c.conn)
	if err != nil {
		fmt.Println("io.Copy error:", err)
	}
}

func (c *Client) Rename() bool {
	fmt.Println(">>>>>please type in new name:")
	_, _ = fmt.Scanln(&c.Name)

	sendMsg := "updateValidName|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("raise from Client.Rename | conn.Write error:", err)
		return false
	}
	return true
}

func (c *Client) SearchOnlineUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("raise from Client.SearchOnlineUsers | conn.Write error:", err)
		return
	}
}

func (c *Client) PublicChat() {
	var chatMsg string
	for strings.ToLower(chatMsg) != "exit" {
		chatMsg = ""
		fmt.Println(">>>>>please type in message, exit to exit:")
		fmt.Scanln(&chatMsg)

		if len(chatMsg) == 0 {
			fmt.Println("message can't be empty")
			continue
		}

		sendMsg := chatMsg + "\n"
		_, err := c.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("raise from Client.PublicChat | conn.Write error:", err)
			return
		}
	}
}

func (c *Client) PrivateChat() {
	var chatWithUser string
	var chatMsg string

	fmt.Println("online users:")
	c.SearchOnlineUsers()

	fmt.Println(">>>>>please select user to chat with, exit to exit:")
	fmt.Scanln(&chatWithUser)

	if strings.ToLower(chatWithUser) == "exit" {
		return
	} else if chatWithUser == c.Name {
		fmt.Println("you can't chat with yourself")
		return
	}

	fmt.Println(">>>>>please type in message, exit to exit:")
	fmt.Scanln(&chatMsg)
	if strings.ToLower(chatMsg) != "exit" {
		return
	}

	for strings.ToLower(chatMsg) != "exit" {
		if len(chatMsg) == 0 {
			fmt.Println("message can't be empty")
			continue
		}

		sendMsg := "to|" + chatWithUser + "|" + chatMsg + "\n"
		_, err := c.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("raise from Client.PrivateChat | conn.Write error:", err)
			return
		}

		chatMsg = ""
		fmt.Println(">>>>>please type in message, exit to exit:")
		fmt.Scanln(&chatMsg)
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

	go client.DealResponse()

	client.Run()
}
