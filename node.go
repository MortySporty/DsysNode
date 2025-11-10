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
	}

	client := proto.NewITUDatabaseClient(conn)
	client.SendMessages(context.Background(), &proto.Message{
		Message: []string{"hello"},
	})

	if err != nil {
		log.Fatalf("WE DID NOT RECIEVE OR FAILED TO SEND")
	}
}
