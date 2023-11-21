## 50.041 Distributed Systems and Computing Programming Assignment 2

### Basic Folder Structure
In order to keep things simple the entirety of the submission that implements these DME protocols is split into 3 simple go files as follows:
1. lamportspqprotocol.go
2. raspqprotocol.go
3. votingprotocol.go

You can simply run the simulations of all the 3 above mentioned protocols by a simple go run command as follows:
```
go run lamportspqprotocol.go
```
```
go run raspqprotocol.go
```
```
go run votingprotocol.go
```

If you're having any issues please contact me on tele @jyotit_kaushal

PS: Please not I haven't added any excessive comments to the code itself as it's quite long and didn't want to make it harder to read/formulate. Most of the functions have been defined quite aptly and relevantly for example all the "handle" functions are used to handle each type of message, sendall  means send a certain message to all other nodes, etc. Please let me know if you have any difficulty with understanding any portion!

### Protocol 1: Lamport's Shared Priority Queue Protocol

In this implementation, we are essentially replicating the shared Priority Queue Protocol that distributed systems use in order to manage critical sections or shared memory which is quite common in practical scenarios. To implement this in Go, we make use of goroutines, waitgroups and implemented queues that help us simulate multiple nodes trying to access a critical section.

In paritcular, we also want to test out the performance of different protocols and how it scales with an increase in number of nodes and that's why we will be running the protocol from 1 through 10 nodes and noting the time taken by each until each node has entered and exited the critical section without any compromises.
To do this, we make use of the constant ```setnumofNodes``` defined at the  head od the code with all the other constants and you can make your own change to it and observer the results as well.

To keep things a bit more interesting, we also have that each node is entering the critical section for a time-frame of "NodeID *mod* 3" seconds. That is Node 9 would enter it for 0 seconds whereas Node 7 would enter it for 1 second and so on.
Now , once you run the above command you will see the output of the code. Below I help you understand what some of the printed statements in the terminal actually mean and how they relate to the protocol.

### Protocol 2: Ricart and Agrawal optimized Shared Priority Queue Protocol

Please note than in this optimized version of Lamport's SPQ Protocol, most of the code behind it is still pretty much the same other than the fact how the shared priority queues are updated and what's the final check done before a node enters the critical section.

The method of testing for nodes from 1 through 10 is again still the same and most of the print statements are also basically the same, and as a result I'll be covering the output of both these protocols together to avoid redundancy.

Following are the types of print statements you will see upon running the simulations and this is how to interpret them:

#### 1. Clients requesting to enter Critical Section
Outputs in this format are simply the nodes requesting to enter the critical section, the order of which has been randomized within the code whis is why you see this haphazard order. All the nodes try to enter the critical section at start of the simulation and this is where the protocol starts up.

```
|XX| Node 10 |XX| is requesting to enter Critical Section |XX| 

|XX| Node 3 |XX| is requesting to enter Critical Section |XX| 

|XX| Node 1 |XX| is requesting to enter Critical Section |XX| 

|XX| Node 2 |XX| is requesting to enter Critical Section |XX| 
```

This is followed up by each node subsequently requesting each of the other nodes individually after putting this request in their own queue, the result of which is presented like as follows. Upon receiving this request, all nodes decide what to do by the logic defined in the  ```handleRequest()``` function.

```
|XX| Node 5 sending a |XX| Request message |XX| to Node 8 |XX| at MemAddr 0x1400011a230 |XX|
|XX| Node 8 sending a |XX| Request message |XX| to Node 1 |XX| at MemAddr 0x1400011a000 |XX|
```

#### 2. Reading the queue status
Each time a change has been made to any of the queues of the node (defined as nodePQ) in the type Node, it's converted to a string and then printed to keep track of what state of the protocol/simulation we are currently in. You will see it as follows:
```
|XX| Node 7 |XX| Queue Status: |X| Priority Queue: Request at [TS: 1] by Node 7 |X| Request at [TS: 1] by Node 10 |X| Request at [TS: 1] by Node 9 |X| Request at [TS: 1] by Node 8 |X|  |XX| 
```

#### 3. Nodes receiving reply message from  other nodes
A bit self explanatory, but once the nodes receive a request and they're okay to reply to the node making the request, they'll send the request and print this message.

```
|XX| Node 7 |XX| sending reply to |XX| Node 1 |XX|
|XX| Node 7 sending a |XX| Reply message |XX| to Node 1 |XX| at MemAddr 0x1400011a000 |XX|
```

#### 4. Nodes entering and exiting the critical section
Once a node has received all needed replies from the other nodes and checks if it's at the head of the SPQ, it enters the critical section and thereby exists it. You can see this in the terminal output as follows:

```
|XX| Node 1 has received all replies for Request at Time Stamp: 1 |XX|

|XX| Node 1 |XX| is entering the Critical Section for |XX| 1 seconds |XX| for Message with |XX| Timestamp 1 |XX| 
 
|XX| Node 6 received a |XX| Reply Message |XX| from Node 8 |XX| 
|XX| Node 2 received a |XX| Reply Message |XX| from Node 6 |XX| 
|XX| Node 7 received a |XX| Reply Message |XX| from Node 10 |XX| 
|XX| Node 5 received a |XX| Reply Message |XX| from Node 6 |XX| 

|XX| Node 1 |XX| is now exiting the Critical Section |XX|
```
#### 5. End of the simulation
This will continue with all other nodes until all nodes have entered and then exited out of the Critical Section. Once this is done, you can see the final few lines of the printed output indicating the time taken for the simulation to run as well as the set number of nodes for reference.

```
Time Taken with 10 nodes : 39.80 seconds 
Simulation Ending 
```

### Protocol 3: Voting Protocol with Deadlock Prevention
Most of the print statements for the implementation of the voting protocol are also the same  just that one point to note is that while the other two protocols were using logical clocks for all their operations, the implementation of the voting protocol uses vector clocks instead. That's why you see the print of the shared priroity queue in this instance be of the form:

```
 Request at [TS: [0 0 0 0 0 1 0 0 0 0]] by Node 5 |X| Request at [TS: [0 0 0 0 1 0 0 0 0 0]] by Node 4 |X| Request at [TS: [0 0 0 1 0 0 0 0 0 0]] by Node 3 |X| Request at [TS: [0 0 1 0 0 0 0 0 0 0]] by Node 2 |X| Request at [TS: [0 1 0 0 0 0 0 0 0 0]] by Node 1 |X| Request at [TS: [1 0 0 0 0 0 0 0 0 0]] by Node 0 |X|  |XX| 
```

Rest all the print statements are basically the same as compared to the ones shown above (from 1. to 5.)

### Resulting Time Evaluation Table

| Number of Nodes| Time Taken By Lamport's SPQ Protocol | Time Taken by RA SPQ Protocol| Time Taken by Voting Protocol|
| ------------- | ----------- |----------|-----------|
| 1 Node    |   1.00 seconds     |    1.00 second    |    0.00 seconds (since Nodes start from 0 in this protocol)     |
| 2 Nodes   | 11.10 seconds       |  12.16 seconds       |     14.24 seconds    |
| 3 Nodes   |   14.65 seconds     |   13.45 seconds      |    18.94 seconds     |
| 4 Nodes   |    17.76 seconds    |     17.87 seconds    |    23.28 seconds     |
| 5 Nodes   |     23.05 seconds   |     21.85 seconds    |    25.97 seconds     |
| 6 Nodes   |    24.74 seconds    |   24.84 seconds      |    28.77 seconds     |
| 7 Nodes   |   28.32 seconds     |    28.95 seconds     |  27.63 seconds       |
| 8 Nodes   |   33.12 seconds     |   33.62 seconds      |     28.74 seconds    |
| 9 Nodes   |  36.27 seconds      |     35.46 seconds    |    27.58 seconds     |
| 10 Nodes   |    40.53 seconds    |   38.80 seconds      |   29.66 seconds      |

Based off of that we see that while both Lamport's SPQ Protocol and Ricart Agrawal's Optimzied SPQ Protocol both see an increase in time for all nodes to enter critical section, for voting protocol as the nodes increase there is more or less bit of a plateau in the running times.

One more thing to note is that while Ricart Agrwal's Protocol is a bit faster than Lamport's one when the number of nodes is 9 or 10, this difference falls apart and even becomes negatie actually when moving on to lower number of nodes. This is quite interesting!



