package main

type user struct {
	reader   *bufio.Reader
	writer   *bufio.Writer
	username string
	join_id  string
}

func newUser(reader *bufio.Reader, writer *bufio.Writer, l4 string) {

	username := strings.Split(l4, "CLIENT_NAME:")[1]

	new_user := &user{reader, writer, username, id}
	return new_user
}
