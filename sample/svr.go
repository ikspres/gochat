package main

import "github.com/ikspres/gochat/server"

func main() {
	cr := server.NewChatRoom(":6666")
	cr.Go()
}
