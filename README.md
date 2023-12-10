## 50.041 Distributed Systems and Computing Programming Assignment 1

### Basic Folder Structure
In order to keep things simple the entirety of the submission that implements lamport's logical clocks using both logical clocks as well as vector clocks; and a simulation of Bully Election Algorithm is split into 4 files:
1. logicalclock.go
2. vectorclock.go
3. bullyalgorithm.go
4. bullyalgoexe.bat

You can simply run the simulations of Logical Clock, Vector Clock, and Bully Algorithm by a simple go run command as follows:
```
go run logicalclock.go
```
```
go run vectorclock.go
```
```
go run bullyalgorithm.go
```

If you're having any issues please contact me on tele @jyotit_kaushal

### Prob 1.1- Implementing Lamport's Logical Clock for Total Ordering of Events

In this implementation of Lamport's Logical Clocks, we are taking essentially a simulate

#### 1. Clients Being Registered
You can see the clients being registered at the start of the simulation. The total number of clients can be set to whatever number you want by changing the ```tnoc``` constant at the top. By default, it is set to 10 clients
```
Registering Client 1 
Registering Client 2
Registering Client 3
Registering Client 4
Registering Client 5
Registering Client 6
Registering Client 7
Registering Client 8
Registering Client 9
Registering Client 10
```

#### 2. Clients Sending Messages --> Server Receiving Messages from Clients
Clients send a set number of messages before shutting down. This number can also be changed based on your preference by again setting the value of the ```tnmtbsbc``` constant. This is how a message sent by the client which is received by the server looks like. Each message is a "Hello" followed by what number message it is from that particular client. The last message is simply a "Last Hello" indicating this is the last message that is going to be sent by this client.
```
|XX| Logical Clock: 3 |XX| Received By: Server   |XX| Message: Hello 1 from Client 3 |XX| Sent By: Client 3 |XX|
```
#### 3. Result of CoinToss by Server --> Messages being Broadcasted
Upon receiving the message, the server does a coin toss and the result is displayed whether it is going to broadcast the message or drop the message. If the coin toss decides a broadcasting of the message that happens and you see a Broadasted message sent by the server and recieved by the client as follows:
```
|XXXX| Server : Broadcasting 'Hello 1 from Client 1' from Client 1 |XXXX|

|XX| Logical Clock: 7 |XX| Received By: Client 6 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Sent By: Server |XX|
|XX| Logical Clock: 7 |XX| Received By: Client 8 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Sent By: Server |XX|
|XX| Logical Clock: 7 |XX| Received By: Client 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Sent By: Server |XX|
```
Similarly if the CoinToss decides dropping of the message you can catch this message in the terminal.
```
|XXXX| Server : Not Broadcasting 'Hello 1 from Client 7' from Client 7 |XXXX|
```
#### 4. Total Order of Events
Once all the messages has been sent and all clients have been shut down (being tracked as a count of active clients in the count), all messages that were being stored in a parent channel as the simulation was going on are pulled out of the channel and put in an array which is then sorted using the logical clock values of the messages and displayed in order as follows. A few of the messages are shown below, but you can run the code as indicated above to see the entire order by yourself.
```
|XXXX| Total Order of All Messages Below: |XXXX|

|XX| 0. |XX| Logical Clock: 2 |XX| Message: Hello 1 from Client 9 |XX| Received By: Server |XX| Sent By: Client 9 |XX|
|XX| 1. |XX| Logical Clock: 3 |XX| Message: Hello 1 from Client 3 |XX| Received By: Server |XX| Sent By: Client 3 |XX|
|XX| 2. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 10 |XX| Sent By: Server |XX|
|XX| 3. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 4 |XX| Sent By: Server |XX|
|XX| 4. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 7 |XX| Sent By: Server |XX|
|XX| 5. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 8 |XX| Sent By: Server |XX|
|XX| 6. |XX| Logical Clock: 5 |XX| Message: Hello 1 from Client 1 |XX| Received By: Server |XX| Sent By: Client 1 |XX|
|XX| 7. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 9 |XX| Sent By: Server |XX|
|XX| 8. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 1 |XX| Sent By: Server |XX|
|XX| 9. |XX| Logical Clock: 5 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 2 |XX| Sent By: Server |XX|
|XX| 10. |XX| Logical Clock: 7 |XX| Message: Hello 1 from Client 2 |XX| Received By: Server |XX| Sent By: Client 2 |XX|
|XX| 11. |XX| Logical Clock: 7 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Received By: Client 5 |XX| Sent By: Server |XX|
|XX| 12. |XX| Logical Clock: 7 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Received By: Client 8 |XX| Sent By: Server |XX|
|XX| 13. |XX| Logical Clock: 7 |XX| Message: <<Broadcasted>> Hello 1 from Client 1 |XX| Received By: Client 6 |XX| Sent By: Server |XX|
|XX| 14. |XX| Logical Clock: 9 |XX| Message: <<Broadcasted>> Hello 1 from Client 2 |XX| Received By: Client 9 |XX| Sent By: Server |XX|
|XX| 15. |XX| Logical Clock: 9 |XX| Message: <<Broadcasted>> Hello 1 from Client 2 |XX| Received By: Client 1 |XX| Sent By: Server |XX|
|XX| 16. |XX| Logical Clock: 9 |XX| Message: <<Broadcasted>> Hello 1 from Client 2 |XX| Received By: Client 7 |XX| Sent By: Server |XX|
|XX| 17. |XX| Logical Clock: 9 |XX| Message: <<Broadcasted>> Hello 1 from Client 2 |XX| Received By: Client 6 |XX| Sent By: Server |XX|
|XX| 18. |XX| Logical Clock: 9 |XX| Message: Hello 1 from Client 6 |XX| Received By: Server |XX| Sent By: Client 6 |XX|
|XX| 19. |XX| Logical Clock: 11 |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX| Received By: Client 2 |XX| Sent By: Server |XX|
|XX| 20. |XX| Logical Clock: 11 |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX| Received By: Client 10 |XX| Sent By: Server |XX|
|XX| 21. |XX| Logical Clock: 11 |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX| Received By: Client 3 |XX| Sent By: Server |XX|
|XX| 22. |XX| Logical Clock: 11 |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX| Received By: Client 5 |XX| Sent By: Server |XX| 
|XX| 23. |XX| Logical Clock: 11 |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX| Received By: Client 9 |XX| Sent By: Server |XX|
```

  



### Prob 1.2- Implementing Vector Clock for Total Ordering of Events

Similar to the above implementation, we have a similar structure of the ouput, just that this time we are using vector clocks and also detecting Potential Causality Violations. All the variables can still be changed in the code as specified. By default, we are using 10 clients again like the problem suggests.

#### 1. Clients Being Registered
```
Registering Client 1 
Registering Client 2
Registering Client 3
Registering Client 4
Registering Client 5
Registering Client 6
Registering Client 7
Registering Client 8
Registering Client 9
Registering Client 10
```
#### 2. Clients Sending Messages --> Server Receiving Messages from Clients
```
|XX| Vector Clock: [3 0 0 1 0 0 0 1 0 0 0] |XX| Received By: Server   |XX| Message: Hello 1 from Client 3 |XX| Sent By: Client 3 |XX|
```

#### 3. Result of CoinToss by Server --> Messages being Broadcasted
Successful Toss:
```
|XXXX| Server : Broadcasting 'Hello 1 from Client 7' from Client 7 |XXXX|

|XX| Vector Clock: [4 0 2 1 0 0 0 1 0 0 0] |XX| Received By: Client 2 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Sent By: Server |XX|
|XX| Vector Clock: [8 1 1 1 1 0 0 4 0 0 0] |XX| Received By: Client 7 |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Sent By: Server |XX|
```
Unsuccessful Toss:
```
|XXXX| Server : Not Broadcasting 'Hello 2 from Client 5' from Client 5 |XXXX|
```
#### 4. Showing Potential Causality Violations
If a Potential Causality Violation occurs when receiving any message, you see an alert in the terminal log, this PCV is then stored in a pcvChannel which is used to display all the detected potential causality violations at once at the end of the code. If after receiving the message there is no detected PCV, no message is displayed. An example of this message is as follows:
```
|XXXX| DETECTED POTENTIAL CAUSALITY VIOLATION |XXXX|
```

#### 5. Total Order of Events
A total ordering of all the messages is then displayed once all the clients have shut down:
```
|XXXX| Total Order of All Messages Below: |XXXX|

|XX| 0. |XX| Logical Clock: [1 0 0 0 0 0 0 1 0 0 0] |XX| Message: Hello 1 from Client 7 |XX| Received By: Server |XX| Sent By: Client 7 |XX|
|XX| 1. |XX| Logical Clock: [3 0 0 1 0 0 0 1 0 0 0] |XX| Message: Hello 1 from Client 3 |XX| Received By: Server |XX| Sent By: Client 3 |XX|
|XX| 2. |XX| Logical Clock: [2 0 0 0 0 0 0 1 1 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX| Received By: Client 8 |XX| Sent By: Server |XX|
|XX| 3. |XX| Logical Clock: [5 1 0 1 0 0 0 1 0 0 0] |XX| Message: Hello 1 from Client 1 |XX| Received By: Server |XX| Sent By: Client 1 |XX|
|XX| 4. |XX| Logical Clock: [2 2 0 0 0 0 0 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX| Received By: Client 1 |XX| Sent By: Server |XX|
|XX| 5. |XX| Logical Clock: [6 1 1 1 0 0 0 1 0 0 0] |XX| Message: Hello 1 from Client 2 |XX| Received By: Server |XX| Sent By: Client 2 |XX|
|XX| 6. |XX| Logical Clock: [4 0 2 1 0 0 0 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX| Received By: Client 2 |XX| Sent By: Server |XX|
|XX| 7. |XX| Logical Clock: [7 1 1 1 1 0 0 1 0 0 0] |XX| Message: Hello 1 from Client 4 |XX| Received By: Server |XX| Sent By: Client 4 |XX|
|XX| 8. |XX| Logical Clock: [8 1 1 1 1 0 0 2 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 4 |XX| Received By: Client 7 |XX| Sent By: Server |XX|
|XX| 9. |XX| Logical Clock: [9 1 1 1 1 1 0 1 0 0 0] |XX| Message: Hello 1 from Client 5 |XX| Received By: Server |XX| Sent By: Client 5 |XX|
|XX| 10. |XX| Logical Clock: [11 1 1 1 1 1 1 1 0 0 0] |XX| Message: Hello 1 from Client 6 |XX| Received By: Server |XX| Sent By: Client 6 |XX|
|XX| 11. |XX| Logical Clock: [10 3 1 1 1 1 0 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 5 |XX| Received By: Client 1 |XX| Sent By: Server |XX|
|XX| 12. |XX| Logical Clock: [2 0 0 0 0 0 2 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX| Received By: Client 6 |XX| Sent By: Server |XX|
|XX| 13. |XX| Logical Clock: [13 1 1 1 1 1 1 1 0 1 0] |XX| Message: Hello 1 from Client 9 |XX| Received By: Server |XX| Sent By: Client 9 |XX|
|XX| 14. |XX| Logical Clock: [15 1 1 1 1 1 1 1 2 1 0] |XX| Message: Hello 1 from Client 8 |XX| Received By: Server |XX| Sent By: Client 8 |XX|
|XX| 15. |XX| Logical Clock: [16 1 1 1 1 1 1 1 2 1 1] |XX| Message: Hello 1 from Client 10 |XX| Received By: Server |XX| Sent By: Client 10 |XX|
|XX| 16. |XX| Logical Clock: [10 1 1 1 1 1 3 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 5 |XX| Received By: Client 6 |XX| Sent By: Server |XX|
|XX| 17. |XX| Logical Clock: [2 0 0 2 0 0 0 1 0 0 0] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX| Received By: Client 3 |XX| Sent By: Server |XX|
```

#### All detected PCVs with the corresponding Message, Message Vector Clock and Local Vector Clock
```
|XXXX| All Potential Causality Violations: |XXXX|

|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [17 1 1 1 1 1 1 1 2 1 1] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 7 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 10 |XX|
|XX| From: Server |XX| To: Client 6 |XX| Message Vector Clock: [8 1 1 1 1 0 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 10 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 4 |XX|
|XX| From: Server |XX| To: Client 3 |XX| Message Vector Clock: [10 1 1 1 1 1 0 1 0 0 0] |XX| Local Vector Clock: [44 7 5 9 5 7 7 8 8 4 5] |XX| Message: <<Broadcasted>> Hello 1 from Client 5 |XX|
|XX| From: Server |XX| To: Client 3 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [44 7 5 10 5 7 7 8 8 4 5] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 6 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [39 7 5 7 5 2 11 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 1 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [29 10 3 3 2 2 5 3 2 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 8 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 6 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 12 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 5 |XX| Message Vector Clock: [12 1 1 1 1 1 1 1 0 0 0] |XX| Local Vector Clock: [29 4 3 3 2 9 5 3 2 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX|
|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [24 4 3 3 2 1 5 1 2 1 1] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 9 3] |XX| Message: <<Broadcasted>> Hello 2 from Client 6 |XX|
|XX| From: Server |XX| To: Client 7 |XX| Message Vector Clock: [17 1 1 1 1 1 1 1 2 1 1] |XX| Local Vector Clock: [35 4 5 3 5 2 5 10 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 10 |XX|
|XX| From: Server |XX| To: Client 2 |XX| Message Vector Clock: [2 0 0 0 0 0 0 1 0 0 0] |XX| Local Vector Clock: [35 4 10 3 5 2 5 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX|
|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [2 0 0 0 0 0 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 10 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX|
|XX| From: Server |XX| To: Client 8 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [29 4 3 3 2 2 5 3 10 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 10 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [35 4 5 3 5 2 5 8 4 4 8] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [8 1 1 1 1 0 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 11 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 4 |XX|
|XX| From: Server |XX| To: Client 1 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [29 11 3 3 2 2 5 3 2 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 9 |XX| Message Vector Clock: [10 1 1 1 1 1 0 1 0 0 0] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 12 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 5 |XX|
|XX| From: Server |XX| To: Client 10 |XX| Message Vector Clock: [12 1 1 1 1 1 1 1 0 0 0] |XX| Local Vector Clock: [35 4 5 3 5 2 5 8 4 4 9] |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX|
|XX| From: Server |XX| To: Client 2 |XX| Message Vector Clock: [17 1 1 1 1 1 1 1 2 1 1] |XX| Local Vector Clock: [35 4 11 3 5 2 5 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 10 |XX|
|XX| From: Server |XX| To: Client 10 |XX| Message Vector Clock: [2 0 0 0 0 0 0 1 0 0 0] |XX| Local Vector Clock: [35 4 5 3 5 2 5 8 4 4 10] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX|
|XX| From: Server |XX| To: Client 5 |XX| Message Vector Clock: [2 0 0 0 0 0 0 1 0 0 0] |XX| Local Vector Clock: [29 4 3 3 2 10 5 3 2 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 7 |XX|
|XX| From: Server |XX| To: Client 2 |XX| Message Vector Clock: [10 1 1 1 1 1 0 1 0 0 0] |XX| Local Vector Clock: [35 4 13 3 5 2 5 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 5 |XX|
|XX| From: Server |XX| To: Client 10 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [35 4 5 3 5 2 5 8 4 4 11] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 8 |XX| Message Vector Clock: [12 1 1 1 1 1 1 1 0 0 0] |XX| Local Vector Clock: [29 4 3 3 2 2 5 3 11 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 6 |XX|
|XX| From: Server |XX| To: Client 3 |XX| Message Vector Clock: [17 1 1 1 1 1 1 1 2 1 1] |XX| Local Vector Clock: [44 7 5 12 5 7 7 8 8 4 5] |XX| Message: <<Broadcasted>> Hello 1 from Client 10 |XX|
|XX| From: Server |XX| To: Client 5 |XX| Message Vector Clock: [4 0 0 1 0 0 0 1 0 0 0] |XX| Local Vector Clock: [35 4 5 3 5 12 5 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 3 |XX|
|XX| From: Server |XX| To: Client 1 |XX| Message Vector Clock: [8 1 1 1 1 0 0 1 0 0 0] |XX| Local Vector Clock: [29 13 3 3 2 2 5 3 2 2 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 4 |XX|
|XX| From: Server |XX| To: Client 10 |XX| Message Vector Clock: [22 4 3 3 2 1 1 1 2 1 1] |XX| Local Vector Clock: [39 7 5 7 5 2 7 8 4 4 13] |XX| Message: <<Broadcasted>> Hello 2 from Client 4 |XX|
|XX| From: Server |XX| To: Client 3 |XX| Message Vector Clock: [24 4 3 3 2 1 5 1 2 1 1] |XX| Local Vector Clock: [44 7 5 13 5 7 7 8 8 4 5] |XX| Message: <<Broadcasted>> Hello 2 from Client 6 |XX|
|XX| From: Server |XX| To: Client 7 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [35 4 5 3 5 2 5 12 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 5 |XX| Message Vector Clock: [14 1 1 1 1 1 1 1 0 1 0] |XX| Local Vector Clock: [35 4 5 3 5 13 5 8 4 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 9 |XX|
|XX| From: Server |XX| To: Client 8 |XX| Message Vector Clock: [17 1 1 1 1 1 1 1 2 1 1] |XX| Local Vector Clock: [35 4 5 3 5 2 5 8 13 4 3] |XX| Message: <<Broadcasted>> Hello 1 from Client 10 |XX|
|XX| From: Server |XX| To: Client 4 |XX| Message Vector Clock: [24 4 3 3 2 1 5 1 2 1 1] |XX| Local Vector Clock: [42 7 5 7 12 7 7 8 8 4 3] |XX| Message: <<Broadcasted>> Hello 2 from Client 6 |XX|
|XX| From: Server |XX| To: Client 3 |XX| Message Vector Clock: [29 4 3 3 2 2 5 3 2 2 3] |XX| Local Vector Clock: [44 7 5 14 5 7 7 8 8 4 5] |XX| Message: <<Broadcasted>> Hello 2 from Client 10 |XX|
|XX| From: Server |XX| To: Client 7 |XX| Message Vector Clock: [24 4 3 3 2 1 5 1 2 1 1] |XX| Local Vector Clock: [44 7 5 7 5 7 7 14 8 4 5] |XX| Message: <<Broadcasted>> Hello 2 from Client 6 |XX|
```





### Problem 2- Implementing Bully Algorithm in a simulated distributed system

####  Implementing Bully Algorithm for leader election in a distributed system.

To run the program, please follow the following steps exactly:
##### Step 1
Run the  command as documented above in 5 different terminals. Wait until all of the terminal messages display ```Enter the Node id [1-5]:```

##### Step 2
One by one, put in numbers 1 through 5 followed by Enter in each of the terminals (1 in the first terminal, 2 in the second terminal, etc.) to initialize all clients. Now we have the simulation set up.
```
Enter the Node id [1-5]: 3
This server is now on: 127.0.0.1:5002
Is this node reentering the system? Type in Yes and press senter if so, if not: 
Press enter to communicate with coordinator.
Sending Request to Coordinator

Not able to connect to coordinator -> Invoking Election to find new coordinator
Sending an Election to Node:  4
Node  4 is not available: Communication Failed
Sending an Election to Node:  5
Node  5 is not available: Communication Failed
Node  3 has been elected as the new coordinator

Is this node reentering the system? Type in Yes and press senter if so, if not: 
Press enter to communicate with coordinator.
```
#### PS1 Please note that at the start, everytime you put in a value that is not 5 (when 5 hasn't been started), the node will try to get in touch with 5 but when it doesn't it'll start an election- This is by design once you have started nodes all through 1 to 5, the system will be ready for the rest of the scenarios (with 5 as the coordinator).


The question of whether it is reentring the system is in case the node is rejoining the system, if yes it will start an election to reinstate the coordinator with everyone. At the start, you don't need to worry about it
##### Step 3
The setup is now complete, with 5 different nodes and Node 5 is the Coordinator. You can test this by hitting Enter in any of the clients to communicate with the coordinator and you should see this pop up.
```
Sending Request to Coordinator

Request to Coordinator Successful 5

Press enter to communicate with coordinator.

```
On coordinator's side we get:
```
Getting Request from Node:  3
```
### At this point the basic Algorithm has been set up (i.e. Point 1)- Please have a look at the code to how it's done. We will now go through the other scenarios as listed.
  

### Point 2.1 Worst Case Scenario
To simulate the worst case scenario, kill the terminal that is node 5 (i.e. the coordinator) and then go to Node 1 and hit enter to communicate with coordinator. 
This will make it realize the coordinator is down since it won't receive an acknowledgement when contacting and it will start an election between all the remaining nodes (since it has the lowest id).
On Node 1 side we get the following logs:
```
Sending Request to Coordinator

Not able to connect to coordinator -> Invoking Election to find new coordinator
Sending an Election to Node:  2
Received acknowledgement from Node 2
Sending an Election to Node:  3
Received acknowledgement from Node 3
Sending an Election to Node:  4
Received acknowledgement from Node 4
Sending an Election to Node:  5
Node  5 is not available: Communication Failed

Press enter to communicate with coordinator.
Node  4 has been elected as the new coordinator
```
This election as per the protocol prompted every node to start their own election after sending their acknowledgement- and for example we can have a look at what the log says in Node 3: You can see it got election from 1 as well as 2 and subsequently it also started its own election:
```
Getting Election Request from Node:  1
Sending Acknowledgement for Election to Node  1
Getting Election Request from Node:  2
Sending Acknowledgement for Election to Node  2
Sending an Election to Node:  4
Received acknowledgement from Node 4
Sending an Election to Node:  5
Node  5 is not available: Communication Failed
Node  4 has been elected as the new coordinator
```
### Point 2.2 Best Case Scenario:
To simulate the best case scenario we kill the Node 5 again, and for that please open a new terminal and run the command ```go run bullyalgorithm.go``` to start it and then enter 5 and click "Yes" so that it starts an election again to reinstate itself as coordinator.
We get the following log for this:
```
Enter the Node id [1-5]: 5
This server is now on: 127.0.0.1:5004
Is this node reentering the system? Type in Yes and press senter if so, if not: 
Press enter to communicate with coordinator.
Sending Request to Coordinator

Getting Request from Node:  5
Request to Coordinator Successful 5

Is this node reentering the system? Type in Yes and press senter if so, if not: 
Press enter to communicate with coordinator.
```
On other servers the fact that 5 has again joined as coordinator is also shown as follows:
```
Node  5 has been elected as the new coordinator
```

Now we kill it again, and have Node 4 try to send a request to Coordinator and we see what happens:
```
Sending Request to Coordinator

Not able to connect to coordinator -> Invoking Election to find new coordinator
Sending an Election to Node:  5
Node  5 is not available: Communication Failed
Node  4 has been elected as the new coordinator
```
Nice and simple, it has been chosen as the coordinator and it's reflected in all other nodes as well.

### At this point, we see Point 2 has been fulfilled by showing the Best and Worst Cases

### Point 3.a Newly Elected Coordinator Failure During Announcement:

Now, originally since there's no need the line which causes a delay when announcing itself as a leader has just a 1 second delay (very minimal). But for the rest of the scenarios (i.e. 3, and 4)- please change this line to be delayed for 5 seconds instead making it easier to see the scenarios listed below- It is Line 102 in the original file, if not you can find it in the ```setCoordinator()`` function. Please contact me on Tele @jyotit_kaushal if you can't find it.
```
time.Sleep(time.Duration(1) * time.Second)
```
Once this is done, you can simulate the case where the elected coordinator fails during announcement. In this case, you'll essentially have two copies of the coordinators based on which node you are. But either way when you hit enter to try to talk with the coordinator no matter which node you are- you'll realize the coordinator is down and this will start a new election. Hence, there will not be any functional issues with the program. 

### Point 3.b The failed node is not the newly elected coordinator.
If the failed node is not the newly elected coordinator then there will not be an issue within the system itself, just that when the newly elected coordinator sends a message to update the node that is now a failure, it will log a message like ```"Node X is not available: Communication Failed```. Everything will still be just fine

### Point 3 is hence satisfied by this implementation.

### Point 4 Multiple Nodes Start the Election
You can simulate this by hitting Enter on two different nodes at relatively the same time (within 5 seconds) and you'll now have started an election in 2 nodes simultaneously. This again would not break the system since the implementation of bully algorithm is not affected by this, all the messages will be displayed in the same fashion, just that there will be more messages now on both ends but in the end the right coordinator will be picked which will be told to every node so no worries. You will basically see the coordinator messages in the correct order.

 ### Point 4 is hence satisfied by this implementation.

### Point 5 An arbitrary node silently leaves the network

This is kind of a subset of Point 3 and is going to be handled in a similar fashion. If a non-coordinator node silently leaves the network, no error or election will be raised since there's no need for that and the system will continue to run smoothly. Just that when it rejoins another election will be started just to make sure it wasn't a node that had a higher id than whatever current coordinator is present.

If a node leaves that was actually the coordinator,  then eventually one of the nodes will realize it's not present anymore when it tries to connect with the coordinator and at this point an election will start and you'll have a new coordinator no issues.

### Point 5 is hence satisified by this implementation.

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
| Ivy Without Fault Tolerance | 19.79 seconds |
| Ivy With Fault Tolerance- Scenario 1|21.45 seconds|
| Ivy With Fault Tolerance- Scenario 2 | 21.28 seconds |
| Ivy With Fault Tolerance- Scenario 3|21.38 seconds|
| Ivy With Fault Tolerance- Scenario 4 |21.87 seconds  |
| Ivy With Fault Tolerance- Scenario 5|22.14 seconds|
