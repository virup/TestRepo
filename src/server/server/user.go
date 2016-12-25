package main

import (
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Given a userKey, return the UserInfo
func getUserFromDB(uKey string) (error, *pb.UserInfo) {
	var err error
	var buf []byte

	var u *pb.UserInfo = new(pb.UserInfo)

	v, err := rdb.GetCF(ro, usersCF, []byte(uKey))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get user from DB")
		return err, u
	}
	log.WithFields(log.Fields{"value": v}).Debug("Read" +
		"user value from DB")

	if v.Size() > 0 {
		buf = make([]byte, v.Size())
		copy(buf, v.Data())
		v.Free()
	} else {
		log.WithFields(log.Fields{"error": err}).Error("invalid key/" +
			"user from DB")
	}
	err = proto.Unmarshal(buf, u)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to unmarshal proto from DB")
		return err, nil
	}
	log.WithFields(log.Fields{"userInfo": u, "key": uKey}).
		Debug("Read from DB")
	return err, u
}

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserReq) (*pb.EnrollUserReply, error) {

	var err error
	var resp pb.EnrollUserReply
	log.Debug("Enroll User request")
	//err := db.Save(&u)

	err, resp.UserKey = postUserDB(*in.User)
	if err != nil {
		log.WithFields(log.Fields{"user": in.User, "error": err}).
			Error("Failed to write to DB for user")
		return &resp, err
	}
	log.WithFields(log.Fields{"user": in.User}).Debug("Added to DB")

	/*
		err, customerPayID := pay.CreatePayingCustomer(
			"mycustomer@gmail.com", "1234-xxxx-xxxx", "06", "19")
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to set up customer payment")
			return &resp, err
		}
		log.WithFields(log.Fields{"customerPayID": customerPayID}).
			Debug("Got new customer payment ID")

		err = pay.StartSubscription(customerPayID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to start subscription")
			return &resp, err
		}
	*/
	return &resp, nil
}

func (s *server) GetUser(ctx context.Context,
	in *pb.GetUserReq) (*pb.GetUserReply, error) {

	var resp pb.GetUserReply
	err, userInfo := getUserFromDB(in.UserKey)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get user from DB")
		return &resp, err
	}

	log.WithFields(log.Fields{"userInfo": userInfo}).Debug("Get user success")
	resp.Info = userInfo
	return &resp, nil
}

func (s *server) GetUsers(ctx context.Context,
	in *pb.GetUsersReq) (*pb.GetUsersReply, error) {

	var resp pb.GetUsersReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get users from DB")
		return &resp, err
	}
	return &resp, nil
}

func (s *server) Login(ctx context.Context,
	in *pb.LoginReq) (*pb.LoginReply, error) {

	var resp pb.LoginReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to login")
		return &resp, err
	}
	return &resp, nil
}

func postUserDB(in pb.UserInfo) (err error, uKey string) {

	log.WithFields(log.Fields{"userInfo": in}).Debug("Adding to DB")
	uKey = GetRandomID()

	byteBuf, err := proto.Marshal(&in)
	if err != nil {
		log.WithFields(log.Fields{"userInfo": in, "error": err}).
			Error("Failed to convert to binary")
		return err, ""
	}

	err = rdb.PutCF(wo, usersCF, []byte(uKey), byteBuf)
	if err != nil {
		log.WithFields(log.Fields{"userInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, ""
	}
	log.WithFields(log.Fields{"userInfo": in, "key": uKey}).
		Debug("Added to DB")
	return nil, uKey
}
