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

var numIter = 16

func testInstructors() error {

	for i := 0; i < numIter; i++ {

		var req pb.EnrollInstructorReq
		var ireq pb.GetInstructorReq

		ins := util.GetNewInstructor()
		req.Instructor = &ins

		log.WithFields(log.Fields{"instructor req": req}).Debug("Enrolling instructor")
		r, err := client.EnrollInstructor(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to enroll instructor")
			return err
		}
		log.WithFields(log.Fields{"instructor response": r}).Debug("Enrolled instructor with key")

		ireq.InstructorKey = r.InstructorKey
		gr, err := client.GetInstructor(context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor")
			return err
		}
		util.RegisterEnrolledInstructorID(ireq.InstructorKey)
		log.WithFields(log.Fields{"instructorInfo": gr.Info}).
			Debug("Get instructor success")
	}
	return nil
}

func testUsers() error {

	for i := 0; i < numIter; i++ {

		var req pb.EnrollUserReq
		var ureq pb.GetUserReq

		u := util.GetNewUser()
		req.User = &u

		log.WithFields(log.Fields{"user req": req}).Debug("Enrolling user")
		r, err := client.EnrollUser(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to enroll user")
			return err
		}
		log.WithFields(log.Fields{"user response": r}).Debug("Enrolled user with key")

		ureq.UserKey = r.UserKey
		gr, err := client.GetUser(context.Background(), &ureq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get user")
			return err
		}
		log.WithFields(log.Fields{"userInfo": gr.Info}).
			Debug("Get user success")
	}
	return nil
}

func getData() error {

	var allreq pb.GetSessionsReq

	gr, err := client.GetSessions(context.Background(), &allreq)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get all sessions")
		return err
	}
	log.WithFields(log.Fields{"allsessionResponse": gr}).
		Debug("Get all session success")

	//var alluser pb.GetUsersReq
	//ur, err := client.GetUsers(context.Background(), &alluser)
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Error("Failed" +
	//		" to get all users")
	//	return err
	//}
	//log.WithFields(log.Fields{"allusersResponse": ur}).
	//	Debug("Get all user success")

	//var allins pb.GetInstructorsReq
	//ir, err := client.GetInstructors(context.Background(), &allins)
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Error("Failed" +
	//		" to get all instructors")
	//	return err
	//}
	//log.WithFields(log.Fields{"allinsResponse": ir}).
	//	Debug("Get all instructor success")

	return nil
}

func main() {

	// open a file
	f, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
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

	err = getData()
	if err != nil {
		log.Error("Session test failed")
		return
	}
	return
}
