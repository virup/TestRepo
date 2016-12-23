package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	pb "server/rpcdef"
)

func GetNewSession() pb.SessionInfo {

	var si pb.SessionInfo
	t := time.Now()
	si.SessionTime = t.String()
	si.SessionType = "stype"
	si.InstructorID = randSeq(10)
	si.SessionDesc = "my session"
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
		fmt.Printf("IP not set")
		os.Exit(-1)
	}
	return "http://" + serv + ":8080"
}
