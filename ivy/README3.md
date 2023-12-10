## 50.041 Distributed Systems and Computing Programming Assignment 3

  

### Basic Folder Structure

In order to keep things simple the submission that implements IVY protocols with and without fault tolerance is split into 3 simple go files as follows:

1. ivybasicfunc.go
2. faulttolerantivy.go

  

You can simply run the simulations of all the 3 above mentioned protocols by a simple go run command as follows:

```

go run ivybasicfunc.go

```

```

go run faulttolerantivy.go

```



  

If you're having any issues please contact me on tele @jyotit_kaushal

 
PS: Please not I haven't added any excessive comments to the code itself as it's quite long and didn't want to make it harder to read/formulate.  Most of the functions have been defined quite aptly and relevantly for example all the "handle" functions are used to handle each type of message, sendall means send a certain message to all other nodes, etc. There are however some comments that have been added to ensure you can understand what each part in the code is doing and where each function has been written in case you want to refer to it at some point of testing. Please let me know if you have any difficulty with understanding any portion!

  

### Problem 1: Implementation of the basic IVY architecture
 In the first fie i.e. the `ivybasicfunc.go` as the name suggests we implement a simple basic version of the IVY protocol where there's a set up of different number of nodes which can easily be changed using the `NumberofProcessors` parameter/constant defined at the head of the file as well as the presence of a single Central Manager that interacts with all the nodes as defined by the IVY architecture.

Key structure of the codes are where we define the different structures of the Processor/Node class and the Central Manager class. Next we go on to define all the different types of messages that are going to be exchanged between a pair of Node and Central Manager as well as those among nodes. There have been different methods defined that guide the Central Manger/other processors by handling on whatever message they receive. 

For the implementation, we have again kept is simple for this basic simulation to show what kind of read/write type of operations happen in the IVY architecture. You can have a look at the main function of the ivybasicfunc.go to see that we first instantiate the central manager and the rest of the nodes, create a map of all the processors to pass to all of the nodes as well as the central manager.

We then simply start go routines for all the processors and central manager to start listening and then make all processors to do reads on pages and write on them twice. We then calculate the time taken for the whole simulation to run.

When you run the script you'll be able to see this time and as for rest of the outputs, I'll be explaining them together further down in this README.

### Problem 2. Implementation of basic IVY architecture that is now fault tolerant

 In the second file i.e. `faulttolerantivy.go`, the. backbone of the code remains the same. This includes things like initialisation of processors, central manager(s), defining methods to handle different type of methods that are exchanged between two entities. 
The difference however as the name suggests, comes from the fact that the IVY architecture defined in this file is now fault tolerant. This means that if the central manager that has been assigned with all the different processors falls into a fault, all operations aren't immediately halted but instead everything keeps working as planned.

This is done by the introduction of now a Backup Central Manager apart from the Primary Central Manager that is available to all processors. In this implementation, we setup all the nodes and the primary central manager as with the other one but now also set up a backup Primary Server.

What happens is that everytime the primary server falls, the processor on account of a timeout realize that their primary server (which in this implementation has an id of 0) has failed and in their own metadata which contains the node id for the Primary Central Manager (i.e. the one InCharge) as well as that of the Backup Central Manager. By default, Central Manager 0 is the primary server and the Central Manager 1 is the backup but if there is a fault where the Central Manager 0 dies, the processors realize this based on a timeout that's implemented and make a switch between the active cm and backup one.

The way this has been implemented is that in our case every 100ms the replicas share a metadata message amongst themselves and whenever one dies the other still has all the associated metadata with the simulation and then hence is able to act as the interim central manager in charge until the original primary central manager comes back up. I have created several experiments as were required in the problem set brief which we will get to below now that the rundown of both the implementations is complete.

### Problem 3: Reasoning as to whether this fault tolerant implementation of IVY maintains sequential consistency or not

Maintaining sequential consistency basically involves ensuring that there's a global order of operations and in the sense that all operations are done in some sort of total order that is consistent with what all the processors observe. A result of that is the condition that there shouldn't be/any possibility for conflicting operations (which in this case would be two processors trying to write to a document at the same time).

I believe my implementation of the fault tolerant ivy because of the way all different types of messages that form the inter-node communication are segregated into different channels. Since all the channels implemented in turn are not buffered, this leads to the eventuality that all messages that are received are processed in the same order that they have been received.

In addition to that, there's a system of acknowledgements which are handled in their own channels as defined above meaning they are also processed in order that they are received by the nodes and the central managers. The same logic holds true for the way the metadata is synched across different replicas in that it's all managed in a separate channel and since the metadata is updated at every short regular interval this ensures that when there is in fact a fault and there's a possibility where the order might be broken or is in danger, this constant updating ensures the metadata is kept sequentially updated along with the primary central manager.

Lastly, to handle the issues of concurrent write operations on a single document, the implementation uses a protocol in which the data is only written to the document by a processor if it has first received a confirmation from the current owner of the file and since the same owner can't give control over to two nodes at once, this means that there are no concurrent write operations on any of the documents happening and they only happen in the order that they are processed by the central manager which as established is sequential. 

### Problem 4: Experimentation of the different implementations


#### Understanding the output in `ivybasicfunc.go`

Now for each of the scenarios/experiments that have been listed below and are able to run, here are some basic points about the log messages that you see in the terminal. For each big message that is exchanged between the nodes/the central manager, there is a descriptive message on what is happening with all the relevant information given.

Below is an example of what you might see in the terminal:

```
|XX| Central Manager |XX| sending |XX| ReadNoOwner |XX| message to |XX| Processor 8 |XX| 
|XX| Processor 8 |XX| attempting to read |XX| Page 8 |XX| that doesn't have any owner yet |XX| 
|XX| Processor 8 |XX| sends |XX| ConfirmRead message |XX| to CM |XX| 
|XX| Processor 9 |XX| sends |XX| ReadRequest message |XX| to CM |XX| 
|XX| Central Manager |XX| receiving |XX| ConfirmRead message |XX| from |XX| Processor 8 |XX| 
|XX| Central Manager |XX| receiving |XX| ReadRequest message |XX| from |XX| Processor 9 |XX| 
|XX| Central Manager |XX| sending |XX| ReadNoOwner |XX| message to |XX| Processor 9 |XX| 
|XX| Processor 9 |XX| attempting to read |XX| Page 9 |XX| that doesn't have any owner yet |XX| 
|XX| Processor 10 |XX| sends |XX| ReadRequest message |XX| to CM |XX| 
```
  As you can see, it is quite descriptive on what's happening especially for the first basic ivy script there isn't much explanation that is required regarding the logs. For example, the first message is the Central Manager sending a ReadNoOwner message to Processor 8 which then attempts to read Page 8 which is what it was trying to read at the start, after reading it sends the CentralManager back a ReadConfirm Message, etc.

In all experiments in interest of time, in each simulation all processor try to read something twice and also write something twice.

Once this is done, the simulation then ends and you see the time taken for everything to finish.

#### Understanding the output from `faulttolerantivy.go`

The basic structure of all messages in the second implementation file follows the same structure as well, in terms of how the messages are displayed and what they are trying to say. There are however a few more things to understand

Firstly, when you run it you will be met with this terminal based interface. As shown, by entering in the desired number and pressing enter, you'll be running one of the 5 scenarios as defined in the problem set. Now, within the terminal you can see what each of them do and they all follow from the problem set so I won't go too deep into explaining what each of the scenarios corresponds to but you get the idea.

```
Trying out IVY with different Experiments as defined in PSET Description:

Experiment 1: Press 1 to run the basic scenario where there are no faults
Experiment 2: Press 2 to run the second scenario where the Primary CM dies and doesn't come back up
Experiment 3: Press 3 to run the third scenario where the Primary CM dies and then comes back up 
Experiment 4: Press 4 to run the fourth scenario where the Primary CM dies and then comes back up multiple times
Experiment 5: Press 5 to run the fifth scenario where both the Primary CM and the Backup CM die and then come back up 
```

Next thing to understand in this implementation of however, is that every 100ms the metadata between the two replicas is updated basis which Central Manager is currently InCharge. It would be really hard to see the rest of the messages going around if the log of each of these synch messages is displayed thereofre I have chosen not to show that but this is what's happening under the hood.

To showcase this further however, you do have the following extra prompts:

1. A central manager being killed and subsequently processors realizing that it's down
```
|XXXX| Primary Central Manager getting killed |XXXX| 

|XX| Processor 10 |XX| sending |XX| ConfirmWrite message |XX| to |XX| Central Manager 0 |XX| 
|XX| Central Manager 0 |XX| receiving |XX| ConfirmWrite message |XX| from |XX| Processor 10 |XX| 

|XXXX| Processor 10 |XX| realizes that Current Coordinator is Dead |XXXX| 


|XX| Processor 10 |XX| updates its pointers to Primary and Backup Central Managers |XX|
```
2. Every time central manager is killed/revived or the simulation is ending a printing of the metadata held by each of the replica as follows
```
|XXXX| Central Manager 1 Metadata |XXXX| 

|XXXX| Page: 1 |XX| Owner: 10 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 3 |XX| Owner: 2 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 4 |XX| Owner: 3 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 5 |XX| Owner: 4 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 6 |XX| Owner: 5 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 9 |XX| Owner: 8 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 2 |XX| Owner: 1 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 7 |XX| Owner: 6 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 8 |XX| Owner: 7 |XX| Access Type: ReadWriteAccess |XXXX| 
|XXXX| Page: 10 |XX| Owner: 9 |XX| Access Type: ReadWriteAccess |XXXX| 
```

Rest all remains the same, and the time is printed as before. In the next session, you see the results from all the scenarios in both the ivy without fault and all scenarios in ivy with fault tolerance. From the results you can see that there is in fact not too much of an overhead in exchanging messages between replicas by design (since the delay has been kept short since they are essentially meant to be together), etc. 

### Resulting Time Evaluation Table

  
| Scenario| Associated time taken by Simulation to run  |
|--|--|
| Ivy Without Fault Tolerance |  |
| Ivy With Fault Tolerance- Scenario 1|21.45 seconds|
| Ivy With Fault Tolerance- Scenario 2 | 21.28 seconds |
| Ivy With Fault Tolerance- Scenario 3|21.38 seconds|
| Ivy With Fault Tolerance- Scenario 4 |21.87 seconds  |
| Ivy With Fault Tolerance- Scenario 5|22.14 seconds|
