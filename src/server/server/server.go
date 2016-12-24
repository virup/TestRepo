package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"pay"
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"

	"github.com/ajain1990/gorocksdb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port      = ":8080"
	soulFitDB = "SoulFitDB"
)

var lastUserUserID uint64

// server is used to implement helloworld.GreeterServer.
type server struct{}

//  Send hello
func (s *server) GetStatus(ctx context.Context,
	in *pb.ServerSvcStatusReq) (*pb.ServerSvcStatusReply, error) {
	return &pb.ServerSvcStatusReply{Message: "Hello " + in.Name}, nil
}

func getAllSessionFromDB() (error, []pb.SessionInfo) {
	var sList []pb.SessionInfo
	var err error

	log.Debug("Reading sessions from DB")
	it := rdb.NewIteratorCF(ro, sessionsCF)
	defer it.Close()

	//it.Seek([]byte("foo"))
	for ; it.Valid(); it.Next() {
		log.WithFields(log.Fields{"key": it.Key().Data(),
			"value": it.Value().Data()}).
			Debug("Reading key from DB")
	}

	if err := it.Err(); err != nil {
		return err, sList
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, nil
	}
	return err, sList
}

func getSessionFromDB(sKey string) (error, pb.SessionInfo) {
	var err error
	var buf []byte

	s := pb.SessionInfo{}

	v, err := rdb.GetCF(ro, sessionsCF, []byte(sKey))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, s
	}
	log.WithFields(log.Fields{"value": v}).Debug("Read" +
		"session value from DB")

	if v.Size() > 0 {
		buf = make([]byte, v.Size())
		copy(buf, v.Data())
		v.Free()
	} else {
		log.WithFields(log.Fields{"error": err}).Error("corrupted" +
			" session from DB")
	}
	err = proto.Unmarshal(buf, &s)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to unmarshal proto from DB")
	}
	log.WithFields(log.Fields{"sessionInfo": s, "key": sKey}).
		Debug("Read from DB")
	return err, s
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsReq) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	var err error
	//err, resp.sessionList := getAllSessionFromDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return &resp, err
	}
	return &resp, nil
}

func (s *server) GetSession(ctx context.Context,
	in *pb.GetSessionReq) (*pb.GetSessionReply, error) {

	var resp pb.GetSessionReply
	err, si := getSessionFromDB(in.SessionKey)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return &resp, err
	}

	log.WithFields(log.Fields{"session": s}).Debug("Get session success")
	resp.Info = &si
	return &resp, nil
}

func postSessionDB(in pb.SessionInfo) (err error, sessionKey string) {

	log.WithFields(log.Fields{"sessionInfo": in}).Debug("Adding to DB")
	sessionKey = getRandomID()

	byteBuf, err := proto.Marshal(&in)
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to convert to binary")
		return err, ""
	}

	log.Debugf("wo %#v sessionsCF %#v, sessionKey %#v byteBuf %#v",
		wo, sessionsCF, sessionKey, byteBuf)
	err = rdb.PutCF(wo, sessionsCF, []byte(sessionKey), byteBuf)
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, ""
	}
	log.WithFields(log.Fields{"sessionInfo": in, "key": sessionKey}).
		Debug("Added to DB")
	return nil, sessionKey
}

func (ser *server) PostSession(ctx context.Context,
	in *pb.PostSessionReq) (*pb.PostSessionReply, error) {

	var err error
	var resp pb.PostSessionReply
	log.WithFields(log.Fields{"sessionInfo": in.Info}).
		Debug("Received post session request")

	err, _ = postSessionDB(*in.Info)
	if err != nil {
		log.WithFields(log.Fields{"session": in.Info, "error": err}).
			Error("Failed to write to DB")
		return &resp, err
	}
	log.WithFields(log.Fields{"session": in.Info}).
		Debug("Post session succeeded")
	return &resp, nil
}

func (s *server) GetInstructors(ctx context.Context,
	in *pb.GetInstructorsReq) (*pb.GetInstructorsReply, error) {

	var resp pb.GetInstructorsReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get instructors from DB")
		return &resp, err
	}
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

func (s *server) Login(ctx context.Context,
	in *pb.LoginReq) (*pb.LoginReply, error) {

	var resp pb.LoginReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to login")
		return &resp, err
	}
	return &resp, nil
}

func (s *server) EnrollInstructor(ctx context.Context,
	in *pb.EnrollInstructorReq) (*pb.EnrollInstructorReply, error) {

	var err error
	var resp pb.EnrollInstructorReply
	log.Debug("Enroll Instructor request")
	//i.ID = getRandomID()
	//err := db.Save(&i)
	if err != nil {
		log.WithFields(log.Fields{"instructor": in.Instructor,
			"error": err}).Error("Failed to write to DB")
		return &resp, err
	}
	log.WithFields(log.Fields{"instructor": in.Instructor}).
		Debug("Added to DB")
	return &resp, nil
}

func getRandomID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *server) GetUsers(ctx context.Context,
	in *pb.GetUsersReq) (*pb.GetUsersReply, error) {

	var resp pb.GetUsersReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get users from DB")
		return &resp, err
	}
	return &resp, nil
}

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserReq) (*pb.EnrollUserReply, error) {

	var err error
	var resp pb.EnrollUserReply
	log.Debug("Enroll User request")
	//err := db.Save(&u)
	if err != nil {
		log.WithFields(log.Fields{"user": in.User, "error": err}).
			Error("Failed to write to DB for user")
		return &resp, err
	}
	log.WithFields(log.Fields{"user": in.User}).Debug("Added to DB")

	err, customerPayID := pay.CreatePayingCustomer(
		"mycustomer@gmail.com", "1234-xxxx-xxxx", "06", "19")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to set up customer payment")
		return &resp, err
	}
	log.WithFields(log.Fields{"customerPayID": customerPayID}).
		Debug("Got new customer payment ID")

	err = pay.StartSubscription(customerPayID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to start subscription")
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
