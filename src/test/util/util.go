package util

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"

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
	sessionID += 1
	return si
}

var certID = 1

func GetNewInstructor() pb.InstructorInfo {

	var ui pb.InstructorInfo
	ui.FirstName = randomdata.FirstName(randomdata.Female)
	ui.LastName = randomdata.LastName()
	ui.Email = randomdata.Email()
	ui.City = randomdata.City()
	ui.Certification = "FitnessCert" + strconv.Itoa(certID)
	certID += 1

	return ui
}

func GetNewUser() pb.UserInfo {

	var ui pb.UserInfo
	ui.FirstName = randomdata.FirstName(randomdata.Male)
	ui.LastName = randomdata.LastName()
	ui.Email = randomdata.Email()
	ui.City = randomdata.City()

	return ui
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
