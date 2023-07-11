package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
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
	// Message Channel is used to broadcast information to users.
	onlineMap sync.Map

	Message chan string
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

	// Asynchronous listens to users' messages.
	go s.ListenMessage()

	// Starting accept request.
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

	// This channel checks whether user is active.
	isActive := make(chan struct{})

	// This Goroutine will read the current user's messages and call the
	// ProcessMessage method after a simple conversation.
	go func() {
		buf := make([]byte, HandlerReadBuffer)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("raise from Server.Handle | conn.Read error:", err)
				return
			}
			if n == 0 {
				u.Offline()
				return
			}

			msg := string(buf[:n-1])
			u.ProcessMessage(msg)
			// Expressing current user is active while receive any message.
			isActive <- struct{}{}
		}
	}()

	// Block this goroutine
	for {
		select {
		case <-isActive:
			// This statement means that current user is active. Timeout will be reset
			// automatically while this statement written above time.After case.
		case <-time.After(5 * time.Minute):
			// Has timed out
			// force user to logout
			u.SendToUser("Your session has timed out.")

			// Release user source
			s.onlineMap.Delete(u.Name)
			close(u.C)
			_ = conn.Close()

			// Exit goroutine
			runtime.Goexit()
		}
	}
}

// UpdateUserMap updates the user map with a new name for an existing user.
func (s *Server) UpdateUserMap(name, newName string) {
	u, _ := s.onlineMap.Load(name)
	s.onlineMap.Store(newName, u.(*User))
	s.onlineMap.Delete(name)
}

// ListenMessage listens to the Message channel of the Server. It will traverse
// the onlineMap using the Range method to send messages to every user. While
// nothing is in the Message channel, it will be blocked.
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		s.onlineMap.Range(func(_, u any) bool {
			u.(*User).C <- msg
			return true
		})
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		onlineMap: sync.Map{},
		Message:   make(chan string),
	}
}
