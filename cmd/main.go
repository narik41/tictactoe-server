package main

import (
	"log"

	"github.com/narik41/tictactoe-helper/core"
	"github.com/narik41/tictactoe-server/internal"
	"github.com/narik41/tictactoe-server/internal/repo"
)

func main() {
	log.Println("!!! Starting the tic tac toe server !!!")

	// session
	sessionManager := internal.NewSessionManager()
	gameSessionManager := internal.NewGameSessionManager()

	responseSender := internal.NewResponseSender(sessionManager)

	queue := internal.NewSessionQueue(gameSessionManager, responseSender)
	queue.Start()

	// repo
	userRepo := repo.NewUserRepo()

	// register msg handler
	router := internal.NewMessageRouter()
	router.RegisterHandler(core.MSG_LOGIN_PAYLOAD, internal.NewLoginHandler(userRepo, gameSessionManager, queue))
	router.RegisterHandler(core.PLAYER_MOVE, internal.NewPlayerMoveHandler(gameSessionManager))
	router.RegisterHandler(core.HEARTBEAT, internal.NewHeartbeatHandler())

	server := internal.NewServer(sessionManager, gameSessionManager, router)
	err := server.Start("localhost:9000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
