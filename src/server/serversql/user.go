package main

import (
	"errors"
	"pay"
	pb "server/rpcdefsql"
	"strings"

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

func (s *server) SubscribeUser(ctx context.Context,
	in *pb.SubscribeUserReq) (*pb.SubscribeUserReply, error) {

	var resp pb.SubscribeUserReply
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

func (s *server) EnrollUser(ctx context.Context,
	in *pb.EnrollUserReq) (*pb.EnrollUserReply, error) {

	var err error
	var resp pb.EnrollUserReply
	log.Debug("Enroll User request")

	if in.User.Email == "" || in.User.PassWord == "" {
		log.WithFields(log.Fields{"user": in.User, "error": err}).
			Error("Invalid email/password for user")
		return &resp, errors.New("Invalid email/password")

	}
	err, resp.UserKey = postUserDB(*in.User)
	if err != nil {
		log.WithFields(log.Fields{"user": in.User, "error": err}).
			Error("Failed to write to DB for user")
		return &resp, err
	}
	log.WithFields(log.Fields{"user": in.User}).Debug("Added to DB")
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
	var uList []pb.UserInfo
	var iList []pb.InstructorInfo
	userFound := true
	insFound := true

	err = UserTable.
		Where(pb.UserInfo{Email: in.Email}).
		Find(&uList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get user from DB with email")
		userFound = false
	}

	err = InsTable.
		Where(pb.InstructorInfo{Email: in.Email}).
		Find(&iList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get ins from DB with email")
		insFound = false
	}

	if !insFound && !userFound {
		log.WithFields(log.Fields{"loginreq": in}).Error("Failed" +
			" to get user or ins from DB with email")
		err = errors.New("Invalid email/password for any role")
		return &resp, err
	}

	if len(iList) > 0 {
		if 0 == strings.Compare(in.PassWord, iList[0].PassWord) {
			log.WithFields(log.Fields{"insLoginInfo": in}).
				Debug("Authenticated instructor")
			resp.Instructor = &iList[0]
			resp.PersonType = pb.PersonRole_ROLE_INSTRUCTOR
		} else {
			log.WithFields(log.Fields{"loginReq": in}).
				Error("Invalid password for instructor")
			return &resp, errors.New("Invalid password for instructor")
		}
	} else if len(uList) > 0 {
		if 0 == strings.Compare(in.PassWord, uList[0].PassWord) {
			log.WithFields(log.Fields{"insLoginInfo": in}).
				Debug("Authenticated user")
			resp.User = &uList[0]
			resp.PersonType = pb.PersonRole_ROLE_USER
		} else {
			log.WithFields(log.Fields{"loginReq": in}).
				Error("Invalid password for user")
			return &resp, errors.New("Invalid password for user")
		}
	} else {
		err = errors.New("Couldn't find user/instructor in DB")
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
