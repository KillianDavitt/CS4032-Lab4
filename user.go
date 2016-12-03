package main

import (
	"bufio"
	"strings"
)

type User struct {
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Username string
	JoinId   string
}

func newUser(reader *bufio.Reader, writer *bufio.Writer, l4 string) *User {

	username := strings.Split(l4, "CLIENT_NAME:")[1]

	id := "3"
	new_user := &User{reader, writer, username, id}
	return new_user
}
