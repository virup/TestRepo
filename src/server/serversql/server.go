package main

import (
	"server/db"
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
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"google.golang.org/grpc/credentials"
	"strconv"
	"time"
	"errors"
)

const (
	port = ":8099"
)

var DATABASE_IP = os.Getenv("SF_DB_PORT_3306_TCP_ADDR")

const DATABASE_PORT = 3306
const DATABASE_NAME = "soulfitdb"
const DATABASE_USER = "root"
const DATABASE_PASSWORD = "password"

// XXX Consolidate user, instructor and session info object
// handling and DB handling through a common interface

// server is used to implement helloworld.GreeterServer.
type server struct{}

//  Send hello
func (s *server) GetStatus(ctx context.Context,
	in *pb.SFServerStatusReq) (*pb.SFServerStatusReply, error) {
	return &pb.SFServerStatusReply{Message: "Hello " + in.Name}, nil
}

func (s *server) CleanupAllDBs(ctx context.Context,
	in *pb.CleanupAllDBsReq) (*pb.CleanupAllDBsReply, error) {

	var err error
	var resp pb.CleanupAllDBsReply

	err = dbConn.DropTableIfExists(&pb.SessionInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete session table")
	}

	err = dbConn.DropTableIfExists(&pb.InstructorInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete instructor table")
	}

	err = dbConn.DropTableIfExists(&pb.UserInfo{}).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to delete user table")
	}

	log.Debugf("Cleaned up all tables")

	return &resp, nil
}

func (s *server) RecordEvent(ctx context.Context, in *pb.RecordEventReq) (*pb.RecordEventReply, error) {

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

	log.Debug("Registering server...")
	pb.RegisterSFServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Debug("Registered server...")
}

var dbConn *gorm.DB

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func createTableIfNotExists(tableType interface{}) {
	if !dbConn.HasTable(tableType) {
		err := dbConn.CreateTable(tableType).Error
		if err != nil {
			panic("Couldn't create table " + getType(tableType))
		}
	}
}

func getDatabaseConnectionString() string {
	return DATABASE_USER + ":" + DATABASE_PASSWORD +
		"@tcp(" + DATABASE_IP + ":" + strconv.Itoa(DATABASE_PORT) + ")/" +
		DATABASE_NAME + "?charset=utf8&parseTime=True&loc=Local"
}

func validate() {
	if DATABASE_IP == "" {
		fmt.Println("Specify the environment variable DATABASE_IP")
		os.Exit(1)
	}
}

func initDB() error {
	log.Debugf("Getting DB connections ...")

	validate()
	dbConnectionString := getDatabaseConnectionString()
	fmt.Println(dbConnectionString)

	err := errors.New("First Start")
	for err != nil {
		dbConn, err = gorm.Open("mysql", dbConnectionString)
		if err != nil {
			fmt.Println("Waiting for DB connection ...")
			log.Errorf("Opening of DB failed with error '%v'", err)
			time.Sleep(3 * time.Second)
		} else {
			fmt.Println("Connected!")
			log.Debug("Connected to the DB")
		}
	}


	dbConn.LogMode(true)
	log.Errorf("DB connection made successfully...")

	err = dbConn.DB().Ping()
	if err != nil {
		panic(err)
	}

	createTableIfNotExists(&db.InstructorInfo{})

	createTableIfNotExists(&db.UserInfo{})

	createTableIfNotExists(&db.SessionInfo{})

	createTableIfNotExists(&db.UserInstructorReview{})

	createTableIfNotExists(&db.UserSessionReview{})

	createTableIfNotExists(&db.CreditCard{})

	createTableIfNotExists(&db.BankAcct{})

	log.Debug("Successfully opened  database and created tables")
	return nil
}

func startLogging() (*os.File, bool) {
	// open a file
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil, false
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	return f, true
}

func main() {
	logfile, ok := startLogging()
	if ok {
		defer logfile.Close()
	}

	err := initDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed to init DB")
		return
	}

	initGprcServer()
	dbConn.Close()
}
