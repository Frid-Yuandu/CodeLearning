package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// start listening message
	go user.ListenMessage()
	return user
}

func (u *User) Online() {
	s := u.server

	// To add user into OnlineMap.
	s.OnlineMap.Store(u.Name, u)

	// Broadcast current user connection message.
	s.Broadcast(u, "Has logged in.")
}

func (u *User) Offline() {
	s := u.server

	// To delete user in OnlineMap.
	s.OnlineMap.Delete(u.Name)

	// Broadcast current user logged out message.
	s.Broadcast(u, "Has logged out.")
}

func (u *User) Rename(newName string) {
	// check whether newName is already in OnlineMap
	s := u.server
	if _, ok := s.UserExist(newName); ok {
		// This statement means newName is already in OnlineMap.
		u.SendToUser("name is already in use")
		return
	}

	s.UpdateUserMap(u.Name, newName)
	u.Name = newName
	u.SendToUser("rename successful:" + newName + "\n")
}

// SendToUser provide an API to send message to current user's client and will
// not send to other users.
func (u *User) SendToUser(msg string) {
	_, err := u.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("user:", u.Name, "conn.Write error:", err)
	}
}

// DoMessage implement business of checking & sending message.
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// Searching online user.
		s := u.server

		s.OnlineMap.Range(func(_, otherUser any) bool {
			onlineMsg := "[" + otherUser.(*User).Addr + "]" + otherUser.(*User).Name + ":" + "is online.\n"
			u.SendToUser(onlineMsg)
			return true
		})

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// Expected message type: rename|newName. Everything after "|" will be the new
		// name. strings.Split method will split the message by "|" and returns a slice
		// of string containing the 2 elements: "rename" and the variable newName.
		newName := strings.Split(msg, "|")[1]
		u.Rename(newName)

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// Expected message type: to|remoteName|msg.
		splitMsg := strings.Split(msg, "|")
		if len(splitMsg) != 3 {
			u.SendToUser("side-text type error, please type in 'to|remoteName|msg'")
			return
		}

		remoteName := splitMsg[1]
		remoteUser, ok := u.server.UserExist(remoteName)
		if !ok {
			u.SendToUser("no such user:" + remoteName)
			return
		}

		content := splitMsg[2]
		if content == "" {
			u.SendToUser("empty message, please type something.")
			return
		}
		remoteUser.SendToUser(u.Name + " send to you:" + content)

	} else {
		u.server.Broadcast(u, msg)
	}
}

// ListenMessage listen to user channel, send to client if receive message.
func (u *User) ListenMessage() {
	for {
		msg, ok := <-u.C
		if !ok {
			fmt.Println("raise from User.ListenMessage | user:", u.Name, "channel close")
			return
		}
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("raise from User.ListenMessage | user:", u.Name, "conn.Write error:", err)
		}
	}
}
