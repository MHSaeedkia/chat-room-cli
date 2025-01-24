package cmd

import (
	"chat-room-cli/cmd/client"
	"chat-room-cli/cmd/server"
	"fmt"
	"strings"
)

const (
	SERVER = iota + 1
	CLIENT
)

type role int

func menu() role {
	var role string
	fmt.Printf("Welcome to cli chat room : \n**1-%-12s**\n**2-%-12s**\nChoise your role : ", "Server", "Client")
	fmt.Scan(&role)

	role = strings.ToLower(role)
	switch role {
	case "1":
		return SERVER
	case "2":
		return CLIENT
	case "server":
		return SERVER
	case "client":
		return CLIENT
	default:
		return CLIENT
	}
}

func Run() error {

	// starting menu
	role := menu()
	switch role {
	case SERVER:
		return server.Run()
	case CLIENT:
		return client.Run()
	}
	return nil
}
