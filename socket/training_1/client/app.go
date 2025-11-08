package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// 1️⃣ Connect to the server and return the connection
func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to server!")
	return conn
}

// 2️⃣ Handle server communication: read & write
func handleServer(conn net.Conn) {
	// Goroutine for reading messages from server
	go func() {
		for {
			lenBuf := make([]byte, 4)
			_, err := io.ReadFull(conn, lenBuf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server disconnected")
					os.Exit(0)
				}
				fmt.Println("Error reading:", err)
				continue
			}

			msgLen := binary.BigEndian.Uint32(lenBuf)
			msg := make([]byte, msgLen)
			_, err = io.ReadFull(conn, msg)
			if err != nil {
				fmt.Println("Error reading full message:", err)
				continue
			}

			fmt.Printf("\rServer:> %s\nClient> ", string(msg))
		}
	}()

	// Loop to read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Client> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Closing client...")
			break
		}

		msgBytes := []byte(text)
		lenBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBuf, uint32(len(msgBytes)))

		if _, err := conn.Write(lenBuf); err != nil {
			fmt.Println("Error sending length:", err)
			continue
		}

		if _, err := conn.Write(msgBytes); err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}
		fmt.Println("Client>")
	}
}

func main() {
	conn := connectToServer()
	defer conn.Close() // Close connection when main exits
	handleServer(conn)
}
