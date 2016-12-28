package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	pb "server/rpcdefsql"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client pb.ServerSvcClient

func testInstructors() error {

	var numIns = 8
	var allreq pb.GetInstructorsReq
	for i := 0; i < numIns; i++ {

		var req pb.EnrollInstructorReq
		var ireq pb.GetInstructorReq

		ins := GetNewInstructor()
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
		RegisterEnrolledInstructorID(ireq.InstructorKey)
		log.WithFields(log.Fields{"instructorInfo": gr.Info}).
			Debug("Get instructor success")
	}

	gr, err := client.GetInstructors(context.Background(), &allreq)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get all instructors")
		return err
	}
	log.WithFields(log.Fields{"allInsResponse": gr}).
		Debug("Got all instructor success")

	return nil
}

func testUsers() error {

	var allreq pb.GetUsersReq
	var numUsers = 4
	for i := 0; i < numUsers; i++ {

		var req pb.EnrollUserReq
		var ureq pb.GetUserReq

		u := GetNewUser()
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

	gr, err := client.GetUsers(context.Background(), &allreq)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get all users")
		return err
	}
	log.WithFields(log.Fields{"allUserResponse": gr}).
		Debug("Got all user success")

	return nil
}

func testSessions() error {

	var allreq pb.GetSessionsReq
	var numSessions = 4
	for i := 0; i < numSessions; i++ {

		var req pb.PostSessionReq
		var greq pb.GetSessionReq

		err, s := GetNewSession()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get new session")
			return err
		}
		req.Info = &s

		log.WithFields(log.Fields{"session req": req}).Debug("Posting session")
		r, err := client.PostSession(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to post session")
			return err
		}
		log.WithFields(log.Fields{"session response": r}).Debug("Posted session with key")

		greq.SessionKey = r.SessionKey
		gr, err := client.GetSession(context.Background(), &greq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get all sessions")
			return err
		}
		log.WithFields(log.Fields{"session": gr.Info}).
			Debug("Get session success")
	}

	gr, err := client.GetSessions(context.Background(), &allreq)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get all sessions")
		return err
	}
	log.WithFields(log.Fields{"allsessionResponse": gr}).
		Debug("Get all session success")

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

	address := GetUrl()
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

	err = testUsers()
	if err != nil {
		log.Error("Users test failed")
		return
	}

	log.Printf("\n\n")
	err = testInstructors()
	if err != nil {
		log.Error("Instructor test failed")
		return
	}

	log.Printf("\n\n")
	err = testSessions()
	if err != nil {
		log.Error("Session test failed")
		return
	}
}
