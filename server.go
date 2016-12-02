package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type user struct {
	reader   *bufio.Reader
	writer   *bufio.Writer
	username string
}

func main() {
	port := os.Args[1]
	port = ":" + port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Fatal Error")
	}
	rooms := make(map[string]chan user)

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

func handleConnection(conn net.Conn, listener *net.Listener, terminate_chan chan bool, rooms map[string]chan user) {
	log.Print("Accepted new conn")
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	l1, _ := reader.ReadString(byte('\n'))

	if l1 == "KILL_SERVICE" {
		os.Exit(0)
	}

	if strings.HasPrefix(l1, "HELO") {
		conn.Write([]byte(l1))
		return
	}
	l2, _ := reader.ReadString(byte('\n'))
	l3, _ := reader.ReadString(byte('\n'))
	l4, _ := reader.ReadString(byte('\n'))

	message := l1 + l2 + l3 + l4
	//n, _ := conn.Read(buf)
	interpretMessage(message, reader, writer, terminate_chan, rooms)
}

func interpretMessage(message string, reader *bufio.Reader, writer *bufio.Writer, terminate_chan chan bool, rooms map[string]chan user) {
	// Dict of rooms with channels, send new connections via the socket to the thread.
	log.Print("Interpreting message " + message)
	attribs := strings.Split(message, "\n")
	log.Print(len(attribs))
	if len(attribs) < 3 {
		return
	}
	log.Print(attribs[0][0:13])
	if attribs[0][0:13] == "JOIN_CHATROOM" {
		log.Print("User is joining a room")
		username := strings.Join(strings.Split(attribs[3], ":")[1:], " ")
		var new_user_obj user
		new_user := &new_user_obj
		new_user.username = username
		new_user.reader = reader
		new_user.writer = writer
		roomName := attribs[0][13:]
		room := rooms[roomName]
		if room == nil {
			rooms[roomName] = make(chan user)
			log.Print("Creating chat room")
			go chatRoom(new_user, rooms[roomName], roomName)
			writer.Write([]byte("chatroom has been created\n"))
		} else {
			// Room already exists, send the conn in  the channel
			room <- *new_user
			writer.Write([]byte("You have been connected\n"))
		}
	}

}

func chatRoom(initial_user *user, room_channel chan user, roomName string) {
	users := make([]user, 0, 0)
	users = append(users, *initial_user)
	for {
		select {
		case newUser := <-room_channel:
			users = append(users, newUser)
		default:
			for i := 0; i < len(users); i++ {
				mesg, _ := users[i].reader.ReadString('\n')
				log.Print("User sent message: " + mesg)
				sendMessages(mesg, users, users[i], roomName)
			}

		}
	}
}

/*
   sendMessages takes a message, a sender and a list of the users connections whom are in the chatroom
*/
func sendMessages(message string, users []user, sender user, roomName string) {
	for i := 0; i < len(users); i++ {
		log.Print("Sending...")
		if users[i] != sender {
			mesg := "CHAT: " + roomName + "\nCLIENT_NAME: " + sender.username + "\nMESSAGE:" + string(message) + "\n\n"
			users[i].writer.Write([]byte(mesg))
		}
	}
}
