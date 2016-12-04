package main

import (
	"bufio"
	"log"
	"strings"
)

type User struct {
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Username string
	JoinId   string
}

func newUser(reader *bufio.Reader, writer *bufio.Writer, l4 string, n int) *User {
	log.Print("l4 is : " + l4)
	username := strings.Split(l4, "CLIENT_NAME:")[1]
	username = strings.Replace(username, "\n", "", 2)
	username = strings.TrimSpace(username)
	id := "3"
	new_user := &User{reader, writer, username, id}
	new_user.JoinId = "2"
	return new_user
}

func disconnectUser(user *User, rooms map[string]*chatroom){
	for _,v := range rooms {
	// For every chatroom
		for j:=0; j<len(v.Users); j++{
		// Check if any of the users match
			if v.Users[j].Username == user.Username {
				log.Print("found a case of a chatroom to discon")
				postDisconnect(v, user)
			}
		}
	}
}
