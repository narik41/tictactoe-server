package main

import (
	"log"

	"github.com/narik41/tictactoe-server/internal"
)

func main() {
	log.Println("!!! Starting the tic tac toe server !!!")

	v1MsgHandler := internal.NewVersion1MsgHandler()
	msgHandler := internal.NewMessageHandler(v1MsgHandler)

	server := internal.NewServer(msgHandler)
	err := server.Start("localhost:9000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
