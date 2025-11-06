package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type User struct {
	Name      string
	Conn      net.Conn
	IsActive  string
	IpAddress string
}
type Room struct {
	Name    string
	Admin   *User
	Members []*User
}

var rooms []Room
var users []User
var mu sync.Mutex

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
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

		mu.Lock()
		remoteAdddr := conn.RemoteAddr().(*net.TCPAddr)
		ipOnly := remoteAdddr.IP.String()
		userTemplate := User{
			Name:      fmt.Sprintf("Client%d", len(users)+1),
			Conn:      conn,
			IsActive:  "Active",
			IpAddress: ipOnly,
		}
		users = append(users, userTemplate)
		fmt.Println(users)
		mu.Unlock()

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
			checkMesssage := string(realFullMsg)
			parts := strings.SplitN(checkMesssage, " ", 2)
			if parts[0] == "room" {
				handleCreateRoom(conn, parts[1])
				continue
			}
			if parts[0] == "Invite" {
				handleInviteMembers(checkMesssage)
				continue
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
		parts := strings.SplitN(serverMsg, " ", 2)

		if parts[0] == "broadcast" {
			broadcastMessage := parts[1]
			broadcastMessage_bytes := []byte(broadcastMessage)
			BroadCastMessage(conn, broadcastMessage_bytes)
			continue
		}
		if serverMsg == "DisplayAll" {
			displayActiveClient()
			continue
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

func BroadCastMessage(conn net.Conn, broadcast_message []byte) {

	mu.Lock()
	usersCopy := make([]User, len(users))
	copy(usersCopy, users)
	mu.Unlock()

	for _, user := range usersCopy {
		msgLength := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLength, uint32(len(broadcast_message)))

		if _, err := user.Conn.Write(msgLength); err != nil {
			fmt.Println("Enable to Write to msg Length Size to Client")
			continue
		}

		_, err := user.Conn.Write(broadcast_message)
		if err != nil {
			fmt.Println("Enable to send Message to Client")
			continue
		}
		fmt.Println("this is the broadcast")
		fmt.Print("Server> ")

	}

}

func handleCreateRoom(conn net.Conn, roomName string) {
	mu.Lock()
	defer mu.Unlock()
	for _, room := range rooms {
		if room.Name == roomName {
			fmt.Println("This is Room already exit chose an other name")
			return
		}
	}
	for i := range users {
		if users[i].Conn == conn {
			roomMap := Room{Name: roomName, Admin: &users[i], Members: []*User{&users[i]}}
			rooms = append(rooms, roomMap)
			fmt.Printf("The Room With Room Name %s was created \n", roomName)
			break
		}

	}
	fmt.Print("Server> ")
}

func handleInviteMembers(listOfMembers string) {
	splitcheck := strings.Split(listOfMembers, " ")
	if len(splitcheck) < 3 {
		fmt.Println("UNFOUNDED COMMAND")
		return
	}

	splitList := strings.SplitN(listOfMembers, " ", 3)
	ipsPart := strings.Split(splitList[1], ",")
	roomName := splitList[2]
	var ptrRoom *Room
	for i := range rooms {
		if rooms[i].Name == roomName {
			ptrRoom = &rooms[i]
			break
		}
	}

	if ptrRoom == nil {
		fmt.Println("The Room Does not Exist")
		return
	}

	for ipsIndex := range ipsPart {
		var userFound *User
		for userIndex := range users {
			if users[userIndex].IpAddress == ipsPart[ipsIndex] {
				userFound = &users[userIndex]
				break
			}
		}

		if userFound == nil {
			fmt.Println("There is no such user found")
			continue // <- i think continue is better like if user is not found i this iteration we can just go to next iteration
		}
		var alreadyInRoom bool
		for _, user := range ptrRoom.Members { //<= i dont understand why should i change like ipsIndex i mean each index is scope to a loop ?
			if user.IpAddress == ipsPart[ipsIndex] { // i have no ip variable i used ipsPart wich contain all ips
				alreadyInRoom = true
				break
			}
		}
		if alreadyInRoom == true {
			continue
		}
		ptrRoom.Members = append(ptrRoom.Members, userFound)
	}

}

func displayActiveClient() {
	mu.Lock()
	defer mu.Unlock()

	if len(users) == 0 {
		fmt.Println("No connected users.")
		return
	}

	fmt.Println("┌───────────────────────────────┐")
	fmt.Println("│        Connected Users        │")
	fmt.Println("├─────────────┬───────────┬───────────────┤")
	fmt.Printf("│ %-11s │ %-9s │ %-30s │\n", "Name", "Status", "IP Address")
	fmt.Println("├─────────────┼───────────┼───────────────┤")

	for _, user := range users {
		fmt.Printf("│ %-11s │ %-9s │ %-30s │\n", user.Name, user.IsActive, user.IpAddress)
	}

	fmt.Println("└─────────────┴───────────┴───────────────┘")
}
