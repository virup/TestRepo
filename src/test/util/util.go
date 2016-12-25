package util

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	pb "server/rpcdef"
)

var sessionID = 1

func GetNewSession() pb.SessionInfo {

	var si pb.SessionInfo
	t := time.Now()
	si.SessionTime = t.String()
	si.SessionType = []pb.FitnessCategory{pb.FitnessCategory_YOGA, pb.FitnessCategory_FAST_YOGA}
	si.InstructorID = randSeq(10)
	si.SessionDesc = "my session" + strconv.Itoa(sessionID)
	return si
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func GetUrl() string {
	serv := os.Getenv("SERVERIP")
	if serv == "" {
		fmt.Printf("IP not set")
		os.Exit(-1)
	}
	return serv + ":8080"
}
func GetHttpUrl() string {
	serv := os.Getenv("SERVERIP")
	if serv == "" {
		fmt.Printf("SERVERIP not set")
		os.Exit(-1)
	}
	return "http://" + serv + ":8080"
}
