package main

import (
	"log"
	"test/util"

	pb "server/rpcdef"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const ()

func main() {
	address := util.GetUrl()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewServerSvcClient(conn)

	r, err := c.GetStatus(context.Background(),
		&pb.ServerSvcStatusReq{Name: "grpctest"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
