package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"

	pb "server/rpcdefsql"

	"crypto/tls"
	"crypto/x509"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var client pb.SFServerClient

func testPostReview() error {

	numSessions := 2
	numIns := 2

	for i := 0; i < numIns; i++ {

		var r pb.UserInstructorReview
		var req pb.PostInstructorReviewReq
		req.Review = &r
		req.Review.InstructorID = int32(i) + 1
		req.Review.UserID = 1
		req.Review.InstructorReview = "my review" + strconv.Itoa(numIns)
		_, err := client.PostInstructorReview(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to push instructor review")
			return err
		}
	}

	for i := 0; i < numSessions; i++ {

		var r pb.UserSessionReview
		var req pb.PostSessionReviewReq
		req.Review = &r
		req.Review.SessionID = int32(i) + 1
		req.Review.UserID = 1
		req.Review.SessionReview = "my session review" +
			strconv.Itoa(numSessions)
		_, err := client.PostSessionReview(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to push session review")
			return err
		}
	}
	return nil
}

func testInstructorImages() error {

	numIns := GetNumInstructorImages()
	for i := 0; i < numIns; i++ {

		log.WithFields(log.Fields{"IMG id ": i}).Debug("Pushing image")
		var req pb.PostInstructorDisplayImgReq
		req.InstructorInfoID = int32(i) + 1
		imgData, err := ioutil.ReadFile(InsImages[i])
		req.Blob = make([]byte, len(imgData))
		copy(req.Blob, imgData)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to read jpeg")
			return err
		}
		_, err = client.PostInstructorDisplayImg(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to push instructor image")
			return err
		}
	}
	return nil
}

func testUserCC(numUsers int) error {

	for i := 0; i < numUsers; i++ {

		var req pb.SubscribeUserReq
		var ireq pb.GetUserCCReq

		cc := GetNewCC(i + 1)
		req.PayCard = &cc

		log.WithFields(log.Fields{"userCC req": req}).Debug("Adding user CC")
		r, err := client.SubscribeUser(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to add user cc")
			return err
		}

		ireq.CcID = r.CcID
		gr, err := client.GetUserCC(
			context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get user cc")
			return err
		}
		log.WithFields(log.Fields{"userCC": gr}).
			Debug("Received user cc with ccid")

		ireq.CcID = 0
		ireq.UserID = int32(i + 1)
		gr, err = client.GetUserCC(
			context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get user cc with userid")
			return err
		}
		log.WithFields(log.Fields{"userCC": gr}).
			Debug("Received user cc with userid")

	}

	return nil
}

func testInsBankAcct() error {

	numIns := GetNumInstructors()
	for i := 0; i < numIns; i++ {

		var req pb.RegisterInstructorBankAcctReq
		var ireq pb.GetInstructorBankAcctReq

		bankAcct := GetNewBankAcct(i + 1)
		req.BankAcct = &bankAcct

		log.WithFields(log.Fields{"ins bank req": req}).Debug("Adding instructor bank")
		r, err := client.RegisterInstructorBankAcct(context.Background(),
			&req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to add instructor bank")
			return err
		}

		ireq.BankAcctID = r.BankAcctID
		gr, err := client.GetInstructorBankAcct(
			context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor bank")
			return err
		}
		log.WithFields(log.Fields{"insBank": gr}).
			Debug("Received bank acct with bank ID")

		ireq.BankAcctID = 0
		ireq.InstructorID = int32(i + 1)
		gr, err = client.GetInstructorBankAcct(
			context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor bank with instructor id")
			return err
		}
		log.WithFields(log.Fields{"insBank": gr}).
			Debug("Received bank acct with ins id")
	}

	return nil
}

func testInstructors() error {

	numIns := GetNumInstructors()
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

		log.WithFields(log.Fields{"instructor req": req}).Debug("Enrolling instructor")
		//r, err := client.PostInstructorDisplayImg(context.Background(),
		//	&req)
		//if err != nil {
		//	log.WithFields(log.Fields{"error": err}).Error("Failed" +
		//		" to enroll instructor")
		//	return err
		//}
		//log.WithFields(log.Fields{"instructor response": r}).Debug("Enrolled instructor with key")

		ireq.InstructorKey = r.InstructorKey
		gr, err := client.GetInstructor(context.Background(), &ireq)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor")
			return err
		}
		RegisterEnrolledInstructor(ireq.InstructorKey, gr.Info)
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

func testLoginIns() error {

	numIns := GetNumInstructors()
	for i := 0; i < numIns; i++ {

		var req pb.LoginReq
		err, uid := getEnrolledInstructorID()
		if err != nil {
			return err
		}

		i := getEnrolledIns(uid)
		req.Email = i.Email
		req.Password = i.Password

		log.WithFields(log.Fields{"user req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to enroll user")
			return err
		}
		log.WithFields(log.Fields{"login response": r}).
			Debug("Logged in user with key")
	}

	// Invalid password
	log.Debug("INCORRECT INS PASSWORD test")
	for i := 0; i < numIns; i++ {

		var req pb.LoginReq
		err, uid := getEnrolledInstructorID()
		if err != nil {
			return err
		}

		u := getEnrolledIns(uid)
		req.Email = u.Email
		req.Password = "incorrectpwd"

		log.WithFields(log.Fields{"ins req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err == nil {
			log.WithFields(log.Fields{"resp": r}).Error("Failed" +
				" to invalidate ins login")
			return errors.New("Invalidate ins login failure")
		}
		log.WithFields(log.Fields{"expected error": err}).
			Debug("invalid ins login succeeded")
	}

	// Invalid email
	log.Debug("INCORRECT INS EMAIL test")
	for i := 0; i < numIns; i++ {

		var req pb.LoginReq
		req.Email = "aa@bcac.come"
		req.Password = "incorrectpwd"

		log.WithFields(log.Fields{"user req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err == nil {
			log.WithFields(log.Fields{"resp": r}).Error("Failed" +
				" to invalidate user login")
			return errors.New("Invalidate user login failure")
		}
		log.WithFields(log.Fields{"expected error": err}).
			Debug("invalid ins login succeeded")
	}

	return nil
}

func testLoginUser(numUsers int) error {

	for i := 0; i < numUsers; i++ {

		var req pb.LoginReq
		err, uid := getEnrolledUsersID()
		if err != nil {
			return err
		}

		u := getEnrolledUser(uid)
		req.Email = u.Email
		req.Password = u.Password

		log.WithFields(log.Fields{"user req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to enroll user")
			return err
		}
		log.WithFields(log.Fields{"login response": r}).
			Debug("Enrolled user with key")
	}

	// Invalid password
	log.Debug("INCORRECT USER PASSWORD test")
	for i := 0; i < numUsers; i++ {

		var req pb.LoginReq
		err, uid := getEnrolledUsersID()
		if err != nil {
			return err
		}

		u := getEnrolledUser(uid)
		req.Email = u.Email
		req.Password = "incorrectpwd"

		log.WithFields(log.Fields{"user req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err == nil {
			log.WithFields(log.Fields{"resp": r}).Error("Failed" +
				" to invalidate user login")
			return errors.New("Invalidate user login failure")
		}
		log.WithFields(log.Fields{"expected error": err}).
			Debug("invalid user login succeeded")
	}

	// Invalid email
	log.Debug("INCORRECT USER EMAIL test")
	for i := 0; i < numUsers; i++ {

		var req pb.LoginReq
		req.Email = "aa@bcac.come"
		req.Password = "incorrectpwd"

		log.WithFields(log.Fields{"user req": req}).
			Debug("loggin in user")
		r, err := client.Login(context.Background(), &req)
		if err == nil {
			log.WithFields(log.Fields{"resp": r}).Error("Failed" +
				" to invalidate user login")
			return errors.New("Invalidate user login failure")
		}
		log.WithFields(log.Fields{"expected error": err}).
			Debug("invalid user login succeeded")
	}

	return nil
}

func testUsers(numUsers int) error {

	var allreq pb.GetUsersReq
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
		RegisterEnrolledUser(ureq.UserKey, gr.Info)
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

	numSessions := GetNumSessions()
	var allreq pb.GetSessionsReq
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

func testTwilioJwT() error {
	req := &pb.TwilioJwtReq{"10", "secret"}
	reply, err := client.GetTwilioJwtToken(context.Background(), req)
	if err != nil {
		return err
	}
	log.Printf("JWT token %s", reply.JwToken)
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

	url := getUrl()
	addressAndPort := getAddressAndPort()
	certificate, err := tls.LoadX509KeyPair(
		"../cert/127.0.0.1.crt",
		"../cert/127.0.0.1.key",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("../cert/My_Root_CA.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   url,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial(addressAndPort, dialOption)
	if err != nil {
		log.Fatalf("failed to dial server: %s", err)
	}
	defer conn.Close()

	client = pb.NewSFServerClient(conn)

	r, err := client.GetStatus(context.Background(),
		&pb.SFServerStatusReq{Name: "grpctest"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	//var req pb.CleanupAllDBsReq
	//_, err = client.CleanupAllDBs(context.Background(), &req)
	//if err != nil {
	//	log.Fatalf("could not cleanup DB: %v", err)
	//}
	//log.Printf("Cleaned up DB", r.Message)

	err = testUsers(4)
	if err != nil {
		log.Error("Users test failed")
		return
	}

	log.Printf("\n\n")
	log.Debug("USER TEST")
	err = testLoginUser(4)
	if err != nil {
		log.Error("user login test failed")
		return
	}

	log.Printf("\n\n")
	log.Debug("INS TEST")
	err = testInstructors()
	if err != nil {
		log.Error("Instructor test failed")
		return
	}

	//log.Printf("\n\n")
	//log.Debug("LOGIN INS TEST")
	//err = testLoginIns()
	//if err != nil {
	//	log.Error("user login test failed")
	//	return
	//}

	err = testInstructorImages()
	if err != nil {
		log.Error("Instructor image test failed")
		return
	}

	err = testPostReview()
	if err != nil {
		log.Error("review test FAILED")
		return
	}

	err = testInsBankAcct()
	if err != nil {
		log.Error("ins bank test FAILED")
		return
	}

	err = testUserCC(4)
	if err != nil {
		log.Error("user cc test FAILED")
		return
	}

	log.Printf("\n\n")
	log.Debug("SESSIONS TEST")
	err = testSessions()
	if err != nil {
		log.Error("Session test failed")
		return
	}

	err = testTwilioJwT()
	if err != nil {
		log.Error("Failed getting Twillio JWT")
		return
	}
}
