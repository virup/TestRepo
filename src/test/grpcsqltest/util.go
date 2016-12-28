package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"

	pb "server/rpcdefsql"
)

var sessionID = 1

var enrolledInstructorsID []int32

func RegisterEnrolledInstructorID(instructorID int32) {

	enrolledInstructorsID = append(enrolledInstructorsID, instructorID)
}

func getEnrolledInstructorID() (error, int32) {
	numIns := len(enrolledInstructorsID)
	if numIns == 0 {
		return errors.New("Instructors not registered"), 0
	}
	return nil, enrolledInstructorsID[rand.Intn(len(enrolledInstructorsID))]
}

func GetNewSession() (error, pb.SessionInfo) {

	var si pb.SessionInfo
	var err error
	t := time.Now()
	si.SessionTime = t.String()
	si.SessionType = pb.FitnessCategory_YOGA
	err, si.InstructorInfoID = getEnrolledInstructorID()
	if err != nil {
		return err, si
	}

	//si.TagList = []pb.SessionTag{pb.SessionTag_CALMING, pb.SessionTag_RELAXING}
	si.DifficultyLevel = pb.SessionDifficulty(rand.Intn(3)) //pb.SessionDifficulty_MODERATE
	si.SessionDesc = "my session" + strconv.Itoa(sessionID)
	sessionID += 1
	return nil, si
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
	return serv + ":8099"
}
func GetHttpUrl() string {
	serv := os.Getenv("SERVERIP")
	if serv == "" {
		fmt.Printf("SERVERIP not set")
		os.Exit(-1)
	}
	return "http://" + serv + ":8099"
}
