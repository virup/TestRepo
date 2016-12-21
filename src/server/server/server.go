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
	//err := db.All(&sList)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, nil
	}
	return err, sList
}

func getSessionFromDB(sKey string) (error, pb.Session) {
	var s pb.Session
	//err := db.One("ID", sKey, &s)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, s
	}
	return err, s
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsRequest) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	//err := db.All(&resp.Session)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		resp.ErrData = &pb.ErrorData{internalError, err.Error()}
		return &resp, nil
	}
	resp.ErrData = &pb.ErrorData{successError, successError}
	return &resp, nil
}

func postSessionDB(in *pb.SessionInfo) (error, string) {

	log.WithFields(log.Fields{"sessionInfo": in}).Debug("Adding to DB")
	s := new(pb.Session)
	s.Info = in
	s.ID = getRandomID()

	log.WithFields(log.Fields{"session": s}).Debug("Before Adding to DB")
	//err := db.Save(*s)
	if err != nil {
		log.WithFields(log.Fields{"session": s, "error": err}).Error("Failed" +
			" to write to DB")
		return err, ""
	}
	log.WithFields(log.Fields{"session": s}).Debug("Added to DB")
	return nil, s.ID
}

func (ser *server) PostSession(ctx context.Context,
	in *pb.PostSessionRequest) (*pb.PostSessionReply, error) {

	var resp pb.PostSessionReply
	log.Debug("Post Session grpc request")
	var s pb.Session
	s.Info = in.Info
	s.ID = getRandomID()
	err := db.Save(&s)
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

	var resp pb.EnrollInstructorResponse
	log.Debug("Enroll Instructor request")
	var i pb.Instructor
	i.Person = in.Instructor
	i.ID = getRandomID()
	err := db.Save(&i)
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

	var resp pb.EnrollUserResponse
	log.Debug("Enroll User request")
	var u pb.User
	u.Person = in.User
	u.ID = getRandomID()
	err := db.Save(&u)
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
	router.HandleFunc("/session/{sessionKey}", handleSession).Methods("GET",
		"DELETE", "POST")
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

func postSession(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var err error

	log.Debugf("postSession req %s", req.Body)
	info := new(pb.SessionInfo)
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
	res.WriteHeader(http.StatusCreated)
	result := sessionID + " created"
	fmt.Fprint(res, result)
}

func handleSession(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	sessionKey := vars["sessionKey"]

	switch req.Method {
	case "GET":
		//movie, ok := movies[sessionKey]
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
	case "DELETE":
		//delete(movies, sessionKey)
		res.WriteHeader(http.StatusNoContent)
	}
}

func initGprcServer() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := grpc.NewServer()
	log.Debug("registering server...")
	pb.RegisterServerSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Debug("registered server...")

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
	initRestServer()
}
