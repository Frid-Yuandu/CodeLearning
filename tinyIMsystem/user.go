package main

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
)

type User struct {
	name           string
	addr           string
	ReceiveMessage chan string
	conn           net.Conn
	server         *Server
}

func (u *User) Online() {
	u.server.onlineMap.Store(u.name, u)
	u.server.Broadcast(u, "Has logged in.")
}

func (u *User) Offline() {
	u.server.onlineMap.Delete(u.name)
	u.server.Broadcast(u, "Has logged out.")
}

func (u *User) dealTimeout() {
	u.SendToUser("Your session has timed out.")
	u.releaseSource()
	runtime.Goexit()
}

func (u *User) releaseSource() {
	u.server.onlineMap.Delete(u.name)
	close(u.ReceiveMessage)
	_ = u.conn.Close()
}

// Rename splits a provided message to get a new name, then call the rename
// method to update the username and to handle the error.
func (u *User) Rename(msg string) {
	newName := strings.Split(msg, "|")[1]
	u.rename(newName)
}

func (u *User) rename(newName string) {
	err := u.updateValidName(newName)
	if err != nil {
		u.SendToUser(err.Error())
		return
	}
	u.SendToUser("rename successful:" + newName + "\n")
}

// updateValidName, firstly, detects whether the new name exists or not firstly.
// It updates the username and returns `nil` while throwing an error if the new
// name already exists.
func (u *User) updateValidName(newName string) error {
	if u.server.UserExists(newName) {
		return fmt.Errorf("name %s is already in use", newName)
	}
	u.server.UpdateMapUsername(u.name, newName)
	u.name = newName
	return nil
}

// SendToUser provide an API to send message to current user's client and will
// not send to other users.
func (u *User) SendToUser(msg string) {
	_, err := u.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("user:", u.name, "SendToUser error:", err)
	}
}

func (u *User) selectMessageProcess(msg string) {
	if msg == "who" {
		u.SearchOnlineUsers()
	} else if strings.HasPrefix(msg, "updateValidName|") {
		// Expected message type: updateValidName|newName.
		u.Rename(msg)
	} else if strings.HasPrefix(msg, "to|") {
		// Expected message type: to|remoteName|msg
		u.PrivateChat(msg)
	} else {
		u.server.Broadcast(u, msg)
	}
}

// ListenMessage listens to user channel, and sends message to client if receive message.
func (u *User) ListenMessage() {
	for {
		msg, ok := <-u.ReceiveMessage
		if !ok {
			fmt.Println("user:", u.name, "channel close")
			return
		}
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("raise from User.ListenMessage | user:", u.name, "conn.Write error:", err)
		}
	}
}

func (u *User) SearchOnlineUsers() {
	s := u.server
	s.onlineMap.Range(func(_, otherUser any) bool {
		onlineMsg := "[" + otherUser.(*User).addr + "]" + otherUser.(*User).name + ":" + "is online.\n"
		u.SendToUser(onlineMsg)
		return true
	})
}

func (u *User) PrivateChat(msg string) {
	u.decodeToPrivateChat(msg)
}

func (u *User) decodeToPrivateChat(msg string) {
	splitMsg := strings.Split(msg, "|")
	err := u.checkPrivateChatType(splitMsg)
	if err != nil {
		u.SendToUser(err.Error())
		return
	}
	u.privateChat(splitMsg)
}

func (u *User) checkPrivateChatType(splitMsg []string) error {
	if len(splitMsg) != 3 {
		return errors.New("please type in 'to|remoteName|msg'")
	}
	remoteName, content := splitMsg[1], splitMsg[2]
	if !u.server.UserExists(remoteName) {
		return fmt.Errorf("no such user:" + remoteName)
	} else if content == "" {
		return errors.New("empty message, please type something")
	}
	return nil
}

func (u *User) privateChat(splitMsg []string) {
	remoteName, content := splitMsg[1], splitMsg[2]
	remoteUser, _ := u.server.onlineMap.Load(remoteName)
	remoteUser.(*User).SendToUser(u.name + " send to you:" + content)
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		name:           userAddr,
		addr:           userAddr,
		ReceiveMessage: make(chan string),
		conn:           conn,
		server:         server,
	}

	go user.ListenMessage()
	return user
}
