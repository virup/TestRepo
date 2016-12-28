package main

import (
	pb "server/rpcdefsql"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

var lastSessionID int32 = 1

func getSessionID() int32 {

	ret := lastSessionID
	lastSessionID++
	return ret
}

// Given a sessionKey, return the SessionInfo
func getSessionFromDB(sKey int32) (error, *pb.SessionInfo) {
	var err error
	var s *pb.SessionInfo = new(pb.SessionInfo)
	err = SessionTable.First(s, sKey).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, s
	}
	log.WithFields(log.Fields{"sessionInfo": s, "key": sKey}).
		Debug("Read from DB")
	return err, s
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsReq) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	var err error
	var sList []pb.SessionInfo

	err = SessionTable.Find(&sList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		return &resp, err
	}

	log.Printf("\n")
	for i, _ := range sList {
		resp.SessionList = append(resp.SessionList, &sList[i])
	}

	for _, session := range resp.SessionList {
		err, insInfo := GetInstructorFromDB(session.InstructorID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor from session from DB")
			return &resp, err
		}
		resp.InstructorList = append(resp.InstructorList, insInfo)
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
	sKey = getSessionID()
	in.ID = sKey
	err = SessionTable.Save(&in).Error
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, 0
	}
	log.WithFields(log.Fields{"sessionInfo": in, "key": sKey}).
		Debug("Added to DB")
	return nil, sKey
}

func (ser *server) PostSession(ctx context.Context,
	in *pb.PostSessionReq) (*pb.PostSessionReply, error) {

	var err error
	var resp pb.PostSessionReply
	log.WithFields(log.Fields{"sessionInfo": in.Info}).
		Debug("Received post session request")

	err, resp.SessionKey = postSessionDB(*in.Info)
	if err != nil {
		log.WithFields(log.Fields{"session": in.Info, "error": err}).
			Error("Failed to write to DB")
		return &resp, err
	}
	log.WithFields(log.Fields{"session": in.Info}).
		Debug("Post session succeeded")
	return &resp, nil
}
