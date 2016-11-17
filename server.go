package main;

import (
    "fmt"
    "net"
    "os"
)

func main(){
    port := os.Args[1]
    port = ":" + port
	listener, err := net.Listen("tcp", port)
	if err != nil {
	    fmt.Println("Fatal Error")
    }
    terminate_chan := make(chan bool)
	for {
			conn, err := listener.Accept()
			defer conn.Close()
            if err != nil {
			    fmt.Println("Fatal Error")
            }
			go handleConnection(conn, &listener, terminate_chan)
	}

}

func handleConnection(conn net.Conn, listener *net.Listener, terminate_chan chan bool){
    for {
        buf := make([]byte, 1024)
        n, _ := conn.Read(buf)
        interpretMessage(buf[:n], conn, listener, terminate_chan)
    }
}

func interpretMessage(byte_message []byte, conn net.Conn, listener *net.Listener, terminate_chan chan bool){
	// Dict of rooms with channels, send new connections via the socket to the thread.

    ip, _ := os.Hostname()
    ip = "10.62.0.83"
    str := string(byte_message)
    fmt.Print(str)
    if str == "HELO BASE_TEST\n" {
        conn.Write([]byte("HELO BASE_TEST\nIP:" + ip + "\nPort:8080\nStudentID:13319024\n"))
    } else if str == "KILL_SERVICE\n" {
        conn.Close()
        //list.Close()
        os.Exit(0)
    } else {
        //
    }
    return
}

func chatRoom(initial_user net.Conn, room_channel chan net.Conn){
	users := make([]net.Conn, 10)
	users[0] = initial_user
	for {
		newUser := <- room_channel
		if newUser != nil {

		}

		for i:=0; i<len(users); {
			buf := make([]byte, 1024)
			n, err := users[i].Read(buf)
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
	return
}


