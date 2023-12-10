package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NumberofProcessors int = 10
	NumberofPages      int = 20
	MaxSendDelay       int = 300
)

type CentralManager struct {
	processors             map[int]*Processor
	waitGroup              *sync.WaitGroup
	ownerMapping           map[int]int
	copysetMapping         map[int][]int
	incomingreqChannel     chan Message
	outgoingconfirmChannel chan Message
}

type Processor struct {
	pid              int
	cm               *CentralManager
	processorMapping map[int]*Processor
	waitGroup        *sync.WaitGroup
	accessMapping    map[int]AccessType
	contentMapping   map[int]string
	contentToWrite   string
	generalChannel   chan Message
	responseChannel  chan Message
}

type Message struct {
	messageType MessageType
	senderID    int
	requesterID int
	page        int
	content     string
}

type MessageType int
type AccessType int

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

const (
	ReadOnlyAccess AccessType = iota
	ReadWriteAccess
)

func (msg MessageType) String() string {
	var messagetypes = [...]string{"ReadRequest", "WriteRequest", "ConfirmRead", "ConfirmWrite", "ConfirmInvalidate", "ForwardRead", "ForwardWrite", "Invalidate", "ReadNoOwner", "WriteNoOwner", "ReadPageContent", "WritePageContent"}
	return messagetypes[msg]
}

func (msg AccessType) String() string {
	var accesstypes = [...]string{"ReadOnlyAccess", "ReadWriteAccess"}
	return accesstypes[msg]
}

func checkifinArray(e int, array []int) bool {
	for _, element := range array {
		if element == e {
			return true
		}
	}
	return false
}

func (cm *CentralManager) Send(msg Message, recipientID int) {

	fmt.Printf("|XX| Central Manager |XX| sending |XX| %s |XX| message to |XX| Processor %d |XX| \n", msg.messageType, recipientID)

	sendingDelay := rand.Intn(MaxSendDelay)
	time.Sleep(time.Millisecond * time.Duration(sendingDelay))

	receiver := cm.processors[recipientID]

	if msg.messageType == ReadNoOwner || msg.messageType == WriteNoOwner {
		receiver.responseChannel <- msg
	} else {
		receiver.generalChannel <- msg
	}
}

func (cm *CentralManager) handleReadRequest(req Message) {

	page := req.page
	requesterID := req.requesterID

	if _, check := cm.ownerMapping[page]; !check {

		noOwnerMsg := Message{
			messageType: ReadNoOwner,
			senderID:    0,
			requesterID: requesterID,
			page:        page,
			content:     "",
		}

		go cm.Send(noOwnerMsg, requesterID)

		confirmMsg := <-cm.outgoingconfirmChannel

		fmt.Printf("|XX| Central Manager |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", confirmMsg.messageType, confirmMsg.senderID)

		cm.waitGroup.Done()
		return
	}

	currentOwner := cm.ownerMapping[page]
	copyset := cm.copysetMapping[page]
	forwardMsg := Message{
		messageType: ForwardRead,
		senderID:    0,
		requesterID: requesterID,
		page:        page,
		content:     "",
	}

	go cm.Send(forwardMsg, currentOwner)

	if !checkifinArray(requesterID, copyset) {
		copyset = append(copyset, requesterID)
	}

	cm.copysetMapping[page] = copyset

	confirmMsg := <-cm.outgoingconfirmChannel

	fmt.Printf("|XX| Central Manager |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", confirmMsg.messageType, confirmMsg.senderID)

	cm.waitGroup.Done()
}

func (cm *CentralManager) handleWriteRequest(req Message) {

	page := req.page
	requesterID := req.requesterID

	if _, check := cm.ownerMapping[page]; !check {

		cm.ownerMapping[page] = requesterID
		noOwnerMsg := Message{
			messageType: WriteNoOwner,
			senderID:    0,
			requesterID: requesterID,
			page:        page,
			content:     "",
		}

		go cm.Send(noOwnerMsg, requesterID)

		confirmMsg := <-cm.outgoingconfirmChannel

		fmt.Printf("|XX| Central Manager |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", confirmMsg.messageType, confirmMsg.senderID)

		cm.waitGroup.Done()
		return
	}

	currentOwner := cm.ownerMapping[page]
	copyset := cm.copysetMapping[page]

	InvalidateMsg := Message{
		messageType: Invalidate,
		senderID:    0,
		requesterID: requesterID,
		page:        page,
		content:     "",
	}

	numToInvalidate := len(copyset)

	for _, pid := range copyset {
		go cm.Send(InvalidateMsg, pid)
	}

	for i := 0; i < numToInvalidate; i++ {
		<-cm.outgoingconfirmChannel
	}

	fwd := Message{
		messageType: ForwardWrite,
		senderID:    0,
		requesterID: requesterID,
		page:        page,
		content:     "",
	}
	go cm.Send(fwd, currentOwner)

	confirmMsg := <-cm.outgoingconfirmChannel

	fmt.Printf("|XX| Central Manager |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", confirmMsg.messageType, confirmMsg.senderID)
	cm.ownerMapping[page] = requesterID
	cm.copysetMapping[page] = []int{}
	cm.waitGroup.Done()
}

func (cm *CentralManager) listen() {

	for {
		req := <-cm.incomingreqChannel

		fmt.Printf("|XX| Central Manager |XX| receiving |XX| %s message |XX| from |XX| Processor %d |XX| \n", req.messageType, req.senderID)

		if req.messageType == ReadRequest {
			cm.handleReadRequest(req)
		} else if req.messageType == WriteRequest {
			cm.handleWriteRequest(req)
		}

	}
}

func (p *Processor) Send(msg Message, recipientID int) {

	if recipientID == 0 {
		fmt.Printf("|XX| Processor %d |XX| sends |XX| %s message |XX| to CM |XX| \n", p.pid, msg.messageType)
	} else {
		fmt.Printf("|XX| Processor %d |XX| sends |XX| %s message |XX| to Processor %d |XX| \n", p.pid, msg.messageType, recipientID)
	}

	sendDelay := rand.Intn(MaxSendDelay)
	time.Sleep(time.Millisecond * time.Duration(sendDelay))

	if recipientID == 0 {

		if msg.messageType == ReadRequest || msg.messageType == WriteRequest {
			p.cm.incomingreqChannel <- msg
		} else if msg.messageType == ConfirmRead || msg.messageType == ConfirmWrite || msg.messageType == ConfirmInvalidate {
			p.cm.outgoingconfirmChannel <- msg
		}
	} else {
		p.processorMapping[recipientID].responseChannel <- msg
	}
}

func (p *Processor) handleForwardRead(forwardMsg Message) {

	page := forwardMsg.page
	requester := forwardMsg.requesterID

	contentMsg := Message{
		messageType: ReadPageContent,
		senderID:    p.pid,
		requesterID: requester,
		page:        page,
		content:     p.contentMapping[page],
	}

	go p.Send(contentMsg, requester)
}

func (p *Processor) handleForwardWrite(forwardMsg Message) {

	page := forwardMsg.page
	requester := forwardMsg.requesterID

	contentMsg := Message{
		messageType: WritePageContent,
		senderID:    p.pid,
		requesterID: requester,
		page:        page,
		content:     p.contentMapping[page],
	}

	go p.Send(contentMsg, requester)

	delete(p.accessMapping, page)
}

func (p *Processor) handleInvalidate(InvalidateMsg Message) {

	page := InvalidateMsg.page
	delete(p.accessMapping, page)

	confirmMsg := Message{
		messageType: ConfirmInvalidate,
		senderID:    p.pid,
		requesterID: InvalidateMsg.requesterID,
		page:        page,
		content:     "",
	}

	go p.Send(confirmMsg, 0)
}

func (p *Processor) handleReadNoOwner(noOwnerMsg Message) {

	page := noOwnerMsg.page

	fmt.Printf("|XX| Processor %d |XX| attempting to read |XX| Page %d |XX| that doesn't have any owner yet |XX| \n", p.pid, page)

	confirmMsg := Message{
		messageType: ConfirmRead,
		senderID:    p.pid,
		requesterID: noOwnerMsg.requesterID,
		page:        page,
		content:     "",
	}
	go p.Send(confirmMsg, 0)
}

func (p *Processor) handleWriteNoOwner(noOwnerMsg Message) {

	page := noOwnerMsg.page
	p.contentMapping[page] = p.contentToWrite
	p.accessMapping[page] = ReadWriteAccess

	fmt.Printf("|XX| Processor %d |XX| writes to |XX| Page %d |XX| that doesn't have any owner yet |XX| with content |XX| %s |XX|  \n", p.pid, page, p.contentToWrite)

	confirmMsg := Message{
		messageType: ConfirmWrite,
		senderID:    p.pid,
		requesterID: noOwnerMsg.requesterID,
		page:        page,
		content:     "",
	}
	go p.Send(confirmMsg, 0)
}

func (p *Processor) handleReadPageContent(contentMsg Message) {

	page := contentMsg.page
	content := contentMsg.content
	p.contentMapping[page] = content
	p.accessMapping[page] = ReadOnlyAccess

	fmt.Printf("|XX| Processor %d |XX| reading Page %d |XX| with content |XX| %s |XX| \n", p.pid, page, content)

	confirmMsg := Message{
		messageType: ConfirmRead,
		senderID:    p.pid,
		requesterID: contentMsg.requesterID,
		page:        page,
		content:     "",
	}

	go p.Send(confirmMsg, 0)
}

func (p *Processor) handleWritePageContent(contentMsg Message) {

	page := contentMsg.page
	p.contentMapping[page] = p.contentToWrite
	p.accessMapping[page] = ReadWriteAccess

	fmt.Printf("|XX| Processor %d |XX| writes to |XX| Page %d |XX| with content |XX| %s |XX|  \n", p.pid, page, p.contentToWrite)

	confirmMsg := Message{
		messageType: ConfirmWrite,
		senderID:    p.pid,
		requesterID: contentMsg.requesterID,
		page:        page,
		content:     "",
	}
	go p.Send(confirmMsg, 0)

}

func (p *Processor) listen() {

	for {
		msg := <-p.generalChannel
		fmt.Printf("|XX| Processor %d |XX| handling |XX| %s message |XX| from CM |XX| \n", p.pid, msg.messageType)
		if msg.messageType == ForwardRead {
			p.handleForwardRead(msg)
		} else if msg.messageType == ForwardWrite {
			p.handleForwardWrite(msg)
		} else if msg.messageType == Invalidate {
			p.handleInvalidate(msg)
		}
	}

}

func (p *Processor) ReadPage(page int) {

	p.waitGroup.Add(1)

	if _, check := p.accessMapping[page]; check {

		content := p.contentMapping[page]

		fmt.Printf("|XX| Processor %d |XX| reading from cached |XX| Page %d |XX| which has content |XX| {%s} |XX| \n", p.pid, page, content)

		p.waitGroup.Done()
		return
	}

	req := Message{
		messageType: ReadRequest,
		senderID:    p.pid,
		requesterID: p.pid,
		page:        page,
		content:     "",
	}

	go p.Send(req, 0)

	msg := <-p.responseChannel

	if msg.messageType == ReadNoOwner {
		p.handleReadNoOwner(msg)
	} else if msg.messageType == ReadPageContent {
		p.handleReadPageContent(msg)
	}

}

func (p *Processor) WritetoPage(page int, content string) {

	p.waitGroup.Add(1)

	if accessType, check := p.accessMapping[page]; check {

		if accessType == ReadWriteAccess && p.contentMapping[page] == content {
			fmt.Printf("|XX|Processor %d does |XX| Redundant Write |XX| to |XX| Page %d |XX| with content %s |XX| \n", p.pid, page, content)

			p.waitGroup.Done()
			return
		}
	}

	p.contentToWrite = content
	requestMsg := Message{
		messageType: WriteRequest,
		senderID:    p.pid,
		requesterID: p.pid,
		page:        page,
		content:     content,
	}

	go p.Send(requestMsg, 0)
	msg := <-p.responseChannel

	if msg.messageType == WriteNoOwner {
		p.handleWriteNoOwner(msg)
	} else if msg.messageType == WritePageContent {
		p.handleWritePageContent(msg)
	}

}

func main() {

	var wg sync.WaitGroup

	cm := &CentralManager{
		processors:             make(map[int]*Processor),
		waitGroup:              &wg,
		ownerMapping:           make(map[int]int),
		copysetMapping:         make(map[int][]int),
		incomingreqChannel:     make(chan Message),
		outgoingconfirmChannel: make(chan Message),
	}

	allProcessors := make(map[int]*Processor)
	for i := 1; i <= NumberofProcessors; i++ {
		p := &Processor{
			pid:              i,
			cm:               cm,
			processorMapping: make(map[int]*Processor),
			waitGroup:        &wg,
			accessMapping:    make(map[int]AccessType),
			contentMapping:   make(map[int]string),
			contentToWrite:   "",
			generalChannel:   make(chan Message),
			responseChannel:  make(chan Message),
		}
		allProcessors[p.pid] = p
	}
	cm.processors = allProcessors

	for _, p1 := range allProcessors {
		for _, p2 := range allProcessors {
			if p1.pid != p2.pid {
				p1.processorMapping[p2.pid] = p2
			}
		}
	}

	go cm.listen()
	for _, p := range allProcessors {
		go p.listen()
	}

	start := time.Now()

	// initial Read
	for i := 1; i <= NumberofProcessors; i++ {
		allProcessors[i].ReadPage(i)
	}

	for i := 1; i <= NumberofProcessors; i++ {
		content := fmt.Sprintf("Write content by Processor %d", i)
		allProcessors[i].WritetoPage(i, content)
	}

	for i := 1; i <= NumberofProcessors; i++ {
		temp := i + 1
		temp %= (NumberofPages + 1)
		if temp == 0 {
			temp += 1
		}
		allProcessors[i].ReadPage(temp)
	}

	for i := 1; i <= NumberofProcessors; i++ {
		toWrite := fmt.Sprintf("This is written by pid %d", i)
		temp := i + 1
		temp %= (NumberofPages + 1)
		if temp == 0 {
			temp += 1
		}
		allProcessors[i].WritetoPage(temp, toWrite)
	}

	wg.Wait()
	end := time.Now()
	time.Sleep(time.Second * 3)
	fmt.Printf("Time taken = %.2f seconds \n", end.Sub(start).Seconds())

}
