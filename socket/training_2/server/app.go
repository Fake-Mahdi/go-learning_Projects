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
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("The Server is Listening")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Problem listening to the comming connection", err)
		}
		fmt.Println("Connection Was Established :", conn.RemoteAddr())
		go handleClient(conn)
	}
}
func handleClient(conn net.Conn) {
	defer conn.Close()
	go func() {
		for {
			commingMsgLength := make([]byte, 4)
			_, err := io.ReadFull(conn, commingMsgLength)
			if err != nil {
				if err == io.EOF {
					fmt.Println("The Data Has Reached Here end")
					return
				}
				fmt.Println("Problem Reading Here")
				return
			}
			fullMsgLength := binary.BigEndian.Uint32(commingMsgLength)

			realFullMsg := make([]byte, fullMsgLength)
			if _, err := io.ReadFull(conn, realFullMsg); err != nil {
				if err == io.EOF {
					fmt.Println("The Data Has Reached Here end")
					return
				}
				fmt.Println("Problem Reading Here", err)
				return
			}
			fmt.Printf("\rClient:> %s\nServer> ", string(realFullMsg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Server> ")
	for scanner.Scan() {
		serverMsg := scanner.Text()
		if serverMsg == "exit" {
			return
		}
		serverMsgLength := make([]byte, 4)
		fullServerMsg := []byte(serverMsg)
		binary.BigEndian.PutUint32(serverMsgLength, uint32(len(fullServerMsg)))

		if _, err := conn.Write(serverMsgLength); err != nil {
			fmt.Println("Enable to Write to msg Length Size to Client")
			return
		}

		_, err := conn.Write(fullServerMsg)
		if err != nil {
			fmt.Println("Enable to send Message to Client")
		}
		fmt.Print("Server> ")
	}
}
