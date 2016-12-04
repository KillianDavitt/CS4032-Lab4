package main

import (
	"log"
	"strconv"
)

type chatroom struct {
	Users       []User
	NewUser     chan User
	RemoveUsers chan User
	Disconnect  chan User
	Messages    chan struct {
		User
		string
	}
	RoomName string
	RoomId   int
}

func newRoom(roomName string, id int) *chatroom {
	new_room := &chatroom{}
	new_room.NewUser = make(chan User)
	new_room.RemoveUsers = make(chan User)
	new_room.Disconnect = make(chan User)
	new_room.Messages = make(chan struct {
		User
		string
	})
	new_room.RoomName = roomName
	new_room.RoomId = id
	return new_room
}

func getRoomById(roomId string, rooms map[string]*chatroom) *chatroom {
	var room *chatroom
	for _, v := range rooms {
		if strconv.Itoa(v.RoomId) == roomId {
			room = v
		}
	}
	return room

}

func leaveRoom(leaver *User, room *chatroom) {
	room.RemoveUsers <- *leaver
}

func joinRoom(joinee *User, room *chatroom) {
	room.NewUser <- *joinee
}

func postDisconnect(room *chatroom, user *User){
	log.Print("Sending discon to channel")
	room.Disconnect <- *user
}

func messageRoom(message struct {
	User
	string
}, room *chatroom) {
	room.Messages <- message
}

func chatRoom(initial_user *User, room *chatroom) {
	users := make([]User, 0, 50)
	users = append(users, *initial_user)
	room.Users = users
	initial_user.Writer.Write([]byte("JOINED_CHATROOM: " + room.RoomName + "\nSERVER_IP: 10.82.0.63\nPORT: 8000\nROOM_REF: " + strconv.Itoa(room.RoomId) + "\nJOIN_ID: " + initial_user.JoinId + "\n"))
	initial_user.Writer.Flush()
	sendMessages(initial_user.Username+" has joined", room, initial_user)
	for {
		select {
		case newUser := <-room.NewUser:
			room.Users = append(room.Users, newUser)
			newUser.Writer.Write([]byte("JOINED_CHATROOM: " + room.RoomName + "\nSERVER_IP: 10.82.0.63\nPORT: 8000\nROOM_REF: " + strconv.Itoa(room.RoomId) + "\nJOIN_ID: " + initial_user.JoinId + "\n"))
			newUser.Writer.Flush()

			sendMessages(newUser.Username+" has joined", room, &newUser)
			
		case disconUser := <-room.Disconnect:
			log.Print("discon")
			sendMessages(disconUser.Username + " has disconnected", room, &disconUser)
			
		case remUser := <-room.RemoveUsers:
			log.Print("leaving room in goroutine")
			mesg := "LEFT_CHATROOM:" + strconv.Itoa(room.RoomId) + "\nJOIN_ID:" + remUser.JoinId + "\n"
			remUser.Writer.Write([]byte(mesg))
			remUser.Writer.Flush()

			log.Print("Sent leave message back to sender")
			leftMesg := "CHAT: " + strconv.Itoa(room.RoomId) + "\nCLIENT_NAME: " + remUser.Username + "\nMESSAGE: " + remUser.Username + " has left the chatroom\n\n"

			sendToUsers(leftMesg, room)
			i := 0
			for i = 0; i < len(room.Users); i++ {
				if room.Users[i] == remUser {
					break
				}
			}
			room.Users = append(room.Users[:i], room.Users[i+1:]...)

		case message := <-room.Messages:
			sendMessages(message.string, room, &message.User)
		default:
		}
	}
}

/*
   sendMessages takes a message, a sender and a list of the users connections whom are in the chatroom
*/
func sendMessages(message string, room *chatroom, sender *User) {
	mesg := "CHAT: " + strconv.Itoa(room.RoomId) + "\nCLIENT_NAME: " + sender.Username + "\nMESSAGE:" + message + "\n\n"

	sendToUsers(mesg, room)
}

func sendToUsers(message string, room *chatroom) {
	users := room.Users
	log.Print("sending to users: " + string(len(users)))
	for i := 0; i < len(users); i++ {
		log.Print("Sending to a user there Ted")
		users[i].Writer.Write([]byte(message))
		users[i].Writer.Flush()
	}
}
