package main

import (
	"fmt"
	"net"
	"os"
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

	// Read the incoming HTTP request (optional, for debugging)
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// Construct a basic HTTP response with status code 200 OK
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 13\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello, World!"

	// Send the response to the client
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}
