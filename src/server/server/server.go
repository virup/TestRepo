package main

import (
	"fmt"
	"net"
	"os"
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"

	"github.com/ajain1990/gorocksdb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

var RocksDBPath = "/libera/bin/rocksdb/"
var dbname = "mydb"
var rdb *gorocksdb.DB
var sessionsCF *gorocksdb.ColumnFamilyHandle
var usersCF *gorocksdb.ColumnFamilyHandle
var instructorsCF *gorocksdb.ColumnFamilyHandle
var wo *gorocksdb.WriteOptions
var ro *gorocksdb.ReadOptions

// XXX Cleanup pending
// XXX Reopening an existing DB doesn't work for now.
func initRocksDB() error {
	err := os.MkdirAll(RocksDBPath, 0750)
	if err != nil {
		log.Errorf("Can not create DB directory for '%s'; Error '%v'",
			dbname, err)
		return err
	}

	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetCreateIfMissingColumnFamilies(true)
	wo = gorocksdb.NewDefaultWriteOptions()
	ro = gorocksdb.NewDefaultReadOptions()

	/* Read and write options are required for reading and writing the keys */
	rdb, err = gorocksdb.OpenDb(opts, RocksDBPath+dbname)
	if err != nil {
		log.Errorf("Opening of rocks DB '%s' failed with error '%v'",
			dbname, err)
		return err
	}
	log.Debug("Successfully opened RocksDB's database", dbname)
	log.Debug("Successfully opened RocksDB's database %v", rdb)
	sessionsCF, err = rdb.CreateColumnFamily(opts, "sessions")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to create session CF")
		return err
	}
	usersCF, err = rdb.CreateColumnFamily(opts, "users")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to create users CF")
		return err
	}
	instructorsCF, err = rdb.CreateColumnFamily(opts, "instructors")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to create instructors CF")
		return err
	}
	return nil
}

func main() {
	// open a file
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
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

	err = initRocksDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to init rocks DB")
		return
	}
	//err = pay.InitPayPlan()
	initGprcServer()
}
