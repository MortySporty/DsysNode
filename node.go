package main

import (
	proto "DsysNode/grpc"
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUDatabaseServer
	messages []string
}

func main() {
	ID, _ := strconv.ParseInt(os.Args[1], 10, 32)

	nodePorts := map[string]int{
		"node1": 8001,
		"node2": 8002,
		"node3": 8003,
		"node4": 8004,
		"node5": 8005,
	}

	delete(nodePorts, fmt.Sprintf("node%d", ID))

	server := &ITU_databaseServer{messages: []string{}}
	go server.start_server(ID)

	var Channels []chan int

	for _, value := range nodePorts {
		ch := make(chan int)
		go clientConnect(value, ch, ID)
		Channels = append(Channels, ch)
	}

	for {
		state := 0 
		clock := 0;

		// if channel comes back successfully from connected, update clock and state

			sleepDuration := time.Duration(rand.IntN(10)) * time.Second;
			time.Sleep(sleepDuration);
			for {
					for _, ch := range Channels {
					ch <- 1 
			}
		}
	}
}

func (s *ITU_databaseServer) start_server(ID int64) {

	port := fmt.Sprintf(":800%d", ID)

	grpcserver := grpc.NewServer()
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("SERVER WONT WORK")
	}

	proto.RegisterITUDatabaseServer(grpcserver, s)

	fmt.Printf("Starting gRPC server on %s", port)

	if err := grpcserver.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func clientConnect(port int, ch chan int, ID int64) {

	for {
		conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Printf("Could not create client connection to %d: %v. Retrying in 2 seconds...", port, err)
			time.Sleep(2 * time.Second)
			continue
		}

		client := proto.NewITUDatabaseClient(conn)

		_, err = client.TestConnection(context.Background(), &proto.Empty{})

		if err != nil {
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("\nConnected to server on port %d \n", port)
			connected(client, ch, ID)
		}
	}
}

func connected(client proto.ITUDatabaseClient, ch chan int, ID int64) {
	// at a random interval between 0 - 10, a message (1 int) is sent in the channel, 
	// which indicates the node wants to access the critical section,
	// set local variable to await state and send to other,
	// this should course the client to send requests to other clients for access
	ting bool;

	msg := <-ch

	// based on this message, we decide if we want to access the critical section
	client.

}

func (s *ITU_databaseServer) TestConnection(ctx context.Context, in *proto.Empty) (*proto.Message, error) {

	// placeholder
	return &proto.Message{
		Message: []string{""},
		Tick:    100,
	}, nil
}
