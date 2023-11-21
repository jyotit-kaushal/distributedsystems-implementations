package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// functions to define the constants and data structures

const (
	Request     = "request"
	Vote        = "reply"
	Release     = "release"
	Rescind     = "rescind"
	setnumNodes = 10
)

type Receiver struct {
	targetId  int
	timestamp [setnumNodes]int
}

type Message struct {
	timestamp    [setnumNodes]int
	senderId     int
	msgSignature string
	replyTo      Receiver
	status       ExecStatus
}

type ExecStatus string

const (
	REQ  ExecStatus = "req"
	EXE  ExecStatus = "exec"
	DONE ExecStatus = "done"
)

type Node struct {
	id             int
	nodeWG         *sync.WaitGroup
	nodeClock      [setnumNodes]int
	nodeChannel    chan Message
	nodePQ         []Message
	replyCheckList map[int]*Node
	replyHistory   map[string]map[int]bool
	status         map[string]ExecStatus
	vote           Message
}

// helper functions

func checkforEarlierTS(t1 [setnumNodes]int, t2 [setnumNodes]int) bool {

	for i, _ := range t1 {
		if t2[i] < t1[i] {
			return false
		}
	}
	return true
}

func checkifConcurrent(t1 [setnumNodes]int, t2 [setnumNodes]int) bool {

	t1inFront := false
	t2inFront := false

	for i, _ := range t1 {
		if t1[i] > t2[i] {
			t1inFront = true
		}
		if t2[i] > t1[i] {
			t2inFront = true
		}
		if t1inFront && t2inFront {
			return true
		}
	}
	return false
}

func tostringTS(timestamps [setnumNodes]int) string {

	res := ""

	for _, timestamp := range timestamps {
		res += strconv.Itoa(timestamp) + " "
	}
	return res
}

func combineClocks(clientid int, t1 [setnumNodes]int, t2 [setnumNodes]int) [setnumNodes]int {
	res := [setnumNodes]int{}

	for i, v := range t1 {
		if v > t2[i] {
			res[i] = v
		} else {
			res[i] = t2[i]
		}
	}

	res[clientid] += 1
	return res
}

func instantiateNode(id int) *Node {

	nodeChannel := make(chan Message)
	pq := []Message{}
	replyHistory := make(map[string]map[int]bool)

	node := Node{id, &sync.WaitGroup{}, [setnumNodes]int{}, nodeChannel, pq, map[int]*Node{}, replyHistory, map[string]ExecStatus{}, Message{}}

	return &node
}

func instantiateMessage(msgSignature string, senderId int, timestamp [setnumNodes]int) *Message {

	msg_receiver := Receiver{-1, [setnumNodes]int{}}
	msg := Message{timestamp, senderId, msgSignature, msg_receiver, REQ}

	return &msg
}

func (node *Node) dequeue(senderId int) {

	for i, msg := range node.nodePQ {

		if msg.senderId == senderId {
			if i < len(node.nodePQ)-1 {
				node.nodePQ = append(node.nodePQ[:i], node.nodePQ[i+1:]...)
			} else {
				node.nodePQ = node.nodePQ[:i]
			}
			return
		}
	}
}

func tostringPQ(pq []Message) string {

	if len(pq) == 0 {
		return "Priority Queue: |X| Empty |X|"
	}

	var res string = "|X| Priority Queue: "

	for _, msg := range pq {
		res += fmt.Sprintf("Request at [TS: %d] by Node %d |X| ", msg.timestamp, msg.senderId)
	}
	return res
}

func (node *Node) enqueue(message Message) {

	for _, msg := range node.nodePQ {
		if message.senderId == msg.senderId {
			return
		}
	}

	node.nodePQ = append(node.nodePQ, message)

	sort.SliceStable(node.nodePQ, func(i, j int) bool {
		if checkforEarlierTS(node.nodePQ[i].timestamp, node.nodePQ[j].timestamp) {
			return true
		} else if node.nodePQ[j].senderId < node.nodePQ[i].senderId && checkifConcurrent(node.nodePQ[i].timestamp, node.nodePQ[j].timestamp) {
			return true
		} else {
			return false
		}
	})

	fmt.Printf("|XX| Node %d |XX| Queue Status: %s |XX| \n", node.id, tostringPQ(node.nodePQ))
}

func (node *Node) send(msg Message, receiverId int) {

	fmt.Printf("|XX| Node %d sending a |XX| %s message |XX| to Node %d |XX| at MemAddr %p |XX|\n", node.id, msg.msgSignature, receiverId, node.replyCheckList[receiverId])

	randomLatency := rand.Intn(1000) + 2000

	time.Sleep(time.Duration(randomLatency) * time.Millisecond)

	receiver := node.replyCheckList[receiverId]

	receiver.nodeChannel <- msg
}

func (node *Node) reply(msg Message) {

	fmt.Printf("|XX| Node %d |XX| sending reply to |XX| Node %d |XX|\n", node.id, msg.senderId)

	node.nodeClock = combineClocks(node.id, node.nodeClock, msg.timestamp)
	replyMsg := instantiateMessage(Vote, node.id, node.nodeClock)

	replyMsg.replyTo.targetId = msg.senderId
	replyMsg.replyTo.timestamp = msg.timestamp

	node.vote = msg

	node.dequeue(msg.senderId)
	node.send(*replyMsg, msg.senderId)
}

func (node *Node) sendAll(msg Message) {

	for id, _ := range node.replyCheckList {
		if id == node.id {
			continue
		}
		go node.send(msg, id)
	}
}

func (node *Node) createreplyHistory() map[int]bool {

	res := map[int]bool{}

	for i, _ := range node.replyCheckList {
		if i == node.id {
			continue
		}
		res[i] = false
	}
	return res
}

// base node functionality functions

func (node *Node) enterCS(msg Message) {

	defer node.nodeWG.Done()

	timestamp := tostringTS(msg.timestamp)

	node.status[timestamp] = EXE
	node.dequeue(msg.senderId)

	numSeconds := int(math.Mod(float64(node.id), 3))

	fmt.Printf("\n|XX| Node %d |XX| is entering the Critical Section for |XX| %d seconds |XX| for Message with |XX| Timestamp %d |XX| \n \n", node.id, numSeconds, msg.timestamp)
	time.Sleep(time.Duration(numSeconds) * time.Millisecond)
	fmt.Printf("\n|XX| Node %d |XX| is now exiting the Critical Section |XX|\n \n", node.id)

	node.nodeClock = combineClocks(node.id, node.nodeClock, msg.timestamp)

	if _, check := node.replyHistory[timestamp]; !check {
		node.replyHistory[timestamp] = node.createreplyHistory()
	}

	for id, req := range node.replyHistory[timestamp] {
		if id != node.id {
			if req {
				releaseMsg := instantiateMessage(Release, node.id, node.nodeClock)
				releaseMsg.status = DONE
				releaseMsg.replyTo.targetId = node.id
				releaseMsg.replyTo.timestamp = msg.timestamp
				node.send(*releaseMsg, id)
			}
		}
	}

	delete(node.replyHistory, timestamp)
	node.status[timestamp] = DONE

	if len(node.nodePQ) > 0 {
		head := node.nodePQ[0]
		node.vote = head
		go node.reply(head)
	}

}

func (node *Node) enterCSRequest() {

	node.nodeClock[node.id] += 1

	reqMsg := instantiateMessage(Request, node.id, node.nodeClock)

	if setnumNodes == 1 {
		node.enterCS(Message{
			node.nodeClock,
			node.id,
			Request,
			Receiver{},
			"",
		})

	}

	fmt.Printf("\n|XX| Node %d |XX| is requesting to enter Critical Section |XX| \n", node.id)

	time.Sleep(time.Duration(500) * time.Millisecond)

	otherNodes := map[int]bool{}

	timestamp := tostringTS(node.nodeClock)

	node.status[timestamp] = REQ

	for id, _ := range node.replyCheckList {

		if id == node.id {
			if node.vote == (Message{}) {
				otherNodes[id] = true
				node.vote = *reqMsg
			} else {
				node.enqueue(*reqMsg)
			}
		} else {
			otherNodes[id] = false
		}
	}

	node.replyHistory[timestamp] = otherNodes
	node.sendAll(*reqMsg)
}

func (node *Node) checkallreplies(timestamp [setnumNodes]int) bool {

	timestampStr := tostringTS(timestamp)

	total := 0

	if _, ok := node.replyHistory[timestampStr]; !ok {
		return false
	}

	for _, replyStatus := range node.replyHistory[timestampStr] {
		if replyStatus == true {
			total += 1
		}
	}

	if total >= setnumNodes/2+1 {
		return true
	}
	return false
}

func (node *Node) rescind(targetId int, timestamp [setnumNodes]int) {

	fmt.Printf("|XX| Node %d |XX| is Rescinding it's vote |XX|", node.id)

	node.nodeClock[node.id] += 1

	replyMsg := instantiateMessage(Rescind, node.id, node.nodeClock)

	replyMsg.replyTo.targetId = targetId
	replyMsg.replyTo.timestamp = timestamp

	node.send(*replyMsg, targetId)
}

func (node *Node) handleRequest(msg Message) {

	time.Sleep(time.Duration(500) * time.Millisecond)

	var earlierReq bool = false

	fmt.Printf("|XX| Node %d |XX| has it's current vote as |XX| %d at |XX| Time Stamp: %d |XX|\n", node.id, node.vote.senderId, node.vote.timestamp)

	for _, reqMsg := range node.nodePQ {

		if checkforEarlierTS(msg.timestamp, reqMsg.timestamp) {
			earlierReq = true
			break
		} else if checkifConcurrent(reqMsg.timestamp, msg.timestamp) && reqMsg.senderId < msg.senderId {
			earlierReq = true
			break
		}
	}

	if checkforEarlierTS(msg.timestamp, node.vote.timestamp) {
		earlierReq = true
	} else if checkifConcurrent(node.vote.timestamp, msg.timestamp) && node.vote.senderId < msg.senderId {
		earlierReq = true
	}

	if earlierReq {
		if node.vote != (Message{}) {
			if node.vote.senderId != node.id {
				node.rescind(node.vote.senderId, node.vote.timestamp)
			} else {
				node.replyHistory[tostringTS(node.vote.timestamp)][node.id] = false
			}
			fmt.Printf("|XX| Node %d rescinded it's earlier vote and is now voting to |XX| Node %d at |XX| Time Stamp: %d |XX| \n", node.id, node.vote.senderId, node.vote.timestamp)
			node.enqueue(Message{node.vote.timestamp, node.vote.senderId, Request, Receiver{}, REQ})

			node.vote = Message{}
		}
		go node.reply(msg)

	} else {
		node.enqueue(msg)
	}

	fmt.Printf("|XX| Node %d |XX| Queue Status: %s |XX| \n", node.id, tostringPQ(node.nodePQ))
}

func (node *Node) handleReply(msg Message) {

	msgTime := msg.replyTo.timestamp
	timestamp := tostringTS(msgTime)

	if _, check := node.replyHistory[timestamp]; !check {
		node.replyHistory[timestamp] = node.createreplyHistory()
	}

	if node.status[timestamp] == EXE || node.status[timestamp] == DONE {
		releaseMsg := instantiateMessage(Release, node.id, node.nodeClock)
		releaseMsg.status = DONE
		releaseMsg.replyTo.targetId = node.id
		releaseMsg.replyTo.timestamp = msg.timestamp
		node.send(*releaseMsg, msg.senderId)
		return
	}
	node.replyHistory[timestamp][msg.senderId] = true
	if node.vote == (Message{}) {
		node.vote = *instantiateMessage(Request, node.id, msgTime)
	}

	if node.checkallreplies(msg.replyTo.timestamp) {
		fmt.Printf("\n|XX| Node %d has received all replies for Request at Time Stamp: %d |XX|\n", node.id, msg.replyTo.timestamp)

		head := node.nodePQ[0]
		if head.timestamp == msgTime {
			if node.status[timestamp] == REQ {
				node.enterCS(head)
			}

		} else {
			if node.status[timestamp] == REQ {
				node.enterCS(Message{msg.replyTo.timestamp, node.id, Request, Receiver{}, REQ})

			}

		}
	}
}

func (node *Node) handleRelease(msg Message) {

	if msg.status == REQ {
		node.enqueue(node.vote)
	}

	node.vote = Message{}

	if len(node.nodePQ) > 0 {
		head := node.nodePQ[0]
		if head.senderId != node.id {
			go node.reply(head)
		} else {
			timestamp := tostringTS(head.timestamp)
			node.replyHistory[timestamp][node.id] = true
			if node.checkallreplies(head.timestamp) {
				fmt.Printf("[Node %d]  Request with TS: %d has recieved majority replies \n", node.id, msg.replyTo.timestamp)

				if node.status[timestamp] == REQ {
					node.enterCS(head)
				}

			}
		}
	}

}

func (node *Node) handleRescind(msg Message) {

	timestamp := tostringTS(msg.replyTo.timestamp)
	node.nodeClock = combineClocks(node.id, node.nodeClock, msg.timestamp)

	if node.status[timestamp] == EXE || node.status[timestamp] == DONE {
		return
	} else {
		node.replyHistory[timestamp][msg.senderId] = false
		releaseMsg := instantiateMessage(Release, node.id, node.nodeClock)
		releaseMsg.status = REQ
		releaseMsg.replyTo.targetId = node.id
		releaseMsg.replyTo.timestamp = msg.replyTo.timestamp

		go node.send(*releaseMsg, msg.senderId)
	}
}

func (node *Node) handleReceive(msg Message) {

	node.nodeClock = combineClocks(node.id, node.nodeClock, msg.timestamp)

	fmt.Printf("|XX| Node %d received a |XX| %s Message |XX| from Node %d |XX| \n", node.id, msg.msgSignature, msg.senderId)

	if msg.msgSignature == Request {
		node.handleRequest(msg)
	} else if msg.msgSignature == Vote {
		node.handleReply(msg)
	} else if msg.msgSignature == Release {
		node.handleRelease(msg)
	} else if msg.msgSignature == Rescind {
		node.handleRescind(msg)
	} else {
		fmt.Println("Invalid Message Signature")
	}

}

func (node *Node) listen() {

	for {
		select {
		case msg := <-node.nodeChannel:
			go node.handleReceive(msg)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	globalNodeMap := map[int]*Node{}

	for i := 0; i < setnumNodes; i++ {
		node := instantiateNode(i)
		globalNodeMap[i] = node
	}

	for i := 0; i < setnumNodes; i++ {
		wg.Add(1)
		globalNodeMap[i].replyCheckList = globalNodeMap
		globalNodeMap[i].nodeWG = &wg
	}

	for i := 0; i < setnumNodes; i++ {
		go globalNodeMap[i].listen()
	}

	start := time.Now()

	for i := 0; i < setnumNodes; i++ {
		go globalNodeMap[i].enterCSRequest()
	}

	wg.Wait()

	tEnd := time.Now()

	time.Sleep(time.Duration(3) * time.Second)

	fmt.Printf("\nTime Taken with %d nodes : %.2f seconds \n", setnumNodes, tEnd.Sub(start).Seconds())

	fmt.Printf("All Nodes have entered entered and exited the Critical Section \n")

	os.Exit(0)
}
