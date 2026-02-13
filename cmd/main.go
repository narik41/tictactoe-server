package main

import (
	"log"

	"github.com/narik41/tictactoe-server/internal"
)

func main() {
	log.Println("!!! Starting the tic tac toe server !!!")

	server := internal.NewServer()
	err := server.Start("localhost:9000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
