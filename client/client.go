package main

import (
	"context"
	"fmt"
	"time"

	api "github.com/gannonbarnett/cloud/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connect(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Error dialing: %v\n", err)
		return
	}
	defer conn.Close()

	stream, err := api.NewCloudClient(conn).Handle(context.Background())
	if err != nil {
		fmt.Printf("Error creating stream: %v\n", err)
		return
	}

	go func() {
		stream.Send(&api.FromDevice{
			Name: "client1",
		})
	}()

	for {
		reply, err := stream.Recv()
		if err != nil {
			fmt.Printf("Error reading from stream: %v\n", err)
			return
		}

		fmt.Printf("received: %v\n", reply)
	}
}

func main() {
	addr := "0.0.0.0:9000"
	fmt.Printf("Connecting to cloud at %v\n", addr)
	for range time.NewTicker(5 * time.Second).C {
		connect(addr)
		fmt.Printf("...\n")
	}
}
