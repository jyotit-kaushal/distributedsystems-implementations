package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// defining all the constants and data structures

const (
	Request     msgSignature = "Request"
	Reply       msgSignature = "Reply"
	Release     msgSignature = "Release"
	setnumNodes int          = 10
)

type Node struct {
	id             int
	logicalClock   int
	nodeChannel    chan Message
	nodePQ         []Message
	replyCheckList map[int]*Node
	replyHistory   map[int]map[int]bool
	repliesPending map[int][]Message
	mainWaitGroup  *sync.WaitGroup
	quit           chan int
}

type msgSignature string

type Message struct {
	msgSignature msgSignature
	message      string
	senderID     int
	timestamp    int
	receiver     Receiver
}

type Receiver struct {
	targetID  int
	timestamp int
}

// helper functions

func instantiateNode(id int) *Node {

	nodeChannel := make(chan Message)
	var nodePQ []Message
	var replyHistory = map[int]map[int]bool{}

	node := Node{id, 0, nodeChannel, nodePQ, nil,
		replyHistory, map[int][]Message{}, &sync.WaitGroup{}, make(chan int)}

	return &node
}

func toString(nodePQ []Message) string {

	if len(nodePQ) == 0 {
		return "Priority Queue: |X| Empty |X|"
	}

	var res string = "|X| Priority Queue: "

	for _, msg := range nodePQ {
		res += fmt.Sprintf("Request at [TS: %d] by Node %d |X| ", msg.timestamp, msg.senderID)
	}
	return res
}

func (node *Node) createReplyMap() map[int]bool {

	rmap := map[int]bool{}
	for i, _ := range node.replyCheckList {
		if i == node.id {
			continue
		}
		rmap[i] = false
	}

	return rmap
}

func (node *Node) dequeue(senderID int) {

	for i, msg := range node.nodePQ {

		if msg.senderID == senderID {
			if i >= len(node.nodePQ)-1 {
				node.nodePQ = node.nodePQ[:i]
			} else {
				node.nodePQ = append(node.nodePQ[:i], node.nodePQ[i+1:]...)
			}
			return
		}
	}
}

func (node *Node) updatereplyCheckList(replyCheckList map[int]*Node) {
	node.replyCheckList = replyCheckList
}

func (node *Node) enqueue(message Message) {

	for _, msg := range node.nodePQ {
		if message.senderID == msg.senderID {
			return
		}
	}

	node.nodePQ = append(node.nodePQ, message)

	sort.SliceStable(node.nodePQ, func(i, j int) bool {

		if node.nodePQ[i].timestamp < node.nodePQ[j].timestamp {
			return true
		} else if node.nodePQ[i].timestamp == node.nodePQ[j].timestamp && node.nodePQ[i].senderID < node.nodePQ[j].senderID {
			return true
		} else {
			return false
		}
	})

	fmt.Printf("|XX| Node %d |XX| Queue Status: %s |XX| \n", node.id, toString(node.nodePQ))
}

func (node *Node) sendAll(msg Message) {

	for nodeId, _ := range node.replyCheckList {
		if nodeId == node.id {
			continue
		}
		go node.sendmessage(msg, nodeId)
	}

}

func (node *Node) sendmessage(msg Message, receiverID int) {

	fmt.Printf("|XX| Node %d sending a |XX| %s message |XX| to Node %d |XX| at MemAddr %p |XX|\n", node.id, msg.msgSignature, receiverID, node.replyCheckList[receiverID])

	randomLatency := rand.Intn(1000) + 2000
	time.Sleep(time.Duration(randomLatency) * time.Millisecond)

	receiver := node.replyCheckList[receiverID]

	receiver.nodeChannel <- msg

}

func (node *Node) checkallreplies(timestamp int) bool {

	if _, check := node.replyHistory[timestamp]; !check {
		return false
	}

	for _, replyStatus := range node.replyHistory[timestamp] {
		if replyStatus == false {
			return false
		}
	}
	return true
}

// base node functionality related functions

func (node *Node) enterCSRequest() {

	node.logicalClock += 1

	if setnumNodes == 1 {
		node.enterCS(Message{
			msgSignature: "Request",
			message:      "",
			senderID:     node.id,
			timestamp:    node.logicalClock,
			receiver:     Receiver{},
		})
	}

	fmt.Printf("\n|XX| Node %d |XX| is requesting to enter Critical Section |XX| \n", node.id)

	time.Sleep(time.Duration(500) * time.Millisecond)

	req := Message{
		msgSignature: "Request",
		message:      "",
		senderID:     node.id,
		timestamp:    node.logicalClock,
	}

	node.enqueue(req)

	otherNodes := map[int]bool{}

	for nodeId, _ := range node.replyCheckList {
		if nodeId == node.id {
			continue
		}
		otherNodes[nodeId] = false
	}

	node.replyHistory[node.logicalClock] = otherNodes
	node.sendAll(req)
}

func (node *Node) enterCS(msg Message) {

	defer node.mainWaitGroup.Done()

	node.dequeue(msg.senderID)

	numSeconds := int(math.Mod(float64(node.id), 3))

	fmt.Printf("\n|XX| Node %d |XX| is entering the Critical Section for |XX| %d seconds |XX| for Message with |XX| Timestamp %d |XX| \n \n", node.id, numSeconds, msg.timestamp)
	time.Sleep(time.Duration(numSeconds) * time.Second)
	fmt.Printf("\n|XX| Node %d |XX| is now exiting the Critical Section |XX|\n \n", node.id)

	node.logicalClock += 1

	release := Message{
		msgSignature: "Release",
		message:      "",
		senderID:     node.id,
		timestamp:    node.logicalClock,
		receiver:     Receiver{},
	}
	node.sendAll(release)

}

func (node *Node) sendReply(msg Message) {

	fmt.Printf("|XX| Node %d |XX| sending reply to |XX| Node %d |XX|\n", node.id, msg.senderID)

	node.logicalClock += 1

	reply := Message{
		msgSignature: "Reply",
		message:      "",
		senderID:     node.id,
		timestamp:    node.logicalClock,
		receiver: Receiver{
			targetID:  msg.senderID,
			timestamp: msg.timestamp,
		},
	}

	node.sendmessage(reply, msg.senderID)
}

// functions to handle different messages

func (node *Node) handleRequest(msg Message) {

	var replied bool = false

	if len(node.replyHistory) == 0 {
		go node.sendReply(msg)
		replied = true
	}

	for requestTS, replyMap := range node.replyHistory {

		if requestTS < msg.timestamp {

			if replyMap[msg.senderID] {
				go node.sendReply(msg)
				replied = true
			}

		} else if requestTS == msg.timestamp && node.id < msg.senderID {

			if replyMap[msg.senderID] {
				go node.sendReply(msg)
				replied = true
			}
		} else {
			go node.sendReply(msg)
			replied = true
		}
	}

	if replied == false {
		node.repliesPending[msg.senderID] = append(node.repliesPending[msg.senderID], msg)
	} else {
		node.enqueue(msg)
	}

	fmt.Printf("|XX| Node %d |XX| Queue Status: %s |XX| \n", node.id, toString(node.nodePQ))
}

func (node *Node) handleReply(msg Message) {

	ts := msg.receiver.timestamp

	if _, check := node.replyHistory[ts]; !check {
		node.replyHistory[ts] = node.createReplyMap()
	}

	node.replyHistory[msg.receiver.timestamp][msg.senderID] = true

	for _, reqMsg := range node.repliesPending[msg.senderID] {
		node.handleRequest(reqMsg)
	}

	if node.checkallreplies(msg.receiver.timestamp) {

		fmt.Printf("\n|XX| Node %d has received all replies for Request at Time Stamp: %d |XX|\n", node.id, msg.receiver.timestamp)

		requestatHead := node.nodePQ[0]

		if requestatHead.senderID == node.id && requestatHead.timestamp == msg.receiver.timestamp {

			fmt.Printf("|XX| With Node %d at the head, it is all good to enter Critical Section |XX| \n", node.id)

			delete(node.replyHistory, msg.receiver.timestamp)
			node.enterCS(requestatHead)
		}

	}
}

func (node *Node) handleRelease(msg Message) {

	node.dequeue(msg.senderID)

	if len(node.nodePQ) > 0 {
		requestatHead := node.nodePQ[0]

		if requestatHead.senderID == node.id {
			if node.checkallreplies(node.nodePQ[0].timestamp) {

				delete(node.replyHistory, node.nodePQ[0].timestamp)
				node.enterCS(requestatHead)
			}
		}
	}
}

func (node *Node) handleReceive(msg Message) {

	time.Sleep(time.Duration(50) * time.Millisecond)

	if msg.timestamp >= node.logicalClock {
		node.logicalClock = msg.timestamp + 1
	} else {
		node.logicalClock += 1
	}

	fmt.Printf("|XX| Node %d received a |XX| %s Message |XX| from Node %d |XX| \n", node.id, msg.msgSignature, msg.senderID)

	if msg.msgSignature == "Request" {
		node.handleRequest(msg)
	} else if msg.msgSignature == "Reply" {
		node.handleReply(msg)
	} else if msg.msgSignature == "Release" {
		node.handleRelease(msg)
	} else {
		fmt.Println("Invalid Message Signature")
	}
}

func (node *Node) listen() {

	for {
		select {
		case msg := <-node.nodeChannel:
			go node.handleReceive(msg)
		case <-node.quit:
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup

	globalNodeMap := map[int]*Node{}

	for i := 1; i <= setnumNodes; i++ {
		newNode := instantiateNode(i)
		globalNodeMap[i] = newNode

	}

	for i := 1; i <= setnumNodes; i++ {
		wg.Add(1)
		globalNodeMap[i].updatereplyCheckList(globalNodeMap)
	}

	for i := 1; i <= setnumNodes; i++ {
		globalNodeMap[i].mainWaitGroup = &wg
		go globalNodeMap[i].listen()
	}

	start := time.Now()

	for i := 1; i <= setnumNodes; i++ {
		go globalNodeMap[i].enterCSRequest()
	}

	wg.Wait()
	t := time.Now()
	time.Sleep(time.Duration(3) * time.Second)

	fmt.Printf("\nTime Taken with %d nodes : %.2f seconds \n", setnumNodes, t.Sub(start).Seconds())
	fmt.Println("Simulation Ending \n")

	for i := 1; i <= setnumNodes; i++ {
		globalNodeMap[i].quit <- 1
	}
}
