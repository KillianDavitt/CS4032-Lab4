package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	port := os.Args[1]
	port = ":" + port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Fatal Error")
	}
	rooms := make(map[string]*chatroom)

	terminate_chan := make(chan bool)
	for {
		conn, err := listener.Accept()
		log.Print("New conn")
		defer conn.Close()
		if err != nil {
			fmt.Println("Fatal Error")
		}
		go handleConnection(conn, &listener, terminate_chan, rooms)
	}
}

func handleConnection(conn net.Conn, listener *net.Listener, terminate_chan chan bool, rooms map[string]*chatroom) {
	log.Print("Accepted new conn")
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var new_user *User
	madeUser := false
	for {
		l1, _ := reader.ReadString(byte('\n'))
		if l1 == "KILL_SERVICE" {
			os.Exit(0)
		}
		log.Print(l1)
		if strings.HasPrefix(l1, "HELO") {
			reply := l1 + "IP:10.62.0.83\nPort:8000\nStudentID:13319024"
			conn.Write([]byte(reply))
			return
		}
		_, _ = reader.ReadString(byte('\n'))
		_, _ = reader.ReadString(byte('\n'))

		if strings.HasPrefix(l1, "LEAVE_CHATROOM") {
			roomId := strings.TrimSpace(strings.Split(l1, "LEAVE_CHATROOM:")[1])
			room := getRoomById(roomId, rooms)
			leaveRoom(new_user, room)
			continue
		}

		if strings.HasPrefix(l1, "DISCONNECT") {
			disconnectUser(new_user, rooms)
			conn.Close()
			return
		}
		l4, _ := reader.ReadString(byte('\n'))
		if strings.HasPrefix(l1, "CHAT:") {
			roomId := strings.TrimSpace(strings.Split(l1, "CHAT:")[1])
			message := strings.Replace(strings.Split(l4, "MESSAGE:")[1], "\n", "", 3)
			room := getRoomById(roomId, rooms)
			_, _ = reader.ReadString(byte('\n'))

			messageRoom(struct {
				User
				string
			}{*new_user, message}, room)
		}

		if !madeUser {
			log.Print("made new user")
			madeUser = true
			new_user = newUser(reader, writer, l4, len(rooms))
		}

		
		if strings.HasPrefix(l1, "JOIN_CHATROOM:") {
			roomName := strings.Replace(strings.Split(l1, "JOIN_CHATROOM:")[1], "\n", "", 2)
			room := rooms[roomName]
			if room == nil {
				// Room doesn't exist, make it
				rooms[roomName] = newRoom(roomName, len(rooms)+1)
				go chatRoom(new_user, rooms[roomName])
			} else {
				// Room already exists, send the conn in  the channel
				joinRoom(new_user, room)
			}
		}
	}
}
