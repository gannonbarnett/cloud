package main

import (
	"context"
	"fmt"
	"io"
	"net"

	api "github.com/gannonbarnett/cloud/api"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type CloudServer struct {
	api.UnimplementedCloudServer

	db *pgx.Conn
}

func NewCloudServer(dbAddr string) *CloudServer {
	conn, err := pgx.Connect(context.Background(), dbAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to db: %v\n", dbAddr)

	return &CloudServer{
		db: conn,
	}
}

func (s *CloudServer) Handle(stream api.Cloud_HandleServer) error {
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		err := fmt.Errorf("couldn't get peer")
		fmt.Printf("Stream rejected: %v\n", err)
		return err
	}
	addr := p.Addr
	fmt.Printf("%v: connected\n", addr)
	go func() {
		err := stream.Send(&api.ToDevice{Name: "server"})
		if err != nil {
			fmt.Printf("%v: error sending: %v\n", addr, err)
		}
	}()

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("EOF\n")
		} else if err != nil {
			fmt.Printf("Stream: %v\n", err)
			return err
		}
		fmt.Printf("%v: '%+v'\n", addr, recv)
	}

	return nil
}

func main() {
	dbAddr := "postgres://admin@localhost:5432/pgdata"

	addr := "0.0.0.0:9000"
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	handler := NewCloudServer(dbAddr)
	grpcServer := grpc.NewServer()
	api.RegisterCloudServer(grpcServer, handler)

	fmt.Printf("Cloud server running on %v\n", addr)
	grpcServer.Serve(conn)
	fmt.Printf("Cloud server stopped.")
}
