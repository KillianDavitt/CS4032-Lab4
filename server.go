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
    //max_threads := 5
    terminate_chan := make(chan bool)
    //num_threads := 0
	for {
			conn, err := listener.Accept()
			defer conn.Close()
            if err != nil {
			    fmt.Println("Fatal Error")
            }
			go handleConnection(conn, &listener, terminate_chan)
            //term := <-terminate_chan
            //if term {
            //    listener.Close()
            //    os.Exit(0)
            //}   
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
    ip, _ := os.Hostname()
    ip = "10.62.0.83"
    str := string(byte_message)
    fmt.Print(str)
    if str == "HELO BASE_TEST\n" {
        conn.Write([]byte("HELO BASE_TEST\nIP:" + ip + "\nPort:8080\nStudentID:13319024\n"))
    } else if str == "KILL_SERVICE\n" {
        conn.Close()
        //list := *listener
        //list.Close()
        os.Exit(0)
    } else {
        //
    }
    return
}
