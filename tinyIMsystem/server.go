package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	HandlerReadBuffer = 4096
)

type Server struct {
	IP   string
	Port int

	// Each client is a user instance. User instance should contain a channel, a
	// connection(, and a client?). User's channel is blocked until receive message.
	// broadcastMessage Channel is used to broadcast information to users.
	onlineMap        sync.Map
	broadcastMessage chan string
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	go s.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}

		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	u := NewUser(conn, s)
	u.Online()
	isActive := make(chan struct{})

	go s.HandleMessage(u, isActive)

	for {
		select {
		case <-isActive:
			// This statement means that current user is active. Timeout will be reset
			// automatically while this statement written above time.After case.
		case <-time.After(5 * time.Second):
			u.dealTimeout()
		}
	}
}

// HandleMessage reads message from provided user's connection and call the
// selectMessageProcess method after a simple conversation.
func (s *Server) HandleMessage(u *User, isActive chan struct{}) {
	buf := make([]byte, HandlerReadBuffer)
	for {
		n, err := u.conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("raise from Server.Handle | conn.Read error:", err)
			return
		}
		if n == 0 {
			u.Offline()
			return
		}

		msg := string(buf[:n-1])
		u.selectMessageProcess(msg)
		isActive <- struct{}{}
	}
}

func (s *Server) UserExists(userName string) bool {
	_, ok := s.onlineMap.Load(userName)
	return ok
}

// UpdateMapUsername updates the user map with a new name for an existing user.
func (s *Server) UpdateMapUsername(name, newName string) {
	u, _ := s.onlineMap.Load(name)
	s.onlineMap.Store(newName, u.(*User))
	s.onlineMap.Delete(name)
}

// ListenMessage listens to the broadcastMessage channel of the Server. It will traverse
// the onlineMap using the Range method to send messages to every user. While
// nothing is in the broadcastMessage channel, it will be blocked.
func (s *Server) ListenMessage() {
	for {
		msg := <-s.broadcastMessage

		s.onlineMap.Range(func(_, u any) bool {
			u.(*User).ReceiveMessage <- msg
			return true
		})
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.addr + "]" + user.name + ":" + msg
	s.broadcastMessage <- sendMsg
}

func NewServer(ip string, port int) *Server {
	return &Server{
		IP:               ip,
		Port:             port,
		onlineMap:        sync.Map{},
		broadcastMessage: make(chan string),
	}
}
