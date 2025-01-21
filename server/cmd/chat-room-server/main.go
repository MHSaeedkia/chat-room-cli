package main

import (
	"chat-room-cli/server/api/handler"
	"log"
)

func main() {
	log.Fatal(handler.Run())
}
