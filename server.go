package main;

import (
	"fmt"
	"net"
	"os"
	"strings"
	"log"
)


func main(){
	port := os.Args[1]
	port = ":" + port
	listener, err := net.Listen("tcp", port)
	if err != nil {
	    fmt.Println("Fatal Error")
	}
	rooms := make(map[string]chan net.Conn)

	
	terminate_chan := make(chan bool)
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			    fmt.Println("Fatal Error")
		}
		go handleConnection(conn, &listener, terminate_chan, rooms)
	}

}

func handleConnection(conn net.Conn, listener *net.Listener, terminate_chan chan bool, rooms map[string]chan net.Conn){
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	interpretMessage(buf[:n], conn, terminate_chan, rooms)
	
}

func interpretMessage(byte_message []byte, conn net.Conn, terminate_chan chan bool, rooms map[string]chan net.Conn){
	// Dict of rooms with channels, send new connections via the socket to the thread.
	message := string(byte_message)
	log.Print("Interpreting message " + message)
	attribs := strings.Split(message, "\n")
	log.Print(attribs[0][0:13])
	if attribs[0][0:13] == "JOIN_CHATROOM" {
		log.Print("User is joining a room")
		roomName := attribs[0][13:]
		room := rooms[roomName]
		if room == nil {
			rooms[roomName] = make(chan net.Conn)
			log.Print("Creating chat room")
			go chatRoom(conn, rooms[roomName])
			conn.Write([]byte("chatroom has been created"))
		} else {
			// Room already exists, send the conn in  the channel
			room <- conn
			conn.Write([]byte("You have been connected"))
		}
	}
	
}

func chatRoom(initial_user net.Conn, room_channel chan net.Conn){
	users := make([]net.Conn, 0, 0)
	users = append(users, initial_user)
	for {
		newUser := <- room_channel
		if newUser != nil {
			users = append(users, newUser)
		}

		for i:=0; i<len(users); i++ {
			buf := make([]byte, 1024)
			n, err := users[i].Read(buf)
			log.Print("User sent message: " + string(buf))
			if err != nil {
				fmt.Println(err)
			}
			sendMessages(buf[:n], users)
		}
	}
}

/*
   sendMessages takes a message, a sender and a list of the users connections whom are in the chatroom
*/
func sendMessages(message []byte, users []net.Conn){
	for i:=0; i<len(users); i++ {
		users[i].Write(message)
	}
}


