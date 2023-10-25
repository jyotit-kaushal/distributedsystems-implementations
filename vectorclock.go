package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	tnmtbsbc int     = 3  //total number of messages to be sent by client
	tnoc     int     = 10 // total number of clients
	ibem     float32 = 2  // interval between each message
	mddtn    float32 = 10 // max delay due to network
	tnop     int     = 11 // total number of processes
)

type server struct {
	pid         int
	channel     chan message
	clientArray []client
	vectorClock [tnop]int
}

type message struct {
	sender      string
	receiver    string
	text        string
	vectorClock [tnop]int
}

type client struct {
	pid           int
	name          string
	clientChannel chan message
	server        server
	readyChannel  chan int
	killChan      chan int
	vectorClock   [tnop]int
}

func main() {

	server := instantiateServer(0)

	for i := 1; i <= tnoc; i++ {

		client := instantiateClient(i, fmt.Sprintf("Client %d", i), *server)
		server.clientArray = append(server.clientArray, *client)

	}

	var parentChannel chan message = make(chan message, tnoc*tnoc*tnmtbsbc*10)

	var pcvChannel chan string = make(chan string, tnoc*tnoc*tnmtbsbc)

	for _, client := range server.clientArray {

		go client.clientSendRoutine()
		go client.clientMainRoutine(parentChannel, pcvChannel)

	}

	server.listen(parentChannel, pcvChannel)

	allMessages := []message{}

	for _, client := range server.clientArray {
		client.killChan <- 1
	}

	close(parentChannel)
	close(pcvChannel)

	for message := range parentChannel {
		allMessages = append(allMessages, message)
	}

	sort.Slice(allMessages, func(i, j int) bool {
		return sortingHelper(allMessages[i].vectorClock, allMessages[j].vectorClock)
	})

	fmt.Println("\n|XXXX| Total Order of All Messages Below: |XXXX|\n")

	for i, message := range allMessages {
		fmt.Printf("|XX| %d. |XX| Logical Clock: %d |XX| Message: %s |XX| Received By: %s |XX| Sent By: %s |XX| \n", i, message.vectorClock, message.text, message.receiver, message.sender)
	}

	allPCVs := []string{}
	for pcv := range pcvChannel {
		allPCVs = append(allPCVs, pcv)
	}

	fmt.Println("\n|XXXX| All Potential Causality Violations: |XXXX|\n")

	for _, message := range allPCVs {
		fmt.Println(message)
	}

}

func sortingHelper(v1 [tnop]int, v2 [tnop]int) bool {

	for i, _ := range v1 {
		if v1[i] > v2[i] {
			return false
		}
	}
	return true

}

func combineClocks(vectorClock1 [tnop]int, vectorClock2 [tnop]int, receiverPID int) [tnop]int {

	res := [tnop]int{}

	for i, _ := range vectorClock1 {
		res[i] = max(vectorClock2[i], vectorClock1[i])
	}

	res[receiverPID] += 1
	return res
}

func pcvCheck(vectorClock1 [tnop]int, vectorClock2 [tnop]int) bool {

	for i, _ := range vectorClock1 {

		if vectorClock1[i] <= vectorClock2[i] {
			return false
		}
	}
	return true
}

func instantiateServer(pid int) *server {
	channel := make(chan message)
	clientArray := []client{}
	vectorClock := [tnop]int{}

	s := server{pid, channel, clientArray, vectorClock}

	return &s
}

func broadcastMessage(clientChannel chan message, msg message) {
	var numSeconds float32 = rand.Float32() * mddtn
	time.Sleep(time.Duration(numSeconds) * time.Second)
	clientChannel <- msg

	return
}

func (s server) listen(parentChannel chan message, pcvChannel chan string) {

	var activeclientCount int = tnoc

	for {
		msg := <-s.channel

		s.vectorClock = combineClocks(s.vectorClock, msg.vectorClock, s.pid)

		if pcvCheck(s.vectorClock, msg.vectorClock) {
			fmt.Println("\n|XXXX| DETECTED POTENTIAL CAUSALITY VIOLATION |XXXX|\n")
			pcvChannel <- fmt.Sprintf("|XX| From: %s |XX| To: %s |XX| Message Vector Clock: %d |XX| Local Vector Clock: %d |XX| Message: %s |XX|",
				msg.sender, msg.receiver, msg.vectorClock, s.vectorClock, msg.text)
		}

		msg.vectorClock = s.vectorClock

		parentChannel <- msg

		fmt.Printf("|XX| Vector Clock: %d |XX| Received By: Server   |XX| Message: %s |XX| Sent By: %s |XX| \n", s.vectorClock, msg.text, msg.sender)

		var broadcast int = rand.Intn(2)

		if broadcast == 0 {
			fmt.Printf("\n|XXXX| Server : Not Broadcasting '%s' from %s |XXXX| \n \n", msg.text, msg.sender)
		} else {
			s.vectorClock[s.pid] += 1

			fmt.Printf("\n|XXXX| Server : Broadcasting '%s' from %s |XXXX| \n \n", msg.text, msg.sender)

			for _, client := range s.clientArray {
				if client.name == msg.sender {
					continue
				} else {
					newMessagestr := "<<Broadcasted>> " + msg.text
					newMessage := message{"Server", client.name, newMessagestr, s.vectorClock}
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

// Constructor for client
func instantiateClient(pid int, name string, s server) *client {

	clientChannel := make(chan message)
	readyChannel := make(chan int)
	killChan := make(chan int)
	vectorClock := [tnop]int{}
	c := client{pid, name, clientChannel, s, readyChannel, killChan, vectorClock}
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

func (c client) clientMainRoutine(parentChannel chan message, pcvChannel chan string) {

	var nextMessage message

	for {
		select {

		case messagefromServer := <-c.clientChannel:

			c.vectorClock = combineClocks(c.vectorClock, messagefromServer.vectorClock, c.pid)

			if pcvCheck(c.vectorClock, messagefromServer.vectorClock) {
				fmt.Println("\n|XXXX| DETECTED POTENTIAL CAUSALITY VIOLATION |XXXX|\n")
				pcvChannel <- fmt.Sprintf("|XX| From: %s |XX| To: %s |XX| Message Vector Clock: %d |XX| Local Vector Clock: %d |XX| Message: %s |XX|",
					messagefromServer.sender, messagefromServer.receiver, messagefromServer.vectorClock, c.vectorClock, messagefromServer.text)
			}

			messagefromServer.vectorClock = c.vectorClock

			parentChannel <- messagefromServer

			fmt.Printf("|XX| Vector Clock: %d |XX| Received By: %s |XX| Message: %s |XX| Sent By: %s |XX| \n", c.vectorClock, c.name, messagefromServer.text, messagefromServer.sender)

		case nextmessageNo := <-c.readyChannel:

			c.vectorClock[c.pid] += 1

			if nextmessageNo == tnmtbsbc {
				nextMessage = message{c.name, "Server", fmt.Sprintf("Last Hello from %s", c.name), c.vectorClock}
			} else {
				nextMessage = message{c.name, "Server", fmt.Sprintf("Hello %d from %s", nextmessageNo, c.name), c.vectorClock}
			}
			c.server.channel <- nextMessage

		case <-c.killChan:
			return
		default:

		}
	}
}
