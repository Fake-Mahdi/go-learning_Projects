package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func StartConnection() net.Conn {
	conn, err := net.Dial("tcp", "192.168.1.2:8080")
	if err != nil {
		fmt.Print("This Something Bad Happen or connection problem")
		panic(err)
	}
	return conn
}
func handleServer(conn net.Conn) {
	stop := make(chan struct{})
	msgChan := make(chan []byte)
	defer conn.Close()
	go func() {
		for {
			commingMsgLength := make([]byte, 4)
			_, err := io.ReadFull(conn, commingMsgLength)
			if err != nil {
				fmt.Println("Problem connecting with the server :", err)
				close(stop)
				return
			}
			fullComingMsg := binary.BigEndian.Uint32(commingMsgLength)
			msg := make([]byte, fullComingMsg)
			if _, err := io.ReadFull(conn, msg); err != nil {
				fmt.Println("Problem Reading The full Real Message")
				close(stop)
				return
			}
			select {
			case msgChan <- msg:
			case <-stop:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Have A Good Day Sir")
				return
			case serverMsg := <-msgChan:
				fmt.Printf("\rServer:> %s\nClient> ", string(serverMsg))
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Client> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Closing client...")
			close(stop)
			return
		}

		msgBytes := []byte(text)
		lenBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBuf, uint32(len(msgBytes)))

		if _, err := conn.Write(lenBuf); err != nil {
			fmt.Println("Error sending length:", err)
			close(stop)
			conn.Close()
		}

		if _, err := conn.Write(msgBytes); err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}
		fmt.Print("Client>")
	}
}
func main() {
	conn := StartConnection()
	defer conn.Close()
	handleServer(conn)
}
