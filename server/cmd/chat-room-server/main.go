package main

import (
	"chat-room-cli/client/api/handler"
	"log"
)

func main() {
	log.Fatal(handler.Run())
}
