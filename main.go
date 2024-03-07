package main

import (
	"fmt"
	a "ZK-leader-election/leaderElection"
	b "ZK-leader-election/udpcommunication"
	"time"
)

func main() {
	server1 := a.NewLE(0)
	server2 := a.NewLE(1)
	server3 := a.NewLE(2)
	sender,_ := b.NewUDPC()
	sender.AddNode(0)
	startTime := time.Now()
	server1.StartServer(sender)
	sender.AddNode(1)
	server2.StartServer(sender)
	sender.AddNode(2)
	server3.StartServer(sender)


	for server1.Leader() == server2.Leader() == server3.Received()  {
	}
	endTime := time.Now()
	fmt.Printf("Leader %d elected in %d milliseconds\n", server1.Leader(), endTime.Sub(startTime).Milliseconds())

	time.Sleep(5 * time.Second)
	for server1.Leader() == server2.Leader() == server3.Received()  {
		fmt.Printf("Leaders %d %d %d elected \n", server1.Leader(),server2.Leader(),server3.Leader())
		time.Sleep(5 * time.Second)
	}
	for true {
		fmt.Printf("Leader %d elected in %d milliseconds\n", server1.Leader(), endTime.Sub(startTime).Milliseconds())
		time.Sleep(1 * time.Second)
	}
	return

}
