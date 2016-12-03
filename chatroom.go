package main

import (
	"net"
)

type chatroom struct{
	Users []net.Conn
	NewUser chan user
	RemoveUsers chan user
	Messages chan string
}

func newRoom(){
	new_room := &chatroom{}
	new_room.NewUser := make(chan user)
	new_room.RemoveUsers := make(chan user)
	new_room.Messages := make(chan string)
	return new_room
}

func leaveRoom(leaver user, room chatroom){
	room.RemoveUsers <- leaver
}

func joinRoom(joinee user, room chatroom){
	room.NewUsers <- joinee
}

func chatRoom(initial_user *user, room_channel chan user) {
	users := make([]user, 0, 0)
	users = append(users, *initial_user)
	sendMessages(initial_user.username+" has joined", users, *initial_user, strings.Split(roomName, "room")[1])
	for {
		select {
		case newUser := <-room.NewUser:
			room.Users = append(room.Users, newUser)
			sendMessages(newUser.username+" has joined", room, newUser)

		case remUser := <-room.RemoveUsers:
			mesg := "LEFT_CHATROOM:" + strings.Split(room.roomName, "room")[1] + "\nJOIN_ID:" + room.Users[i].join_id + "\n"
			sendToUSers(mesg, room)
			leftMesg := "CHAT:" + strings.Split(roomName, "room")[1] + "\nCLIENT_NAME:" + remUser.Username + "\nMESSAGE: " + remUser.Username + " has left the chatroom"
			sendToUsers(leftMesg, room)

			i := 0
			for i=0; i<len(room.Users); i++ {
				if room.Users[i] == remUser {
					break
				} 
			}
			room.Users = append(room.Users[:i], room.Users[i+1:]...)

		case message := <-room.Messages:
		        sendMessages(message, room)
		default:
	        }
	}
}

/*
   sendMessages takes a message, a sender and a list of the users connections whom are in the chatroom
*/
func sendMessages(message string, room chatroom, sender user) {
    mesg := message + "Hi"
    sendToUsers(mesg, room)
}

func sendToUsers(message string, room chatroom){
	users := room.Users
	for i := 0; i < len(users); i++ {
		users[i].Writer.Write(message)
		users[i].Writer.Flush()
	}
}
