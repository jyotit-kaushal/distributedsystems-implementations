package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

type Message struct {
	Data string
}

type DistributedSystemSim struct {
	node_id        int
	coordinator_id int
	ids_ip         map[int]string
}

var distsystem = DistributedSystemSim{ // creating an instance of the main class of distributedsystem that contains id-ip mapping and what current ip is there
	node_id:        1,
	coordinator_id: 5,
	ids_ip: map[int]string{
		1: "127.0.0.1:5000",
		2: "127.0.0.1:5001",
		3: "127.0.0.1:5002",
		4: "127.0.0.1:5003",
		5: "127.0.0.1:5004",
	},
}

var no_election_invoked = true
var superiorNodeAvailable = false

func (distsystem *DistributedSystemSim) Election(invoker_id int, reply *Message) error { // function that receives and acknowledges an election from a lower node than itself
	fmt.Println("Getting Election Request from Node: ", invoker_id)

	if invoker_id < distsystem.node_id {

		fmt.Println("Sending Acknowledgement for Election to Node ", invoker_id)
		reply.Data = "OK"

		if no_election_invoked {
			no_election_invoked = false
			go invokeElection()
		}
	}
	return nil
}

func invokeElection() { // function to invoke election in all nodes more than itself
	for id, ip := range distsystem.ids_ip {

		reply := Message{""}

		if id > distsystem.node_id {

			fmt.Println("Sending an Election to Node: ", id)

			client, error := rpc.Dial("tcp", ip)

			if error != nil {
				fmt.Println("Node ", id, "is not available: Communication Failed")
				continue
			}

			err := client.Call("DistributedSystemSim.Election", distsystem.node_id, &reply)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if reply.Data == "OK" {
				fmt.Println("Received acknowledgement from Node", id)
				superiorNodeAvailable = true
			}
		}
	}

	if !superiorNodeAvailable {
		setCoordinator()
	}

	superiorNodeAvailable = false

	no_election_invoked = true
}

func setCoordinator() { // called when new coordinator has been decided

	reply := Message{""}

	for _, ip := range distsystem.ids_ip {

		client, error := rpc.Dial("tcp", ip)

		if error != nil {
			continue
		}
		time.Sleep(time.Duration(5) * time.Second)

		client.Call("DistributedSystemSim.SetCoordinatorHelper", distsystem.node_id, &reply)
	}
}

func (distsystem *DistributedSystemSim) SetCoordinatorHelper(id int, reply *Message) error {

	distsystem.coordinator_id = id

	fmt.Println("Node ", distsystem.coordinator_id, "has been elected as the new coordinator")
	return nil
}

func (distsystem *DistributedSystemSim) ReplyCheck(req_id int, reply *Message) error { // Checking if there has been an acknowledhement for whatever message was sent
	fmt.Println("Getting Request from Node: ", req_id)
	reply.Data = "OK"
	return nil
}

func SendRequesttoCoordinator() { // function to send request to coordinator- will invoke election if it doesn't receive an acknowledgement

	coord_id := distsystem.coordinator_id
	coord_ip := distsystem.ids_ip[coord_id]

	fmt.Println("Sending Request to Coordinator\n")

	node_id := distsystem.node_id
	reply := Message{""}

	client, error := rpc.Dial("tcp", coord_ip)

	if error != nil {
		fmt.Println("Not able to connect to coordinator -> Invoking Election to find new coordinator")
		invokeElection()
		return
	}

	error = client.Call("DistributedSystemSim.ReplyCheck", node_id, &reply)

	if error != nil || reply.Data != "OK" {
		fmt.Println("Not able to connect to coordinator -> Invoking Election to find new coordinator")
		invokeElection()
		return
	}

	fmt.Println("Request to Coordinator Successful", coord_id)
}

func main() {

	node_id := 0

	fmt.Printf("Enter the Node id [1-5]: ")
	fmt.Scanf("%d", &node_id)

	distsystem.node_id = node_id
	my_ip := distsystem.ids_ip[distsystem.node_id]

	ip, err := net.ResolveTCPAddr("tcp", my_ip) // Resolve the IP address for the RPC server

	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}

	rpc.Register(&distsystem)

	fmt.Println("This server is now on:", ip)
	go rpc.Accept(inbound)

	for {
		linecheck := ""
		fmt.Printf("Is this node reentering the system? Type in Yes and press senter if so, if not: \n")

		fmt.Printf("Press enter to communicate with coordinator.\n")

		fmt.Scanf("%s", &linecheck)

		if linecheck == "Yes" {
			invokeElection() // Invoke election if node is reentering the system
		} else {
			SendRequesttoCoordinator() // Send request to coordinator when hit Enter
		}

		fmt.Println("")
	}
}
