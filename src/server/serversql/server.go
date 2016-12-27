package main

import (
	"fmt"
	"net"
	"os"
	pb "server/rpcdefsql"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Debug("registering server...")
	pb.RegisterServerSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Debug("registered server...")

}

var db *gorm.DB
var InsTable *gorm.DB
var UserTable *gorm.DB
var SessionTable *gorm.DB

func initDB() error {
	//db, _ := gorm.Open("sqlite3", "./gorm.db")
	db, err := gorm.Open("sqlite3", "./gorm.db")
	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Errorf("Opening of DB failed with error '%v'",
			err)
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	InsTable = db.CreateTable(&pb.InstructorInfo{})
	if InsTable == nil {
		panic("Couldn't create ins table")
	}
	UserTable = db.CreateTable(&pb.UserInfo{})
	if UserTable == nil {
		panic("Couldn't create user table")
	}

	SessionTable = db.CreateTable(&pb.SessionInfo{})
	if SessionTable == nil {
		panic("Couldn't create session table")
	}

	log.Debug("Successfully opened  database and created tables")
	return nil
}

func main() {
	// open a file
	f, err := os.OpenFile("serversql.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
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
