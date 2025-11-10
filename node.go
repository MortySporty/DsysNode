package main

import (
	proto "DsysNode/grpc"
	"context"
	"fmt"
	"log"
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

	for _, value := range nodePorts {
		go clientConnect(value)
	}

	for {
		time.Sleep(2 * time.Second)
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

func clientConnect(port int) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Count not connect to port: %d", port)
		return
	}

	client := proto.NewITUDatabaseClient(conn)

	_, err = client.ServerSend(context.Background(), &proto.Empty{})

	if err != nil {
		fmt.Printf("\nCould not connect to server on port %d: %v", port, err)
	} else {
		fmt.Printf("\nconnected to server on port %d", port)
	}

	// log.Printf("Successfully connected and called server on port %d", port)

}

func send(client proto.ITUDatabaseClient, ID string, time int, input string) {

	msgText := fmt.Sprintf("%s: %s (LOGICAL TIME: %d)", ID, input, time)

	client.ClientSend(context.Background(), &proto.Message{
		Message: []string{msgText},
		Tick:    int32(time),
	})
}

func (s *ITU_databaseServer) ServerSend(ctx context.Context, in *proto.Empty) (*proto.Message, error) {

	// placeholder
	return &proto.Message{
		Message: []string{"hej"},
		Tick:    100,
	}, nil
}
