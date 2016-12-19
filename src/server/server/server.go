package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port          = ":50051"
	soulFitDB     = "SoulFitDB"
	internalError = "internalError"
	successError  = "success"
)

var lastUserUserID uint64
var db *storm.DB

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
	err := db.All(&resp.Session)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, nil
	}
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func (ser *server) PostSession(ctx context.Context,
	in *pb.PostSessionRequest) (*pb.PostSessionReply, error) {

	var resp pb.PostSessionReply
	log.Debug("Enroll Instructor request")
	var s pb.Session
	s.Info = in.Info
	s.ID = getRandomID()
	err := db.Save(&s)
	if err != nil {
		log.WithFields(log.Fields{"session": s, "error": err}).Error("Failed" +
			" to write to DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"session": s}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func (s *server) EnrollInstructor(ctx context.Context,
	in *pb.EnrollInstructorRequest) (*pb.EnrollInstructorResponse, error) {

	var resp pb.EnrollInstructorResponse
	log.Debug("Enroll Instructor request")
	var i pb.Instructor
	i.Person = in.Instructor
	i.ID = getRandomID()
	err := db.Save(&i)
	if err != nil {
		log.WithFields(log.Fields{"instructor": i, "error": err}).Error("Failed" +
			" to write to DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"instructor": i}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func getRandomID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserRequest) (*pb.EnrollUserResponse, error) {

	var resp pb.EnrollUserResponse
	log.Debug("Enroll User request")
	var u pb.User
	u.Person = in.User
	u.ID = getRandomID()
	err := db.Save(&u)
	if err != nil {
		log.WithFields(log.Fields{"user": u, "error": err}).Error("Failed" +
			" to write to DB for user")

		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"user": u}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func main() {
	// open a file
	f, err := os.OpenFile("server.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err = storm.Open(soulFitDB)

	s := grpc.NewServer()
	log.Debug("registering server...")
	pb.RegisterServerSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
