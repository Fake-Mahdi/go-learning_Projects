package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", ":8080") // <= this is the server listner mean type of connection plus port that will listen on
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("server started on port 8080...")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	go func() {
		for {
			lenBuf := make([]byte, 4)
			_, err := io.ReadFull(conn, lenBuf)
			if err != nil {
				if err == io.EOF {
					return
				}
				fmt.Println("error reading : ", err)
				break
			}
			msgLen := binary.BigEndian.Uint32(lenBuf)

			fullMessage := make([]byte, msgLen)
			if _, err := io.ReadFull(conn, fullMessage); err != nil {
				fmt.Println("This is a Huge error : ", err)
				return
			}

			fmt.Printf("\rClient:> %s\nServer> ", string(fullMessage))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin) // <= this is create a buffer using stream input
	fmt.Println("Server> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			print("Have a good day")
			return
		}

		// convert the message into bytes then slice the 4 bytes for length getting
		msgByte := []byte(text)
		msgLen := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLen, uint32(len(msgByte)))

		// sending data length of data into the client
		if _, err := conn.Write(msgLen); err != nil {
			fmt.Println("Error Sending Length : ", err)
			return
		}
		// finally sending the data into the client
		_, err := conn.Write(msgByte)
		if err != nil {
			fmt.Println("error sending Message :", err)
			return
		}
		fmt.Println("Server> ")
	}
}
