package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listenPort, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 4221...")

	for {
		conn, err := listenPort.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		sendResponse(conn, "400 Bad Request", "Invalid request")
		return
	}

	method, path := parts[0], parts[1]
	var userAgent string
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}

		// Check if it's the User-Agent header
		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			userAgent = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}
	}

	body := `Hello World!!`

	if method == "GET" && path == "/" {
		sendResponse(conn, "200 OK", body)
	} else if method == "GET" && path == "echo/abc" {
		sendResponse(conn, "200 OK", body)
	} else if method == "GET" && path == "/user-agent" {
		sendResponse(conn, "200 OK", "User-Agent: "+userAgent)
	} else {
		sendResponse(conn, "404 Not Found", "Page Not Found\n")
	}

}

func sendResponse(conn net.Conn, status, body string) {
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Length: %d\r\nContent-Type: text/plain\r\n\r\n%s", status, len(body), body)
	conn.Write([]byte(response))

	// Write the response to the connection
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}
