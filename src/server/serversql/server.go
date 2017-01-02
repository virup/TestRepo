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

	//if !db.HasTable(&pb.InstructorImage{}) {
	//	err = db.CreateTable(&pb.InstructorImage{}).Error
	//	if err != nil {
	//		panic("Couldn't create image table")
	//	}
	//}

	//if !db.HasTable(&pb.SessionVideo{}) {
	//	err = db.CreateTable(&pb.SessionVideo{}).Error
	//	if err != nil {
	//		panic("Couldn't create video table")
	//	}
	//}

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
