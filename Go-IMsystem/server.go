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
	ReadBuffer = 4096
)

type Server struct {
	IP   string
	Port int

	// Each client is a user instance. User instance should contain a channel, a
	// connection(, and a client?). User's channel is blocked until receive message.
	// Message Channel is used to broadcast information to users.
	OnlineMap sync.Map

	Message chan string
}

func (s *Server) UserExist(name string) (*User, bool) {
	u, ok := s.OnlineMap.Load(name)
	return u.(*User), ok
}

// UpdateUserMap updates the user map with a new name for an existing user.
//
// Parameters:
// - name: the current name of the user.
// - newName: the new name to be assigned to the user.
func (s *Server) UpdateUserMap(name, newName string) {
	u, _ := s.OnlineMap.Load(name)
	s.OnlineMap.Store(newName, u.(*User))
	s.OnlineMap.Delete(name)
}

func (s *Server) ListenMessage() {

	for {
		// Get the next message from the server
		msg := <-s.Message

		// To prevent any changes while we are processing the message
		s.OnlineMap.Range(func(_, u any) bool {
			u.(*User).C <- msg
			return true
		})
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	u := NewUser(conn, s) // Conn is bind to u
	u.Online()

	// This channel checks whether user is active.
	isActive := make(chan struct{})

	// Goroutines which are in charge of receive client message and broadcast it.
	go func() {
		buf := make([]byte, ReadBuffer)
		for {
			n, err := conn.Read(buf)
			// While finish read, conn.Read will raise io.EOF error.
			if err != nil && err != io.EOF {
				fmt.Println("raise from Server.Handle | conn.Read error:", err)
				return
			}
			if n == 0 {
				u.Offline()
				return
			}

			msg := string(buf[:n-1])
			u.DoMessage(msg)
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
			s.OnlineMap.Delete(u.Name)
			close(u.C)
			_ = conn.Close()

			// Exit goroutine
			runtime.Goexit()
		}
	}
}

func (s *Server) Start() {

	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", s.IP, s.Port),
	)
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// Goroutine which belong to server
	// is keep listening Message and transmit to clients.
	go s.ListenMessage()

	// Starting accept request.
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("raise from Server.Start | listener.Accept error:", err)
			continue
		}

		go s.Handler(conn)
	}
}

func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: sync.Map{},
		Message:   make(chan string),
	}
}
