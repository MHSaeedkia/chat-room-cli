# CLI Chat Room

## Overview
This is a simple Command-Line Interface (CLI) chat room application built with Go. The project allows users to run as either a **server** or a **client**, with messaging handled via [NATS](https://nats.io/). The application also includes a health check endpoint for monitoring the server's status.

---

## Features
- **Server and Client Roles**:  
  Upon running the application, you can choose your role as either a server or a client. Only one server is allowed at a time.

- **NATS for Messaging**:  
  The project uses NATS as the messaging tool for communication between clients and the server.

- **Commands in Chat Room**:  
  You can use predefined commands in the chat room by prefixing your text with `#`.  

  ### Supported Commands:
  - `#users`: Displays the list of online users.

- **Health Check Endpoint**:  
  A REST API endpoint is available to check if a server is running.  
  **Endpoint**: `/api/v1/health`

- **Custom Command Syntax**:  
  To send commands in the chat room, always prefix the message with `#` (e.g., `#users`). Messages without `#` will be sent as normal text.

---

## Getting Started

### Prerequisites
1. **Go**: Ensure you have Go installed (version 1.20+ recommended).  
   [Download Go here](https://go.dev/dl/).

2. **NATS Server**:  
   The application requires a running NATS server. You can start one locally using the following command:
   ```bash
   nats-server -p 4222
   ```

# Installation

### Clone the repository:
```bash
git clone https://github.com/your-username/your-repo-name.git

cd your-repo-name
```

# Run the application
```bash
go run . 
```

# Usage

## Running the Application

### Choose Your Role:
When you run the application, you'll be prompted to choose between **server** or **client** roles:
- **Server**: Hosts the chat room and handles client connections.
- **Client**: Joins the chat room to send/receive messages.

### Commands in Chat Room:
Use `#` before text to issue commands. For example:
- `#users`: Displays a list of online users.

### Sending Messages:
Simply type your message and hit `Enter` to send it to the chat room.

### Error Handling:
If multiple servers are detected, the application will terminate with an error message, as only one server is allowed at a time.
If server go down , all client will be down.



