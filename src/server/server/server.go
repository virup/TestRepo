package main

import (
	"log"
	"net"

	pb "server/rpcdef"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

//  Send hello
func (s *server) GetStatus(ctx context.Context,
	in *pb.ServerSvcStatusRequest) (*pb.ServerSvcStatusResponse, error) {
	return &pb.ServerSvcStatusResponse{Message: "Hello " + in.Name}, nil
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsRequest) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	return &resp, nil
}

func (s *server) PostSession(ctx context.Context,
	in *pb.PostSessionRequest) (*pb.PostSessionReply, error) {

	var resp pb.PostSessionReply
	return &resp, nil
}

func (s *server) EnrollInstructor(ctx context.Context,
	in *pb.EnrollInstructorRequest) (*pb.EnrollInstructorResponse, error) {

	var resp pb.EnrollInstructorResponse
	return &resp, nil
}

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserRequest) (*pb.EnrollUserResponse, error) {

	var resp pb.EnrollUserResponse
	return &resp, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServerSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
