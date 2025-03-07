package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const rootDir = "/tmp"

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
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read the request line
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

	headers := make(map[string]string)
	var contentLength int

	// Read headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}

		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) == 2 {
			key := strings.ToLower(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value

			// Capture Content-Length
			if key == "content-length" {
				contentLength, _ = strconv.Atoi(value)
			}
		}
	}

	switch method {
	case "GET":
		handleGet(conn, path, headers["user-agent"])
	case "POST":
		handlePost(conn, reader, path, contentLength)
	default:
		sendResponse(conn, "405 Method Not Allowed", "Only GET and POST are supported\n")
	}
}

func handleGet(conn net.Conn, path, userAgent string) {
	if path == "/" {
		sendResponse(conn, "200 OK", "Hello World!!")
	} else if path == "/user-agent" {
		sendResponse(conn, "200 OK", "User-Agent: "+userAgent)
	} else if strings.HasPrefix(path, "/files/") {
		filename := strings.TrimPrefix(path, "/files/")
		filePath := filepath.Join(rootDir, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			sendResponse(conn, "404 Not Found", "File not found\n")
			return
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			sendResponse(conn, "500 Internal Server Error", "Error reading file\n")
			return
		}
		sendResponse(conn, "200 OK", string(fileContent))
	} else {
		sendResponse(conn, "404 Not Found", "Page not found\n")
	}
}

func handlePost(conn net.Conn, reader *bufio.Reader, path string, contentLength int) {
	if !strings.HasPrefix(path, "/files/") {
		sendResponse(conn, "404 Not Found", "Invalid path\n")
		return
	}

	filename := strings.TrimPrefix(path, "/files/")
	filePath := filepath.Join(rootDir, filename)

	// Read the exact body length
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		sendResponse(conn, "400 Bad Request", "Error reading request body\n")
		return
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		sendResponse(conn, "500 Internal Server Error", "Error writing file\n")
		return
	}

	sendResponse(conn, "201 Created", "File created successfully\n")
}

func sendResponse(conn net.Conn, status, body string) {
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Length: %d\r\n\r\n%s", status, len(body), body)
	conn.Write([]byte(response))
}
