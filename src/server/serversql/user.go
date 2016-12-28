package main

import (
	"pay"
	pb "server/rpcdefsql"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var lastUserID int32 = 1

// XXX Retrieve from persisted
func getUserID() int32 {

	ret := lastUserID
	lastUserID++
	return ret
}

// Given a userKey, return the UserInfo
func getUserFromDB(uKey int32) (error, *pb.UserInfo) {
	var err error
	var u *pb.UserInfo = new(pb.UserInfo)

	err = UserTable.First(u, uKey).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get user from DB")
		return err, u
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

	// Enable payment
	if false {
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
	}
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

	var uList []pb.UserInfo
	var resp pb.GetUsersReply
	var err error

	err = UserTable.Find(&uList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get users from DB")
		return &resp, err
	}

	log.Printf("\n")
	log.WithFields(log.Fields{"users": uList}).Debug("Get alluser success")
	for i, _ := range uList {
		resp.UserList = append(resp.UserList, &uList[i])
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

func postUserDB(in pb.UserInfo) (err error, uKey int32) {

	log.WithFields(log.Fields{"userInfo": in}).Debug("Adding to DB")
	uKey = getUserID()
	in.ID = uKey
	err = UserTable.Save(&in).Error
	if err != nil {
		log.WithFields(log.Fields{"userinfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, 0
	}
	log.WithFields(log.Fields{"userInfo": in, "key": uKey}).
		Debug("Added to DB")
	return nil, uKey
}
