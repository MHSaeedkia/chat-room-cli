package cmd

import (
	"chat-room-cli/cmd/client"
	"chat-room-cli/cmd/server"
	"fmt"
	"net/http"
	"strings"
)

const (
	HttpPort = "12345"

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

// this function check , if is any server up or not .
func checkServer(port string) bool {
	url := fmt.Sprintf("http://localhost:%s/api/v1/health", port)
	response, err := http.Get(url)

	if err != nil {
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}

func Run() error {
	serverStatus := checkServer(HttpPort)

	// starting menu
	role := menu()
	switch role {
	case SERVER:
		if !serverStatus {
			return server.Run()
		}
		fmt.Printf("There is another server here ; only one server is allowed in this chat room\n")
		return client.Run()
	case CLIENT:
		if !serverStatus {
			fmt.Printf("Server is down !\n")
			return fmt.Errorf("Server is Down")
		}
		return client.Run()
	}
	return nil
}
