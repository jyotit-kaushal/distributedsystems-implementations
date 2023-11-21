package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

// setting all the parameters- can be changed based on what you want to check

const (
	tnmtbsbc int     = 3  //total number of messages to be sent by client
	tnoc     int     = 10 // total number of clients
	ibem     float32 = 2  // interval between each message
	mddtn    float32 = 2  // max delay due to network
)

// setting basic classes
type server struct {
	channel     chan message
	clientArray []client
	logicalTS   int
}

type message struct {
	sender    string
	receiver  string
	text      string
	logicalTS int
}

type client struct {
	name          string
	clientChannel chan message
	server        server
	logicalTS     int
	readyChannel  chan int
	killChan      chan int
}

type logTS struct {
	numTS int
	mux   sync.Mutex
}

func main() {

	// initializing server and all clients
	server := instantiateServer()

	for i := 1; i <= tnoc; i++ {

		client := instantiateClient(fmt.Sprintf("Client %d", i), *server)
		server.clientArray = append(server.clientArray, *client)

	}

	var parentChannel chan message = make(chan message, tnoc*tnoc*tnmtbsbc*10)

	for _, client := range server.clientArray {

		go client.clientSendRoutine()              // starting client routine that puts itself in the ready channel indicating it's going to send messages
		go client.clientMainRoutine(parentChannel) // main client routine that actually sends the messages and also listens for broadcasted msgs from server and prints them out

	}

	server.listen(parentChannel)

	allMessages := []message{}

	for _, client := range server.clientArray {
		client.killChan <- 1
	}

	close(parentChannel)

	for message := range parentChannel {
		allMessages = append(allMessages, message)
	}

	// sorting the messages based on their logical clocks
	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].logicalTS < allMessages[j].logicalTS
	})

	fmt.Println("\n|XXXX| Total Order of All Messages Below: |XXXX|")

	for i, message := range allMessages {
		fmt.Printf("|XX| %d. |XX| Logical Clock: %d |XX| Message: %s |XX| Received By: %s |XX| Sent By: %s |XX| \n", i, message.logicalTS, message.text, message.receiver, message.sender)
	}

}

func instantiateServer() *server {
	channel := make(chan message)
	clientArray := []client{}

	s := server{channel, clientArray, 0}

	return &s
}

func broadcastMessage(clientChannel chan message, msg message) {
	var numSeconds float32 = rand.Float32() * mddtn
	time.Sleep(time.Duration(numSeconds) * time.Second)
	clientChannel <- msg

	return
}

// function that handles coin toss and is basically the server listening for msgs from clients
func (s server) listen(parentChannel chan message) {

	var activeclientCount int = tnoc

	for {
		msg := <-s.channel

		s.logicalTS = max(s.logicalTS, msg.logicalTS) + 1
		msg.logicalTS = s.logicalTS

		parentChannel <- msg

		fmt.Printf("|XX| Logical Clock: %d |XX| Received By: Server   |XX| Message: %s |XX| Sent By: %s |XX| \n", s.logicalTS, msg.text, msg.sender)

		var broadcast int = rand.Intn(2)

		if broadcast == 0 {
			fmt.Printf("\n|XXXX| Server : Not Broadcasting '%s' from %s |XXXX| \n \n", msg.text, msg.sender)
		} else {
			s.logicalTS += 1

			fmt.Printf("\n|XXXX| Server : Broadcasting '%s' from %s |XXXX| \n \n", msg.text, msg.sender)

			for _, client := range s.clientArray {
				if client.name == msg.sender {
					continue
				} else {
					newMessagestr := "<<Broadcasted>> " + msg.text
					newMessage := message{"Server", client.name, newMessagestr, s.logicalTS}
					go broadcastMessage(client.clientChannel, newMessage)
				}

			}
		}

		if strings.Contains(msg.text, "Last Hello") {
			activeclientCount -= 1

			if activeclientCount == 0 {
				time.Sleep(time.Duration(mddtn) * time.Second)

				fmt.Println("|XXXX| All clients done sending messages |XXXX|")

				return
			}

		}
	}
}

func (s *server) registerClient(c client) []client {
	fmt.Printf("Registering %s \n", c.name)
	s.clientArray = append(s.clientArray, c)
	return s.clientArray
}

func instantiateClient(name string, s server) *client {

	clientChannel := make(chan message)
	readyChannel := make(chan int)
	killChan := make(chan int)
	c := client{name, clientChannel, s, 0, readyChannel, killChan}
	c.server = s
	s.clientArray = s.registerClient(c)
	return &c

}

func (c client) clientSendRoutine() {

	for i := 1; i <= tnmtbsbc; i++ {
		time.Sleep(time.Duration(ibem) * time.Second)
		c.readyChannel <- i
	}

	return
}

func (c client) clientMainRoutine(parentChannel chan message) {
	var nextMessage message

	for {
		select {

		case messagefromServer := <-c.clientChannel:

			c.logicalTS = max(c.logicalTS, messagefromServer.logicalTS) + 1
			messagefromServer.logicalTS = c.logicalTS
			parentChannel <- messagefromServer
			fmt.Printf("|XX| Logical Clock: %d |XX| Received By: %s |XX| Message: %s |XX| Sent By: %s |XX| \n", c.logicalTS, c.name, messagefromServer.text, messagefromServer.sender)

		case nextmessageNo := <-c.readyChannel:

			c.logicalTS += 1

			if nextmessageNo == tnmtbsbc {
				nextMessage = message{c.name, "Server", fmt.Sprintf("Last Hello from %s", c.name), c.logicalTS}
			} else {
				nextMessage = message{c.name, "Server", fmt.Sprintf("Hello %d from %s", nextmessageNo, c.name), c.logicalTS}
			}
			c.server.channel <- nextMessage

		case <-c.killChan:
			return

		default:

		}
	}
}
