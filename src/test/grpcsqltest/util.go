package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	randomdata "github.com/Pallinder/go-randomdata"

	pb "server/rpcdefsql"
)

var sessionID = 1

var enrolledInstructorsID []int32
var enrolledUsersID []int32

var insMap = make(map[int32]*pb.InstructorInfo)
var userMap = make(map[int32]*pb.UserInfo)

func RegisterEnrolledInstructor(instructorID int32,
	iInfo *pb.InstructorInfo) {

	enrolledInstructorsID = append(enrolledInstructorsID, instructorID)
	insMap[instructorID] = iInfo
}

func RegisterEnrolledUser(userID int32,
	uInfo *pb.UserInfo) {

	enrolledUsersID = append(enrolledUsersID, userID)
	userMap[userID] = uInfo
}

func getEnrolledUser(uid int32) *pb.UserInfo {
	uInfo, ok := userMap[uid]
	if !ok {
		panic("invalid user")
	}
	return uInfo
}

func getEnrolledIns(insID int32) *pb.InstructorInfo {
	iInfo, ok := insMap[insID]
	if !ok {
		panic("invalid ins")
	}
	return iInfo
}

func getEnrolledInstructorID() (error, int32) {
	numIns := len(enrolledInstructorsID)
	if numIns == 0 {
		return errors.New("Instructors not registered"), 0
	}
	return nil, enrolledInstructorsID[rand.Intn(len(enrolledInstructorsID))]
}

func getEnrolledUsersID() (error, int32) {
	numIns := len(enrolledUsersID)
	if numIns == 0 {
		return errors.New("Users not registered"), 0
	}
	return nil, enrolledUsersID[rand.Intn(len(enrolledUsersID))]
}

var nextSessionId = 0

func GetNewSession() (error, pb.SessionInfo) {
	if nextSessionId == len(allSessions) {
		panic("return session")
	}
	retId := nextSessionId
	nextSessionId++
	return nil, allSessions[retId]

}

/*
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
*/

var certID = 1

var allSessions = []pb.SessionInfo{
	{
		SessionDesc:      "Alright, lets be real. The winter months are coming to an end and it's time to transition from bears to gazelles! Whether you are looking to get comfortable in your swimwear or just more agile to move with the spring breeze let's start with this 20 min practice that moves us swiftly but mindfully. This 20 minute vinyasa practice is designed to help you build strength and endurance - mindfully. Yoga Tone invites strong breath to help tone the body! Invite your mind and body to start working for you instead of against you. Get strong with regular practice and comfortable in your beautiful body. Stay present. Let's move!",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=4O-b24WKdYA",
		InstructorName:   "Adriene Smith",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 1,
	},
	{
		SessionDesc:      "This PURE, Cardio only workout will be intense right out of the gate! We are going to do 3 Tabata Intervals of easy to follow, effective but modifiable cardio exercises. Let's get our sweat ON!",
		DifficultyLevel:  pb.SessionDifficulty_DIFFICULT,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=TQPedXNPcBQ",
		InstructorName:   "Shelly Dose",
		SessionType:      pb.FitnessCategory_CARDIO,
		InstructorInfoID: 2,
	},
	{
		SessionDesc:      "This is a thirty minute power yoga class. There are options for modifications throughout, so it's okay for people new to yoga (though I wouldn't say it's appropriate for absolute beginners). This is a vinyasa yoga class. Vinyasa means 'movemen't, so we will be constantly flowing and moving through the practice. If you're looking for a slower style with long holds I would recommend a hatha yoga practice. As always, honor your body, and work within a pain-free range.",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=qy_oIKf1ByM",
		InstructorName:   "Carolina B",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 3,
	},
	{
		SessionDesc:      "A 30 minutes zumba dance workout that you'll be able to do at home. Try this routine three times a week to keep in good shape and help you lose weight. Have fun dancing while you work out your whole body to the best zumba beats in this ultimate zumba tutorial!",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=qAJ6EQtGZ28",
		InstructorName:   "Monica B",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 4,
	},
	{
		SessionDesc: `This is a powerful guided self hypnosis trance experience designed to allow you to sweep away your own subconscious negativity and negative blocks. Clear out all of your subconscious or unconscious negative thoughts, old habits, and emotional baggage with your own positive mind control. 

	With regular self hypnosis, you can truly allow your powerful, positive thinking self to emerge for your best present and brighter future.`,
		DifficultyLevel:  pb.SessionDifficulty_EASY,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=FiPDV9L5qpQ",
		InstructorName:   "Michael Sealy",
		SessionType:      pb.FitnessCategory_MEDITATION,
		InstructorInfoID: 5,
	},
}

var allIns = []pb.InstructorInfo{
	{
		FirstName:   "Adriene",
		LastName:    "Smith",
		FitnessType: pb.FitnessCategory_YOGA,
		Desc:        "Welcome all levels, all bodies, all genders, all souls! Find a practice that suits your mood or start a journey toward healing. Work up a sweat, or calm and relieve a tired mind and body. Create space. Tone and trim. Cultivate self love. Make time for you. Go deeper, have fun. Connect. Fall off the horse and then get back on. Reconnect. Do your best, be authentic and FIND WHAT FEELS GOOD.I got your back and this community rocks. Jump on in! You don't even have to leave your house.",
	},

	{
		FirstName:   "Shelley",
		LastName:    "Dose",
		FitnessType: pb.FitnessCategory_CARDIO,
		Desc:        "Hi, my name is Shelly Dose and I want to move with you! My goal is to move you in as many ways as possible with athletic, low impact, high impact, HIIT, full body workouts, sculpt... you name it, I can teach it. I am a Certified Group Fitness Instructor with a long resume in the fitness industry. I have owned and operated a successful Outdoor Bootcamp Business, I teach Group Exercise/Fitness for Lifetime Fitness and much more. I enjoy teaching in the classroom immensely but my favorite place to sweat is in the comfort of my home.",
	},

	{
		FirstName:   "Candace",
		LastName:    "Rose",
		FitnessType: pb.FitnessCategory_YOGA,
		Desc:        "Candace is an international yoga instructor and the writer behind the popular yoga lifestyle blog http://www.yogabycandace.com",
	},

	{
		FirstName:   "Carolina",
		LastName:    "B",
		FitnessType: pb.FitnessCategory_DANCE,
		Desc:        "I feel blessed for being able to do what I love. Dance is my passion, and it has brought so many good things to my life, such as health, friends, and many good times. I hope I can share that with you thru my dancing and my favorite Dance Fitness and Zumba® routines. Thank you so much for the support and for watching these videos I make with so much love. Mwaaahhhh!",
	},

	{
		FirstName:   "Michael",
		LastName:    "Sealy",
		FitnessType: pb.FitnessCategory_MEDITATION,
		Desc: `Hypnosis - Hypnotherapy - Guided Meditation - Sleep Relaxation

		Hi, my name is Michael and welcome to my channel, where I hope you can stop by to relax, listen in, and see for yourself the power of positive hypnosis."+

		"Hypnosis is a completely natural state of often deeply felt relaxation and focused attention, where positive suggestions can be more easily accepted by our subconscious minds. Imagine a fantastic and tranquil state of daydreaming, and that is very close to hypnosis! 

		Peace & Enjoy`,
	},
}

var InsImages = []string{
	"/libera/bin/sessioninsimages/yoga1.jpg",
	"/libera/bin/sessioninsimages/hiit1.jpeg",
	"/libera/bin/sessioninsimages/yoga2.jpeg",
	"/libera/bin/sessioninsimages/zumba.jpg",
	"/libera/bin/sessioninsimages/meditation1.jpg",
}

//func GetNewInstructor() pb.InstructorInfo {

//}

var nextIns = 0

func GetNewInstructor() pb.InstructorInfo {

	if false {
		var ui pb.InstructorInfo
		ui.FirstName = randomdata.FirstName(randomdata.Female)
		ui.LastName = randomdata.LastName()
		ui.Email = randomdata.Email()
		ui.City = randomdata.City()
		ui.Certification = "FitnessCert" + strconv.Itoa(certID)
		certID += 1
		return ui
	}
	if nextIns == len(allIns) {
		panic("return ins")
	}
	retId := nextIns
	nextIns++
	return allIns[retId]
}

func GetNewUser() pb.UserInfo {

	var ui pb.UserInfo
	ui.FirstName = randomdata.FirstName(randomdata.Male)
	ui.LastName = randomdata.LastName()
	ui.Email = randomdata.Email()
	ui.PassWord = randomdata.LastName()
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
