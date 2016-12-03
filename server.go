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

		if strings.HasPrefix(l1, "HELO") {
			reply := l1 + "IP:10.62.0.83\nPort:8000\nStudentID:13319024"
			conn.Write([]byte(reply))
			return
		}
		l2, _ := reader.ReadString(byte('\n'))
		l3, _ := reader.ReadString(byte('\n'))
		l4, _ := reader.ReadString(byte('\n'))
		if !madeUser {
			log.Print("made new user")
			madeUser = true
			new_user = newUser(reader, writer, l4)
		}
		message := l1 + l2 + l3 + l4

		// Dict of rooms with channels, send new connections via the socket to the thread.
		log.Print("Interpreting message " + message)
		lines := strings.Split(message, "\n")
		log.Print(len(lines))
		if len(lines) < 3 {
			return
		}
		if strings.HasPrefix(lines[0], "CHAT:") {
			roomName := strings.Split(lines[0], "CHAT:")[1]
			room := rooms[roomName]
			messageRoom("HI", room)
		} else if strings.HasPrefix(lines[0], "JOIN_CHATROOM:") {
			log.Print("Joining chatroom")
			roomName := lines[0][14:]
			log.Print("found name " + roomName)
			room := rooms[roomName]
			log.Print("found room")
			if room == nil {
				// Room doesn't exist, make it
				rooms[roomName] = newRoom(roomName, len(rooms)+1)
				go chatRoom(new_user, rooms[roomName])
			} else {
				// Room already exists, send the conn in  the channel
				joinRoom(new_user, room)
			}
		} else if strings.HasPrefix(lines[0], "LEAVE_CHATROOM") {
			log.Print("Leaving chatroom")
			roomName := strings.Split(lines[0], "LEAVE_CHATROOM:")[1]
			room := rooms[roomName]
			leaveRoom(new_user, room)
		}
	}
}
