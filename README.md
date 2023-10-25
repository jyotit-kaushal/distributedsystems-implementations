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

```