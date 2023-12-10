package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const NumberofNodes int = 10
const NumberofDocuments int = 10
const MaxSendDelay int = 300

type CentralManager struct {
	id               int
	cmStatus         cmStatusType
	processorMapping map[int]*Node
	wgCentralManager *sync.WaitGroup
	ownerMapping     map[int]int
	copysetMapping   map[int][]int
	reqMessage       chan Message
	resMessage       chan Message
	killChannel      chan int
	replicaChannel   chan synchMessage
}

type Node struct {
	id               int
	cm               *CentralManager
	backup           *CentralManager
	processorMapping map[int]*Node
	wgProcessor      *sync.WaitGroup
	accessMapping    map[int]AccessType
	contentMapping   map[int]string
	tbwContent       string
	reqMessage       chan Message
	resMessage       chan Message
	cmkillChannel    chan int
}

type Message struct {
	senderId    int
	requesterId int
	messageType MessageType
	page        int
	content     string
}

type synchMessage struct {
	senderId         int
	processorMapping map[int]*Node
	ownerMapping     map[int]int
	copysetMapping   map[int][]int
}

type MessageType int

const (
	ReadRequest MessageType = iota
	WriteRequest
	ConfirmRead
	ConfirmWrite
	ConfirmInvalidate
	ForwardRead
	ForwardWrite
	Invalidate
	ReadNoOwner
	WriteNoOwner
	ReadPageContent
	WritePageContent
)

type AccessType int

const (
	ReadOnlyAccess AccessType = iota
	ReadWriteAccess
)

type cmStatusType int

const (
	Coordinator cmStatusType = iota
	notInCharge
)

func (msg MessageType) String() string {
	var messagetypes = [...]string{"ReadRequest", "WriteRequest", "ConfirmRead", "ConfirmWrite", "ConfirmInvalidate", "ForwardRead", "ForwardWrite", "Invalidate", "ReadNoOwner", "WriteNoOwner", "ReadPageContent", "WritePageContent"}
	return messagetypes[msg]
}

func (msg AccessType) String() string {
	var accesstypes = [...]string{"ReadOnlyAccess", "ReadWriteAccess"}
	return accesstypes[msg]
}

func instantiateProcessor(id int, cm CentralManager, backup CentralManager) *Node {
	node := Node{
		id:               id,
		cm:               &cm,
		backup:           &backup,
		processorMapping: make(map[int]*Node),
		wgProcessor:      &sync.WaitGroup{},
		accessMapping:    make(map[int]AccessType),
		contentMapping:   make(map[int]string),
		tbwContent:       "",
		reqMessage:       make(chan Message),
		resMessage:       make(chan Message),
		cmkillChannel:    make(chan int),
	}

	return &node
}

func instantiateCMReplica(id int, cmStatus cmStatusType) *CentralManager {
	replica := CentralManager{
		id:               id,
		cmStatus:         cmStatus,
		processorMapping: make(map[int]*Node),
		wgCentralManager: &sync.WaitGroup{},
		ownerMapping:     make(map[int]int),
		copysetMapping:   make(map[int][]int),
		reqMessage:       make(chan Message),
		resMessage:       make(chan Message),
		killChannel:      make(chan int),
		replicaChannel:   make(chan synchMessage),
	}
	return &replica
}

func checkifinArray(e int, array []int) bool {
	for _, element := range array {
		if element == e {
			return true
		}
	}
	return false
}

func instatiateMessage(messageType MessageType, senderId int, requesterId int, page int, content string) *Message {
	msg := Message{
		messageType: messageType,
		senderId:    senderId,
		requesterId: requesterId,
		page:        page,
		content:     content,
	}

	return &msg
}

func (cm *CentralManager) showMetadata() {
	fmt.Printf("\n|XXXX| Central Manager %d Metadata |XXXX| \n\n", cm.id)

	for page, owner := range cm.ownerMapping {
		fmt.Printf("|XXXX| Page: %d |XX| Owner: %d |XX| Access Type: %s |XXXX| \n", page, owner, cm.processorMapping[owner].accessMapping[page])

	}
}

func (cm *CentralManager) Send(msg Message, recieverId int) {

	fmt.Printf("|XX| Central Manager %d |XX| sending |XX| %s |XX| message to |XX| Processor %d |XX| \n", cm.id, msg.messageType, recieverId)

	sendingDelay := rand.Intn(MaxSendDelay)
	time.Sleep(time.Millisecond * time.Duration(sendingDelay))

	recieverNode := cm.processorMapping[recieverId]

	if msg.messageType == ReadNoOwner || msg.messageType == WriteNoOwner {
		recieverNode.resMessage <- msg
	} else {
		recieverNode.reqMessage <- msg
	}
}

func (cm *CentralManager) sendMetadata(reciever *CentralManager) {

	synchMessage := synchMessage{
		senderId:         cm.id,
		processorMapping: cm.processorMapping,
		ownerMapping:     cm.ownerMapping,
		copysetMapping:   cm.copysetMapping,
	}

	reciever.replicaChannel <- synchMessage
}

func (cm *CentralManager) handleReadRequest(msg Message) {

	page := msg.page
	requesterId := msg.requesterId

	if _, check := cm.ownerMapping[page]; !check {
		replyMsg := instatiateMessage(ReadNoOwner, 0, requesterId, page, "")
		go cm.Send(*replyMsg, requesterId)
		confirmMsg := <-cm.resMessage

		fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, confirmMsg.messageType, confirmMsg.senderId)

		cm.wgCentralManager.Done()
		return
	}

	ownerMapping := cm.ownerMapping[page]
	pageCopyset := cm.copysetMapping[page]

	replyMsg := instatiateMessage(ForwardRead, 0, requesterId, page, "")

	if !checkifinArray(requesterId, pageCopyset) {
		pageCopyset = append(pageCopyset, requesterId)
	}

	go cm.Send(*replyMsg, ownerMapping)
	confirmMsg := <-cm.resMessage

	fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, confirmMsg.messageType, confirmMsg.senderId)

	cm.copysetMapping[page] = pageCopyset
	cm.wgCentralManager.Done()
}

func (cm *CentralManager) handleWriteRequest(msg Message) {

	page := msg.page
	requesterId := msg.requesterId

	if _, check := cm.ownerMapping[page]; !check {
		cm.ownerMapping[page] = requesterId
		replyMsg := instatiateMessage(WriteNoOwner, 0, requesterId, page, "")
		go cm.Send(*replyMsg, requesterId)
		confirmMsg := <-cm.resMessage
		fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, confirmMsg.messageType, confirmMsg.senderId)

		cm.wgCentralManager.Done()
		return
	}

	ownerMapping := cm.ownerMapping[page]
	pgCopySet := cm.copysetMapping[page]

	invalidationMsg := instatiateMessage(Invalidate, 0, requesterId, page, "")
	invalidationMsgCount := len(pgCopySet)

	for _, nodeid := range pgCopySet {
		go cm.Send(*invalidationMsg, nodeid)
	}

	for i := 0; i < invalidationMsgCount; i++ {
		msg := <-cm.resMessage
		fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, msg.messageType, msg.senderId)

	}

	responseMsg := instatiateMessage(ForwardWrite, 0, requesterId, page, "")
	go cm.Send(*responseMsg, ownerMapping)

	confirmMsg := <-cm.resMessage
	fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, confirmMsg.messageType, confirmMsg.senderId)

	cm.ownerMapping[page] = requesterId
	cm.copysetMapping[page] = []int{}
	cm.wgCentralManager.Done()
}

func (cm *CentralManager) handlesynchMessage(msg synchMessage) {

	cm.processorMapping = msg.processorMapping
	cm.copysetMapping = msg.copysetMapping
	cm.ownerMapping = msg.ownerMapping

}

func (cm *CentralManager) handleIncomingMessages() {

	for {
		select {
		case req := <-cm.reqMessage:
			cm.cmStatus = Coordinator
			fmt.Printf("|XX| Central Manager %d |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", cm.id, req.messageType, req.senderId)

			if req.messageType == ReadRequest {
				cm.handleReadRequest(req)
			} else if req.messageType == WriteRequest {
				cm.handleWriteRequest(req)

			}

		case synchMessage := <-cm.replicaChannel:
			cm.cmStatus = notInCharge

			cm.handlesynchMessage(synchMessage)

		case <-cm.killChannel:

			cm.cmStatus = notInCharge
			cm.showMetadata()

			return
		}
	}
}

func (node *Node) Send(msg Message, recieverId int) {

	if recieverId != 0 {
		fmt.Printf("|XX| Processor %d |XX| sending |XX| %s message |XX| to |XX| Processor %d |XX| \n", node.id, msg.messageType, recieverId)

	} else {
		fmt.Printf("|XX| Processor %d |XX| sending |XX| %s message |XX| to |XX| Central Manager %d |XX| \n", node.id, msg.messageType, node.cm.id)

	}

	sendingDelay := rand.Intn(MaxSendDelay)
	time.Sleep(time.Millisecond * time.Duration(sendingDelay))

	if recieverId == 0 {
		if msg.messageType == ReadRequest || msg.messageType == WriteRequest {
			node.cm.reqMessage <- msg
		} else if msg.messageType == ConfirmInvalidate || msg.messageType == ConfirmRead || msg.messageType == ConfirmWrite {
			node.cm.resMessage <- msg
		}
	} else {
		node.processorMapping[recieverId].resMessage <- msg
	}
}

func (node *Node) handleForwardRead(msg Message) {
	page := msg.page
	requesterId := msg.requesterId

	fmt.Printf("|XX| Processor %d |XX| handling |XX| %s message |XX| from CM |XX| \n", node.id, msg.messageType)
	fmt.Printf("|XX| Processor %d |XX| has |XX| %s |XX| for |XX| Page %d |XX| \n", node.id, node.accessMapping[page], page)

	if node.accessMapping[page] == ReadWriteAccess {
		node.accessMapping[page] = ReadOnlyAccess
	}

	responseMsg := instatiateMessage(ReadPageContent, node.id, requesterId, page, node.contentMapping[page])

	go node.Send(*responseMsg, requesterId)
}

func (node *Node) handleForwardWrite(msg Message) {

	page := msg.page
	requesterId := msg.requesterId

	responseMsg := instatiateMessage(WritePageContent, node.id, requesterId, page, node.contentMapping[page])

	delete(node.accessMapping, page)

	go node.Send(*responseMsg, requesterId)
}

func (node *Node) handleInvalidate(msg Message) {

	page := msg.page
	delete(node.accessMapping, page)

	responseMsg := instatiateMessage(ConfirmInvalidate, node.id, msg.requesterId, page, "")

	go node.Send(*responseMsg, 0)
}

func (node *Node) handleReadNoOwner(msg Message) {

	page := msg.page
	fmt.Printf("|XX| Processor %d |XX| attempting to read |XX| Page %d |XX| that doesn't have any owner yet |XX| \n", node.id, page)

	responseMsg := instatiateMessage(ConfirmRead, node.id, msg.requesterId, page, "")
	go node.Send(*responseMsg, 0)
}

func (node *Node) handleWriteNoOwner(msg Message) {
	page := msg.page

	node.contentMapping[page] = node.tbwContent
	node.accessMapping[page] = ReadWriteAccess

	responseMsg := instatiateMessage(ConfirmWrite, node.id, msg.requesterId, page, "")
	fmt.Printf("|XX| Processor %d |XX| writes to |XX| Page %d |XX| that doesn't have any owner yet |XX| with content |XX| %s |XX|  \n", node.id, page, node.tbwContent)

	go node.Send(*responseMsg, 0)
}

func (node *Node) handleReadPageContent(msg Message) {
	page := msg.page
	content := msg.content

	node.accessMapping[page] = ReadOnlyAccess
	node.contentMapping[page] = content

	fmt.Printf("|XX| Processor %d |XX| reading Page %d |XX| with content |XX| %s |XX| \n", node.id, page, content)

	responseMsg := instatiateMessage(ConfirmRead, node.id, msg.requesterId, page, "")
	go node.Send(*responseMsg, 0)
}

func (node *Node) handleWritePageContent(msg Message) {
	page := msg.page

	node.accessMapping[page] = ReadWriteAccess
	node.contentMapping[page] = node.tbwContent
	fmt.Printf("|XX| Processor %d |XX| writes to |XX| Page %d |XX| with content |XX| %s |XX|  \n", node.id, page, node.tbwContent)

	responseMsg := instatiateMessage(ConfirmWrite, node.id, msg.requesterId, page, "")
	go node.Send(*responseMsg, 0)
}

func (node *Node) handleIncomingMessage() {

	for {
		select {
		case msg := <-node.reqMessage:
			fmt.Printf("|XX| Processor %d |XX| receiving |XX| %s message |XX| from |XX| Central Manager %d |XX| \n", node.id, msg.messageType, node.cm.id)

			if msg.messageType == ForwardRead {
				node.handleForwardRead(msg)
			} else if msg.messageType == ForwardWrite {
				node.handleForwardWrite(msg)
			} else if msg.messageType == Invalidate {
				node.handleInvalidate(msg)
			}

		case <-node.cmkillChannel:
			fmt.Printf("\n|XXXX| Processor %d |XX| realizes that Current Coordinator is Dead |XXXX| \n\n", node.id)

			temp := node.backup
			node.backup = node.cm

			node.cm = temp

			fmt.Printf("\n|XX| Processor %d |XX| updates its pointers to Primary and Backup Central Managers |XX| \n\n", node.id)

		}

	}
}

func (node *Node) executeRead(page int) {

	node.wgProcessor.Add(1)

	if _, check := node.accessMapping[page]; check {

		content := node.contentMapping[page]

		fmt.Printf("> [Node %d] Reading Cached Page %d Content: %s\n", node.id, page, content)
		fmt.Printf("|XX| Processor %d |XX| reading from cached |XX| Page %d |XX| which has content |XX| %s|XX| \n", node.id, page, content)

		node.wgProcessor.Done()
		return
	}

	req := instatiateMessage(ReadRequest, node.id, node.id, page, "")
	go node.Send(*req, 0)

	msg := <-node.resMessage

	if msg.messageType == ReadNoOwner {
		node.handleReadNoOwner(msg)
	} else if msg.messageType == ReadPageContent {
		node.handleReadPageContent(msg)
	}
}

func (node *Node) WritetoPage(page int, content string) {

	node.wgProcessor.Add(1)

	if accessType, exists := node.accessMapping[page]; exists {

		if accessType == ReadWriteAccess && node.contentMapping[page] == content {

			fmt.Printf("|XX| Processor %d does |XX| Redundant Write |XX| to |XX| Page %d |XX| with content %s |XX| \n", node.id, page, content)
			node.wgProcessor.Done()
			return
		} else if accessType == ReadWriteAccess {
			node.tbwContent = content
			node.accessMapping[page] = ReadWriteAccess
			node.contentMapping[page] = node.tbwContent

			fmt.Printf("|XX| Processor %d does |XX| writes |XX| to |XX| Page %d |XX| with content %s |XX| \n", node.id, page, node.tbwContent)

			responseMsg := instatiateMessage(ConfirmWrite, node.id, node.id, page, "")
			go node.Send(*responseMsg, 0)
			return
		}
	}

	node.tbwContent = content

	WriteRequestMsg := instatiateMessage(WriteRequest, node.id, node.id, page, "")
	go node.Send(*WriteRequestMsg, 0)

	msg := <-node.resMessage

	if msg.messageType == WriteNoOwner {
		node.handleWriteNoOwner(msg)
	} else if msg.messageType == WritePageContent {
		node.handleWritePageContent(msg)
	}

}

func (cm *CentralManager) synchMetaData(reciever *CentralManager) {
	for {
		if cm.cmStatus == Coordinator {
			cm.sendMetadata(reciever)
		} else {
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func Default(nodeMap map[int]*Node, cm CentralManager, backupCM CentralManager, wg sync.WaitGroup) {

	start := time.Now()
	for i := 1; i <= NumberofNodes; i++ {
		nodeMap[i].executeRead(i)
	}
	for i := 1; i <= NumberofNodes; i++ {
		toWrite := fmt.Sprintf("This is written by node id %d", i)
		nodeMap[i].WritetoPage(i, toWrite)
	}
	for i := 1; i <= NumberofNodes; i++ {
		temp := i + 1
		temp %= (NumberofDocuments + 1)
		if temp == 0 {
			temp += 1
		}
		nodeMap[i].executeRead(temp)
	}
	for i := 1; i <= NumberofNodes; i++ {
		toWrite := fmt.Sprintf("This is written by pid %d", i)
		temp := i + 1
		temp %= (NumberofDocuments + 1)
		if temp == 0 {
			temp += 1
		}
		nodeMap[i].WritetoPage(temp, toWrite)
	}
	wg.Wait()
	cm.showMetadata()
	backupCM.showMetadata()
	end := time.Now()
	fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())
}

func main() {
	var wg sync.WaitGroup

	fmt.Printf("Testing with %d Processors.\n", NumberofNodes)

	fmt.Printf("\n\nTrying out IVY with different Experiments as defined in PSET Description:\n\n" +
		"Experiment 1: Press 1 to run the basic scenario where there are no faults\n" +
		"Experiment 2: Press 2 to run the second scenario where the Primary CM dies and doesn't come back up\n" +
		"Experiment 3: Press 3 to run the third scenario where the Primary CM dies and then comes back up \n" +
		"Experiment 4: Press 4 to run the fourth scenario where the Primary CM dies and then comes back up multiple times\n" +
		"Experiment 5: Press 5 to run the fifth scenario where both the Primary CM and the Backup CM die and then come back up \n")

	cm := instantiateCMReplica(0, Coordinator)
	cm.wgCentralManager = &wg

	backupCM := instantiateCMReplica(1, notInCharge)
	backupCM.wgCentralManager = &wg

	nodeMap := make(map[int]*Node)
	for i := 1; i <= NumberofNodes; i++ {
		node := instantiateProcessor(i, *cm, *backupCM)
		node.wgProcessor = &wg
		nodeMap[i] = node
	}

	cm.processorMapping = nodeMap

	for _, nodei := range nodeMap {
		for _, nodej := range nodeMap {
			if nodei.id != nodej.id {
				nodei.processorMapping[nodej.id] = nodej
			}
		}
	}

	go cm.handleIncomingMessages()
	go backupCM.handleIncomingMessages()

	for _, node := range nodeMap {
		go node.handleIncomingMessage()
	}

	go backupCM.synchMetaData(cm)
	go cm.synchMetaData(backupCM)

	time.Sleep(3 * time.Second)

	random := ""
	for {

		fmt.Scanf("%s", &random)

		if random == "1" {

			go Default(nodeMap, *cm, *backupCM, wg)

			time.Sleep(time.Duration(1) * time.Second)
			wg.Wait()
			os.Exit(0)
		}

		if random == "2" {
			go func() {

				start := time.Now()

				for i := 1; i <= NumberofNodes; i++ {
					nodeMap[i].executeRead(i)
				}
				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by node id %d", i)
					nodeMap[i].WritetoPage(i, toWrite)
				}

				fmt.Printf("\n|XXXX| Primary Central Manager getting killed |XXXX| \n\n")

				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}

				time.Sleep(200 * time.Millisecond)

				go cm.synchMetaData(backupCM)

				go cm.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].executeRead(temp)
				}
				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by pid %d", i)
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].WritetoPage(temp, toWrite)
				}

				wg.Wait()
				end := time.Now()
				time.Sleep(time.Duration(2) * time.Second)

				cm.showMetadata()
				backupCM.showMetadata()

				fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())
				os.Exit(0)
			}()
		}

		if random == "3" {
			go func() {

				start := time.Now()

				for i := 1; i <= NumberofNodes; i++ {
					nodeMap[i].executeRead(i)
				}
				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by node id %d", i)
					nodeMap[i].WritetoPage(i, toWrite)
				}
				fmt.Printf("\n|XXXX| Primary Central Manager getting killed |XXXX| \n\n")

				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)

				go cm.synchMetaData(backupCM)

				go cm.handleIncomingMessages()

				fmt.Printf("\n|XXXX| Primary Central Manager getting revived  |XXXX| \n\n")
				backupCM.killChannel <- 1
				cm.showMetadata()
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}

				time.Sleep(200 * time.Millisecond)
				go backupCM.synchMetaData(cm)
				go backupCM.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].executeRead(temp)
				}

				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by pid %d", i)
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].WritetoPage(temp, toWrite)
				}

				wg.Wait()
				end := time.Now()
				time.Sleep(time.Duration(1) * time.Second)

				cm.showMetadata()
				backupCM.showMetadata()
				fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())
				os.Exit(0)
			}()
		}

		if random == "4" {
			go func() {

				start := time.Now()

				for i := 1; i <= NumberofNodes; i++ {
					nodeMap[i].executeRead(i)
				}
				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by node id %d", i)
					nodeMap[i].WritetoPage(i, toWrite)
				}

				fmt.Printf("\n|XXXX| Primary Central Manager getting killed |XXXX| \n\n")

				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go cm.synchMetaData(backupCM)

				fmt.Printf("\n|XXXX| Primary Central Manager being revived  |XXXX| \n\n")

				go cm.handleIncomingMessages()

				fmt.Printf("\n|XXXX| Backup Central Manager being killed |XXXX| \n\n")

				backupCM.killChannel <- 1
				cm.showMetadata()
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}

				time.Sleep(200 * time.Millisecond)

				go backupCM.synchMetaData(cm)

				go backupCM.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].executeRead(temp)
				}

				fmt.Printf("\n|XXXX| Killing Primary Central Manager again while Backup Central Manager comes back up \n\n")

				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go cm.synchMetaData(backupCM)
				go cm.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by pid %d", i)
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].WritetoPage(temp, toWrite)
				}
				wg.Wait()
				end := time.Now()
				time.Sleep(time.Duration(1) * time.Second)
				cm.showMetadata()
				backupCM.showMetadata()
				fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())
				os.Exit(0)
			}()

		}

		if random == "5" {
			go func() {

				start := time.Now()
				for i := 1; i <= NumberofNodes; i++ {
					nodeMap[i].executeRead(i)
				}
				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by node id %d", i)
					nodeMap[i].WritetoPage(i, toWrite)
				}

				fmt.Printf("\n|XXXX| Primary Central Manager getting killed |XXXX| \n\n")
				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go cm.synchMetaData(backupCM)
				go cm.handleIncomingMessages()

				fmt.Printf("\n|XXXX| Primary Central Manager getting revived |XXXX| \n\n")

				fmt.Printf("\n|XXXX| Backup Central Manager getting killed |XXXX| \n\n")

				backupCM.killChannel <- 1
				cm.showMetadata()
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go backupCM.synchMetaData(cm)
				go backupCM.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].executeRead(temp)
				}

				fmt.Printf("\n|XXXX| Killing Primary Central Manager again while Backup Central Manager comes back up \n\n")
				cm.killChannel <- 1
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go cm.synchMetaData(backupCM)
				go cm.handleIncomingMessages()

				fmt.Printf("\n|XXXX| Killing Backup Central Manager again while Primary Central Manager comes back up \n\n")
				backupCM.killChannel <- 1
				cm.showMetadata()
				for _, node := range nodeMap {
					node.cmkillChannel <- 1
				}
				time.Sleep(100 * time.Millisecond)
				go backupCM.synchMetaData(cm)
				go backupCM.handleIncomingMessages()

				for i := 1; i <= NumberofNodes; i++ {
					toWrite := fmt.Sprintf("This is written by pid %d", i)
					temp := i + 1
					temp %= (NumberofDocuments + 1)
					if temp == 0 {
						temp += 1
					}
					nodeMap[i].WritetoPage(temp, toWrite)
				}
				wg.Wait()
				end := time.Now()
				time.Sleep(time.Duration(1) * time.Second)
				cm.showMetadata()
				backupCM.showMetadata()
				fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())
				os.Exit(0)
			}()

		}
	}

}
