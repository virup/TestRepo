package main

import (
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Given a sessionKey, return the SessionInfo
func getSessionFromDB(sKey string) (error, *pb.SessionInfo) {
	var err error
	var buf []byte

	var s *pb.SessionInfo = new(pb.SessionInfo)

	v, err := rdb.GetCF(ro, sessionsCF, []byte(sKey))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, s
	}
	log.WithFields(log.Fields{"value": v}).Debug("Read" +
		"session value from DB")

	if v.Size() > 0 {
		buf = make([]byte, v.Size())
		copy(buf, v.Data())
		v.Free()
	} else {
		log.WithFields(log.Fields{"error": err}).Error("invalid key/" +
			"session from DB")
	}
	err = proto.Unmarshal(buf, s)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to unmarshal proto from DB")
		return err, nil
	}
	log.WithFields(log.Fields{"sessionInfo": s, "key": sKey}).
		Debug("Read from DB")
	return err, s
}

func (s *server) GetSessions(ctx context.Context,
	in *pb.GetSessionsReq) (*pb.GetSessionsReply, error) {

	var resp pb.GetSessionsReply
	var err error
	err, resp.SessionList = getAllSessionFromDB()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get sessions from DB")
		return &resp, err
	}
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

func postSessionDB(in pb.SessionInfo) (err error, sessionKey string) {

	log.WithFields(log.Fields{"sessionInfo": in}).Debug("Adding to DB")
	sessionKey = GetRandomID()

	byteBuf, err := proto.Marshal(&in)
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to convert to binary")
		return err, ""
	}

	err = rdb.PutCF(wo, sessionsCF, []byte(sessionKey), byteBuf)
	if err != nil {
		log.WithFields(log.Fields{"sessionInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, ""
	}
	log.WithFields(log.Fields{"sessionInfo": in, "key": sessionKey}).
		Debug("Added to DB")
	return nil, sessionKey
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

func getAllSessionFromDB() (error, []*pb.SessionInfo) {
	var sList []*pb.SessionInfo
	var err error

	log.Debug("Reading all sessions from DB")
	it := rdb.NewIteratorCF(ro, sessionsCF)
	defer it.Close()

	it.SeekToFirst()
	for ; it.Valid(); it.Next() {

		var s *pb.SessionInfo = new(pb.SessionInfo)
		var buf []byte
		buf = make([]byte, it.Value().Size())
		copy(buf, it.Value().Data())
		err = proto.Unmarshal(buf, s)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to unmarshal proto from DB")
			return err, nil
		}

		log.WithFields(log.Fields{"key": it.Key().Data(),
			"session": s}).Debug("Iterating sessions from DB")

		sList = append(sList, s)
	}

	if err := it.Err(); err != nil {
		return err, sList
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return err, nil
	}
	return err, sList
}
