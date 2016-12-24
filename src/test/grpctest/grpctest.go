package main

import (
	"fmt"
	"os"
	"test/util"

	log "github.com/Sirupsen/logrus"

	pb "server/rpcdef"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client pb.ServerSvcClient

func testSessions() error {
	numSessions := 32

	for i := 0; i < numSessions; i++ {

		var req pb.PostSessionReq
		var greq pb.GetSessionReq

		s := util.GetNewSession()
		req.Info = &s
		_, err := client.PostSession(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to post session")
			return err
		}
		log.WithFields(log.Fields{"session": req}).Debug("Posted session")

		gr, err := client.GetSession(context.Background(), &greq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get all sessions")
			return err
		}
		log.WithFields(log.Fields{"session": gr.Info}).
			Debug("Get session success")
	}
	return nil
}

func main() {

	// open a file
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(os.Stdout)
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	address := util.GetUrl()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewServerSvcClient(conn)

	r, err := client.GetStatus(context.Background(),
		&pb.ServerSvcStatusReq{Name: "grpctest"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	err = testSessions()
	if err != nil {
		log.Error("Session test failed")
		return
	}
}
