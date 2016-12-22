package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	pb "server/rpcdef"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	"github.com/ajain1990/gorocksdb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

const (
	port          = ":50051"
	soulFitDB     = "SoulFitDB"
	internalError = "internalError"
	successError  = "success"
)

var lastUserUserID uint64

// server is used to implement helloworld.GreeterServer.
type server struct{}

//  Send hello
func (s *server) GetStatus(ctx context.Context,
	in *pb.ServerSvcStatusRequest) (*pb.ServerSvcStatusResponse, error) {
	return &pb.ServerSvcStatusResponse{Message: "Hello " + in.Name}, nil
}

func getAllSessionFromDB() (error, []pb.Session) {
	var sList []pb.Session
	var err error
	//err := db.All(&sList)
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
	} else {

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
	}
	return err, s
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsRequest) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, nil
	}
	resp.ErrData = &pb.ErrorData{successError, successError}
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
	log.WithFields(log.Fields{"sessionInfo": in, "key": sessionKey}).Debug("Added to DB")
	return nil, sessionKey
}

func (ser *server) PostSession(ctx context.Context,
	in *pb.PostSessionRequest) (*pb.PostSessionReply, error) {

	var err error
	var resp pb.PostSessionReply
	log.Debug("Post Session grpc request")
	var s pb.Session
	s.Info = in.Info
	s.ID = getRandomID()
	//err := db.Save(&s)
	if err != nil {
		log.WithFields(log.Fields{"session": s, "error": err}).Error("Failed" +
			" to write to DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"session": s}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func (s *server) EnrollInstructor(ctx context.Context,
	in *pb.EnrollInstructorRequest) (*pb.EnrollInstructorResponse, error) {

	var err error
	var resp pb.EnrollInstructorResponse
	log.Debug("Enroll Instructor request")
	var i pb.Instructor
	i.Person = in.Instructor
	i.ID = getRandomID()
	//err := db.Save(&i)
	if err != nil {
		log.WithFields(log.Fields{"instructor": i, "error": err}).Error("Failed" +
			" to write to DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"instructor": i}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func getRandomID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserRequest) (*pb.EnrollUserResponse, error) {

	var err error
	var resp pb.EnrollUserResponse
	log.Debug("Enroll User request")
	var u pb.User
	u.Person = in.User
	u.ID = getRandomID()
	//err := db.Save(&u)
	if err != nil {
		log.WithFields(log.Fields{"user": u, "error": err}).Error("Failed" +
			" to write to DB for user")

		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, err
	}
	log.WithFields(log.Fields{"user": u}).Debug("Added to DB")
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func initRestServer() {
	router := mux.NewRouter()
	router.HandleFunc("/getstatus", getStatus).Methods("GET")
	router.HandleFunc("/getsessions", getSessions).Methods("GET")
	router.HandleFunc("/getsession/{sessionKey}", getSession).Methods("GET")
	router.HandleFunc("/deletesession/{sessionKey}", deleteSession).Methods("DELETE")
	router.HandleFunc("/postsession", postSession).Methods("POST")
	http.ListenAndServe(":8080", router)

	log.Debug("rest server running...")
}

func getStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	fmt.Fprint(res, "running from server!")
}

func getSessions(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	log.Debug("getSessions req")
	err, sessionList := getAllSessionFromDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(sessionList)

	log.WithFields(log.Fields{"json": outgoingJSON}).Debug("Sending" +
		" sessions from DB")
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, string(outgoingJSON))
}

type PostSessionResponse struct {
	sessionID string `json:"sessionID,omitempty"`
}

func postSession(res http.ResponseWriter, req *http.Request) {
	var err error

	log.Debugf("postSession req %s", req.Body)
	var info pb.SessionInfo
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&info)
	if error != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to decode json")
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	err, sessionID := postSessionDB(info)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to post session to DB")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	log.WithFields(log.Fields{"sessionID": sessionID}).Debug("Post response")

	var sid PostSessionResponse
	sid.sessionID = sessionID
	json.NewEncoder(res).Encode(sid)
}

func getSession(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	sessionKey := vars["sessionKey"]

	log.WithFields(log.Fields{"sessionKey": sessionKey}).Debug("getsession request")
	err, session := getSessionFromDB(sessionKey)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, string("Session not found"))
	}
	outgoingJSON, error := json.Marshal(session)
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func deleteSession(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	sessionKey := vars["sessionKey"]

	log.WithFields(log.Fields{"sessionKey": sessionKey}).Debug("deletesession request")
	//delete(movies, sessionKey)
	res.WriteHeader(http.StatusNoContent)
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
	f, err := os.OpenFile("testlogrus.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	//initGprcServer()
	fmt.Printf("init rest")
	err = initRocksDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to init rocks DB")
		return
	}
	initRestServer()
}
