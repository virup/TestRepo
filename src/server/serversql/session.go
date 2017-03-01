package main

import (
	pb "server/rpcdefsql"

	log "github.com/Sirupsen/logrus"

	twilio "github.com/xaviiic/twilioGo"
	"golang.org/x/net/context"
)

const (
	TWILIO_ACCOUNT_ID = "ACcb259fbd219b08efc012786fb4e3fae9"
	TWILLIO_KEY_ID    = "a3a2315244c7080d373422f8801eacd1"
)

// getSessionFromDB - Given a sessionKey, return the SessionInfo
func getSessionFromDB(sKey int32) (error, *pb.SessionInfo) {
	var s pb.SessionInfo
	err := dbConn.First(&s, sKey).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, &s
	}
	log.WithFields(log.Fields{"sessionInfo": s, "key": sKey}).
		Debug("Read from DB")
	return err, &s
}

// GetSessionsForInstructor - Find the sessions created by the specified instructor
func (s *server) GetSessionsForInstructor(ctx context.Context,
	in *pb.GetSessionsForInstructorReq) (*pb.GetSessionsForInstructorReply,
	error) {

	var resp pb.GetSessionsForInstructorReply
	var sList []pb.SessionInfo


	err := dbConn.
		Where(pb.SessionInfo{InstructorInfoID: in.InstructorInfoID}).
		Find(&sList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		return &resp, err
	}

	for i, _ := range sList {
		resp.SessionList = append(resp.SessionList, &sList[i])
	}

	log.WithFields(log.Fields{"sessionFitList": resp.SessionList}).Debug("Returning instructor sessionlist success")
	return &resp, err
}

// GetSessionsForFitnessType - Returns all sessions which are of the particular session type.
func (s *server) GetSessionsForFitnessType(ctx context.Context,
	in *pb.GetSessionsForFitnessReq) (*pb.GetSessionsForFitnessReply,
	error) {

	var resp pb.GetSessionsForFitnessReply
	var err error
	var sList []pb.SessionInfo

	err = dbConn.
		Where(pb.SessionInfo{SessionType: in.FitCategory}).
		Find(&sList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		return &resp, err
	}
	for i, _ := range sList {
		resp.SessionList = append(resp.SessionList, &sList[i])
	}

	log.WithFields(log.Fields{"sessionFitList": resp.SessionList}).Debug("Returning fitness sessionlist success")
	return &resp, err
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsReq) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	var err error
	var sList []pb.SessionInfo

	//err = SessionTable.Find(&sList).Error
	err = dbConn.Find(&sList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		return &resp, err
	}

	log.Printf("\n")
	insKeys := make(map[int32]int32)
	for i, _ := range sList {
		resp.SessionList = append(resp.SessionList, &sList[i])
		insID := sList[i].InstructorInfoID
		insKeys[insID] = 0
	}

	// Get non-duplicate instructors who are offering sessions
	var insKeySlice []int32
	for k, _ := range insKeys {
		insKeySlice = append(insKeySlice, k)
	}

	log.WithFields(log.Fields{"insKeySlice": insKeySlice}).Debug("Query sessions with ins slice")

	var iList []pb.InstructorInfo
	//err = InsTable.Where(insKeySlice).Find(&iList).Error
	err = dbConn.Where(insKeySlice).Find(&iList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get instructor rows from DB")
		return &resp, err
	}
	for i, _ := range iList {
		resp.InstructorList = append(resp.InstructorList, &iList[i])
	}

	log.WithFields(log.Fields{"sessionList": resp.SessionList}).Debug("Returning sessionlist success")
	log.WithFields(log.Fields{"instructorList": resp.InstructorList}).Debug("Returning instructorList success")

	return &resp, nil
}

func (s *server) GetSession(ctx context.Context,
	in *pb.GetSessionReq) (*pb.GetSessionReply, error) {

	var resp pb.GetSessionReply
	err, sessionInfo := getSessionFromDB(in.SessionKey)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return &resp, err
	}

	log.WithFields(log.Fields{"session": sessionInfo}).Debug("Get session success")
	resp.Info = sessionInfo
	return &resp, nil
}

func postSessionDB(in pb.SessionInfo) (err error, sKey int32) {

	log.WithFields(log.Fields{"sessionInfo": in}).Debug("Adding to DB")
	//sKey = getSessionID()
	//in.ID = sKey
	//err = SessionTable.Save(&in).Error
	err = dbConn.Save(&in).Error
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, 0
	}
	sKey = in.ID
	log.WithFields(log.Fields{"sessionInfo": in, "key": sKey}).
		Debug("Added to DB")
	return nil, sKey
}

func (ser *server) PostSessionPreviewVideo(ctx context.Context,
	in *pb.PostSessionPreviewVideoReq) (*pb.PostSessionPreviewVideoReply, error) {

	var err error
	var resp pb.PostSessionPreviewVideoReply
	var session pb.SessionInfo

	log.WithFields(log.Fields{"previewVideoInfo": in.VidUrl}).
		Debug("Received post session preview videorequest")

	err = dbConn.First(&session, in.SessionID).
		Update(pb.SessionInfo{PreviewVideoUrl: in.VidUrl}).Error
	if err != nil {
		log.WithFields(log.Fields{"instructorImage": in,
			"error": err}).Error("Failed to update image ID to ins DB")
		return &resp, err
	}

	log.WithFields(log.Fields{"session": in.VidUrl}).
		Debug("Post session preview video succeeded")
	return &resp, nil
}

func (ser *server) PostSession(ctx context.Context, in *pb.PostSessionReq) (*pb.PostSessionReply, error) {

	log.WithFields(log.Fields{"sessionInfo": in.Info}).Debug("Received post session request")

	err, sessionKey := postSessionDB(*in.Info)
	if err != nil {
		log.WithFields(log.Fields{"session": in.Info, "error": err}).
			Error("Failed to write to DB")
		return nil, err
	}
	log.WithFields(log.Fields{"session": in.Info}).
		Debug("Post session succeeded")

	return &pb.PostSessionReply{SessionKey: sessionKey}, nil
}

func (s *server) PostSessionReview(ctx context.Context,
	in *pb.PostSessionReviewReq) (*pb.PostSessionReviewReply,
	error) {

	err := dbConn.Save(&in.Review).Error
	if err != nil {
		log.WithFields(log.Fields{"sessionReview": in, "error": err}).
			Error("Failed to write review to DB")
		return nil, err
	}
	log.WithFields(log.Fields{"sessionReview": in}).Debug("Added to DB")

	return &pb.PostSessionReviewReply{ReviewID: in.Review.ID}, nil
}

func (s *server) GetTwilioJwtToken(ctx context.Context, in *pb.TwilioJwtReq) (*pb.TwilioJwtReply, error) {

	secret := in.Secret
	identity := in.Identity
	// first create token with twilio api configurations
	token := twilio.NewAccessToken(TWILIO_ACCOUNT_ID, TWILLIO_KEY_ID, secret)
	// setup token identity
	token.SetIdentity(identity)

	// grant token access to progammable video API
	configurationProfileID := "profile-sid"
	grant := twilio.NewConversationGrant(configurationProfileID)
	token.AddGrant(grant)

	jwt, err := token.ToJWT()
	if err != nil {
		return nil, err
	}

	return &pb.TwilioJwtReply{string(jwt)}, nil
}
