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

func GetNumInstructorImages() int {
	imglen := len(InsImages)
	if imglen != GetNumInstructors() {
		panic("Invalid number of ins images")
	}
	return imglen
}

func GetNumInstructors() int {
	return len(allIns)
}

func GetNumSessions() int {
	return len(allSessions)
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
		SessionDesc:      "Celebrity trainer JJ Dancer takes dance workouts to a new level. Get ready to pop, kick, and burn calories.",
		DifficultyLevel:  pb.SessionDifficulty_DIFFICULT,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=EV9GYSWij84",
		InstructorName:   "Popsugar fitness",
		SessionType:      pb.FitnessCategory_CARDIO,
		InstructorInfoID: 7,
	},
	{
		SessionDesc:      "This 30-minute workout will leave you looking and feeling fit, sexy, and strong. Grab a pair of free weights and get ready to work.",
		DifficultyLevel:  pb.SessionDifficulty_DIFFICULT,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=NXoy9ZVYy9I&list=PLI37FJmOtrj20cgTm5hcfZ-0H1PLHCQlj",
		InstructorName:   "Popsugar fitness",
		SessionType:      pb.FitnessCategory_CARDIO,
		InstructorInfoID: 7,
	},

	{
		SessionDesc:      "This is a thirty minute power yoga class. There are options for modifications throughout, so it's okay for people new to yoga (though I wouldn't say it's appropriate for absolute beginners). This is a vinyasa yoga class. Vinyasa means 'movemen't, so we will be constantly flowing and moving through the practice. If you're looking for a slower style with long holds I would recommend a hatha yoga practice. As always, honor your body, and work within a pain-free range.",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=qy_oIKf1ByM",
		InstructorName:   "Candace Rose",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 3,
	},
	{
		SessionDesc:      "A 30 minutes zumba dance workout that you'll be able to do at home. Try this routine three times a week to keep in good shape and help you lose weight. Have fun dancing while you work out your whole body to the best zumba beats in this ultimate zumba tutorial!",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=qAJ6EQtGZ28",
		InstructorName:   "Carolina B",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 4,
	},
	{
		SessionDesc:      "At Doonya, our goal is to take your favorite Bollywood songs and make them a SUPER fun workout! Enjoy a little sweat to Udi Udi Jaye from Raees",
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=VxowNt5T7VM",
		InstructorName:   "Doonya",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 6,
	},
	{
		SessionDesc:      `Enjoy this new dance fitness routine to "Cheap Thrills" by Sia ft. Sean Paul.`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=Tj3fuAw0iRA",
		InstructorName:   "Jessica",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 8,
	},
	{
		SessionDesc:      `Get ready for squats in this dance fitness routine to "Closer" by The Chainsmokers featuring Halsey.`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=7ZHkBnWwCj8",
		InstructorName:   "Jessica",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 8,
	},

	{
		SessionDesc:      `The Bollywood Workout is designed to take your favorite Bollywood songs and allow you to learn choreography that will get you in shape. This routine to "The Humma Song" from Ok Jaanu is an ab routine that will strengthen your core!`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=Lhyh1ggc7CA",
		InstructorName:   "Doonya",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 6,
	},

	{
		SessionDesc: `
		Do you need a relaxing timeout from your busy work day? A short break away to recharge your batteries and refocus your concentration?

		Take a deeply relaxing virtual vacation right now with the imagination of your own subconscious mind as you learn to let go and release all workplace stress and anxiety. 

		By listening along to this guided hypnosis meditation you will be able to create your very own positive daydream experience -- in only 15 minutes -- for a refreshing mental break during your busy workday ... or for a powerfully rejuvenating mind focus at any time. 

		See, feel, and hear yourself living a more positive, relaxed, focused, productive and stress free professional and personal life. `,
		DifficultyLevel:  pb.SessionDifficulty_EASY,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=pZHLOU9cpk8",
		InstructorName:   "Michael Sealy",
		SessionType:      pb.FitnessCategory_MEDITATION,
		InstructorInfoID: 5,
	},

	{
		SessionDesc: `
Welcome to this beginner's guided meditation which uses mindfulness meditation techniques to help you positively enhance your own health and well being in just a few short minutes per day. 

This meditation is designed to teach you how to easily meditate with a straight forward practice of sitting mindfully with stillness and gently focused inner observation. 

This meditation experience may be repeated as often as you choose to reinforce your positive intentions, calmness, and mindful detachment from over-thinking or excessive emotional reactivity.

By your own inner direction may you continue to find your greatest waking potential.`,
		DifficultyLevel:  pb.SessionDifficulty_EASY,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=F0jedwTzIJg",
		InstructorName:   "Michael Sealy",
		SessionType:      pb.FitnessCategory_MEDITATION,
		InstructorInfoID: 5,
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

	{
		SessionDesc: `This is a powerful guided self hypnosis trance experience designed to allow you to sweep away your own subconscious negativity and negative blocks. Clear out all of your subconscious or unconscious negative thoughts, old habits, and emotional baggage with your own positive mind control. 

	With regular self hypnosis, you can truly allow your powerful, positive thinking self to emerge for your best present and brighter future.`,
		DifficultyLevel:  pb.SessionDifficulty_EASY,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=FiPDV9L5qpQ",
		InstructorName:   "Michael Sealy",
		SessionType:      pb.FitnessCategory_MEDITATION,
		InstructorInfoID: 5,
	},

	{
		SessionDesc: `Hone in on deep listening, challenge yourself and uncover the wisdom within.

		Renew the relationship with your gut! Your inner instincts guide the way toward a strong, fierce and conscious practice.

		Make self love cool again!`,
		DifficultyLevel:  pb.SessionDifficulty_EASY,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=_h6wIuyVcBU",
		InstructorName:   "Adriene Smith",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 1,
	},
	{
		SessionDesc:      `Yoga For Strength! Join Adriene for a 40 Minute Vinyasa Flow Yoga. This practice cultivates heat, trims, tones, builds strength and flexibility. Adriene works with a strong foundation and encourages you to integrate a long lasting breath practice in your Vinyasa flow yoga. With full body awareness and strong focus on alignment this practice is swift but offers variations for you to try as you build your practice. Open the hips, the shoulders and tap into your core strength. This vinyasa yoga practice tones the legs and the arms while offering a strong foundation to protect the joints. Be mindful and meet your edge! Return to this practice to experience your growth and deepen your practice. "The journey is the reward." Practice.`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=9kOCY0KNByw",
		InstructorName:   "Adriene Smith",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 1,
	},
	{
		SessionDesc:      `Find freedom within the form. Today we explore fluidity, tap into the lower core and awaken creative energy that serves. Let's play!`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=yTdQCuJwsaA",
		InstructorName:   "Adriene Smith",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 1,
	},
	{
		SessionDesc:      `In this vinyasa flow yoga class, we begin with a centering to help you connect with your breath. I teach a quick guide for how to do the ujjayi breath. From here we move into cat cow poses, then move into downward facing dog, then go into warrior 1 and do a little chest and front body opening. As we continue with the vinyasa, we do some back strengthening, and then slowly come into a savasana to rest and let the body absorb all the good things we just did.`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=hmaI2bg5upU",
		InstructorName:   "Candace Rose",
		SessionType:      pb.FitnessCategory_YOGA,
		InstructorInfoID: 3,
	},
	{
		SessionDesc: `If you're familiar with Bollywood culture, then we don't have to tell you that the music and dancing in these films are completely infectious. Bombay Jam takes the same energizing spirit of Bollywood, but gives it a fitness twist, creating a full-body cardio workout that keeps you moving and smiling the entire way though. Get ready to jump, squat, swivel, and pivot as Bombay Jam Master Trainer Janani Chalaka leads you through two different routines that will definitely get your heart pumping.`,

		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=pzob_Br-IWQ",
		InstructorName:   "BombayJam",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 9,
	},
	{
		SessionDesc: `Are you curios to see what Bombay Jam® is all about? Take a look!
		`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=I2Xz_9UTtRA",
		InstructorName:   "BombayJam",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 9,
	},
	{
		SessionDesc: `Watch our cardio dance routine to the party number, Saturday, Saturday! Have a blast while you melt away the calories with Bombay Jam Bollywood fitness classes!
		`,
		DifficultyLevel:  pb.SessionDifficulty_MODERATE,
		PreviewVideoUrl:  "https://www.youtube.com/watch?v=NNFtD3nKqJU",
		InstructorName:   "BombayJam",
		SessionType:      pb.FitnessCategory_DANCE,
		InstructorInfoID: 9,
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
		Desc:        "Candace is an international yoga instructor and the writer behind the popular yoga lifestyle blog http://www.yogabycandace.com. These videos range from 15 minutes to 60 minutes, and cover a wide variety of themes. All video content published after June 13, 2014 is available for purchase for those of you who travel often or just want the content on your device when you don't have internet connection. Purchase at http://www.yogabycandace.com/dvds",
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

	{
		FirstName:   "Doonya",
		LastName:    "",
		FitnessType: pb.FitnessCategory_DANCE,
		Desc: `
Celebrate your body, mind, and one of the most festive cultures of the world through Doonya*!

At Doonya, you'll follow along to cardio and conditioning intervals of varying intensity, each activating major muscle groups of the abs, arms and legs. The energy and expressions of Bollywood-inspired music and dance will keep you smiling as you burn up to 800 calories while learning dance and fitness fundamentals.

By using your own resistance and muscle control, you'll leave with a stronger core, lengthened limbs and increased stamina to keep you invigorated for the rest of your day. Doonya is your happy workout.`,
	},

	{
		FirstName:   "Popsugar fitness",
		LastName:    "",
		FitnessType: pb.FitnessCategory_CARDIO,
		Desc: `
By using your own resistance and muscle control, you'll leave with a stronger core, lengthened limbs and increased stamina to keep you invigorated for the rest of your day. Doonya is your happy workout.`,
	},

	{
		FirstName:   "Jessica",
		LastName:    "",
		FitnessType: pb.FitnessCategory_DANCE,
		Desc: `
		Hey Ya'll! my name is Jessica! I have been teaching dance classes for over 5 years now. My goal is to help inspire and motivate people to make positive changes in their lives. Everyone deserves to feel beautiful and live a healthy and active lifestyle. My community is so supportive and while we may be a little crazy, we make working out TONS OF FUN! 

		You can visit my website https://www.dancefitnesswithjessica.com for instructional dance fitness workout programs, available on DVD and MORE!`,
	},

	{
		FirstName:   "BombayJam",
		LastName:    "",
		FitnessType: pb.FitnessCategory_DANCE,
		Desc: `
		Leave your inhibitions behind and join our action-packed Bollywood dance fitness program. Sign up for a class in your area with Bombay Jam today!`,
	},
}

var InsImages = []string{
	"/libera/bin/sessioninsimages/adriene.jpg",
	"/libera/bin/sessioninsimages/shelleydose.jpeg",
	"/libera/bin/sessioninsimages/candace.jpeg",
	"/libera/bin/sessioninsimages/carolina.jpg",
	"/libera/bin/sessioninsimages/michael.jpg",
	"/libera/bin/sessioninsimages/dancedoonya.jpg",
	"/libera/bin/sessioninsimages/popsugar.jpg",
	"/libera/bin/sessioninsimages/jessica.jpg",
	"/libera/bin/sessioninsimages/bombayjam.jpg",
}

var nextIns = 0

func GetNewBankAcct(insID int) pb.BankAcct {
	var acct pb.BankAcct
	acct.AcctNum = strconv.Itoa(insID)
	acct.InstructorID = int32(insID)
	acct.BankName = "Bank of America"
	return acct
}

func GetNewCC(userID int) pb.CreditCard {
	var cc pb.CreditCard
	cc.Number = "1111-2222-3333-4444"
	cc.CCV = "123"
	cc.UserID = int32(userID)
	return cc
}

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
func getUrl() string {
	serv := os.Getenv("SERVERIP")
	if serv == "" {
		fmt.Printf("env var $SERVERIP not set\n")
		os.Exit(-1)
	}
	return serv
}

func getAddressAndPort() string {
	return getUrl() + ":8099"
}
