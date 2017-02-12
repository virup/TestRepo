package main

import (
	"fmt"
	"net"
	"os"
	pb "server/rpcdefsql"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc/credentials"
)

const (
	port      = ":8099"
	soulFitDB = "SoulFitDB"
)

var lastUserUserID uint64

// XXX Consolidate user, instructor and session info object
// handling and DB handling through a common interface

// server is used to implement helloworld.GreeterServer.
type server struct{}

//  Send hello
func (s *server) GetStatus(ctx context.Context,
	in *pb.ServerSvcStatusReq) (*pb.ServerSvcStatusReply, error) {
	return &pb.ServerSvcStatusReply{Message: "Hello " + in.Name}, nil
}

func (s *server) CleanupAllDBs(ctx context.Context,
	in *pb.CleanupAllDBsReq) (*pb.CleanupAllDBsReply, error) {

	var err error
	var resp pb.CleanupAllDBsReply

	err = db.DropTableIfExists(&pb.SessionInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete session table")
	}

	err = db.DropTableIfExists(&pb.InstructorInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete instructor table")
	}

	err = db.DropTableIfExists(&pb.UserInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete user table")
	}

	log.Debugf("Cleaned up all tables")

	return &resp, nil
}

func (s *server) RecordEvent(ctx context.Context,
	in *pb.RecordEventReq) (*pb.RecordEventReply, error) {

	var resp pb.RecordEventReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to login")
		return &resp, err
	}
	return &resp, nil
}

func initGprcServer() {

	certificate, err := tls.LoadX509KeyPair(
		"../cert/127.0.0.1.crt",
		"../cert/127.0.0.1.key",
	)
	if err != nil {
		log.Fatalf("failed to load cert: %s", err)
	}

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("../cert/My_Root_CA.crt")
	if err != nil {
		log.Fatalf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append client certs")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))
	s := grpc.NewServer(serverOption)

	log.Debug("registering server...")
	pb.RegisterServerSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Debug("registered server...")

}

var db *gorm.DB

func initDB() error {
	var err error
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		log.Errorf("Opening of DB failed with error '%v'",
			err)
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&pb.InstructorInfo{}) {
		err = db.CreateTable(&pb.InstructorInfo{}).Error
		if err != nil {
			panic("Couldn't create ins table")
		}
	}

	if !db.HasTable(&pb.UserInfo{}) {
		err = db.CreateTable(&pb.UserInfo{}).Error
		if err != nil {
			panic("Couldn't create user table")
		}
	}

	if !db.HasTable(&pb.SessionInfo{}) {
		err = db.CreateTable(&pb.SessionInfo{}).Error
		if err != nil {
			panic("Couldn't create session table")
		}
	}

	if !db.HasTable(&pb.UserInstructorReview{}) {
		err = db.CreateTable(&pb.UserInstructorReview{}).Error
		if err != nil {
			panic("Couldn't create ins review table")
		}
	}

	if !db.HasTable(&pb.UserSessionReview{}) {
		err = db.CreateTable(&pb.UserSessionReview{}).Error
		if err != nil {
			panic("Couldn't create session review table")
		}
	}

	if !db.HasTable(&pb.CreditCard{}) {
		err = db.CreateTable(&pb.CreditCard{}).Error
		if err != nil {
			panic("Couldn't create cc table")
		}
	}

	if !db.HasTable(&pb.BankAcct{}) {
		err = db.CreateTable(&pb.BankAcct{}).Error
		if err != nil {
			panic("Couldn't create bank info table")
		}
	}

	log.Debug("Successfully opened  database and created tables")
	return nil
}

func main() {
	// open a file
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
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

	err = initDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to init DB")
		return
	}
	//err = pay.InitPayPlan()
	initGprcServer()
	db.Close()
}
