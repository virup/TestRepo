package main

import (
	pb "server/rpcdefsql"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// Given a instructKey, return the UserInfo
func GetInstructorFromDB(iKey int32) (error, *pb.InstructorInfo) {
	var err error

	var i *pb.InstructorInfo = new(pb.InstructorInfo)

	err = db.First(i, iKey).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed to get instructor from DB")
		return err, i
	}
	log.WithFields(log.Fields{"instructorInfo": i, "key": iKey}).Debug("Read from DB")
	return err, i
}

func (s *server) GetInstructors(ctx context.Context, in *pb.GetInstructorsReq) (*pb.GetInstructorsReply, error) {

	var iList []pb.InstructorInfo
	var resp pb.GetInstructorsReply
	var err error
	//err = InsTable.Find(&iList).Error
	err = db.Find(&iList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get instructors from DB")
		return &resp, err
	}


	log.WithFields(log.Fields{"instructor": iList}).Debug("Get allinstructor success")
	for i, _ := range iList {
		resp.InstructorList = append(resp.InstructorList, &iList[i])
	}
	return &resp, nil
}

func (s *server) GetInstructor(ctx context.Context,
	in *pb.GetInstructorReq) (*pb.GetInstructorReply, error) {

	var resp pb.GetInstructorReply
	err, iInfo := GetInstructorFromDB(in.InstructorKey)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get session from DB")
		return &resp, err
	}

	log.WithFields(log.Fields{"instructor": iInfo}).Debug("Get instructor success")
	resp.Info = iInfo
	return &resp, nil
}

func (s *server) PostInstructorReview(ctx context.Context,
	in *pb.PostInstructorReviewReq) (*pb.PostInstructorReviewReply,
	error) {

	var resp pb.PostInstructorReviewReply
	err := db.Save(&in.Review).Error
	if err != nil {
		log.WithFields(log.Fields{"instructorReview": in,
			"error": err}).Error("Failed to write review to DB")
		return nil, err
	}
	resp.ReviewID = in.Review.ID
	log.WithFields(log.Fields{"instructorReview": in}).
		Debug("Added to DB")
	return &resp, nil
}

func (s *server) PostInstructorDisplayImg(ctx context.Context,
	in *pb.PostInstructorDisplayImgReq) (*pb.PostInstructorDisplayImgReply,
	error) {

	var err error
	var ins pb.InstructorInfo
	var resp pb.PostInstructorDisplayImgReply
	log.Debug("post Instructor image request")
	//err = db.Save(&in.Img).Error
	//if err != nil {
	//	log.WithFields(log.Fields{"instructorImage": in,
	//		"error": err}).Error("Failed to write image to DB")
	//	return &resp, err
	//}

	//err = db.First(&ins, in.InstructorInfoID).
	//	Update(pb.InstructorInfo{DisplayImageID: in.Img.ID}).Error

	err = db.First(&ins, in.InstructorInfoID).
		Update(pb.InstructorInfo{DisplayImage: in.Blob}).Error
	if err != nil {
		log.WithFields(log.Fields{"instructorImage": in,
			"error": err}).Error("Failed to update image ID to ins DB")
		return &resp, err
	}
	log.WithFields(log.Fields{"postInsImgResponse": resp}).
		Debug("instructor image add response")
	return &resp, nil
}

func (s *server) EnrollInstructor(ctx context.Context,
	in *pb.EnrollInstructorReq) (*pb.EnrollInstructorReply, error) {

	var err error
	var resp pb.EnrollInstructorReply
	log.Debug("Enroll Instructor request")
	err, resp.InstructorKey = postInstructorDB(*in.Instructor)
	if err != nil {
		log.WithFields(log.Fields{"instructor": in.Instructor,
			"error": err}).Error("Failed to write to DB")
		return &resp, err
	}
	log.WithFields(log.Fields{"enrollinstructor_response": resp}).
		Debug("Enrolled instructor response")
	return &resp, nil
}

func postInstructorDB(in pb.InstructorInfo) (err error, iKey int32) {

	log.WithFields(log.Fields{"instructorInfo": in}).Debug("Adding to DB")
	err = db.Save(&in).Error
	if err != nil {
		log.WithFields(log.Fields{"instructorInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, 0
	}
	iKey = in.ID
	log.WithFields(log.Fields{"instructorInfo": in}).
		Debug("Added to DB")
	return nil, iKey
}

func (s *server) RegisterInstructorBankAcct(ctx context.Context,
	in *pb.RegisterInstructorBankAcctReq) (*pb.RegisterInstructorBankAcctReply,
	error) {

	var resp pb.RegisterInstructorBankAcctReply
	err := db.Save(&in.BankAcct).Error
	if err != nil {
		log.WithFields(log.Fields{"instructorBankAcct": in,
			"error": err}).Error("Failed to write bank acct to DB")
		return nil, err
	}
	resp.BankAcctID = in.BankAcct.ID
	log.WithFields(log.Fields{"instructorBankAcct": in}).
		Debug("Added to DB")
	return &resp, nil
}

func (s *server) GetInstructorBankAcct(ctx context.Context,
	in *pb.GetInstructorBankAcctReq) (*pb.GetInstructorBankAcctReply,
	error) {

	var resp pb.GetInstructorBankAcctReply
	var err error

	var i *pb.BankAcct = new(pb.BankAcct)
	if in.BankAcctID > 0 {
		err = db.First(i, in.BankAcctID).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get instructor bank from DB")
			return &resp, err
		}
	} else if in.InstructorID > 0 {
		err = db.
			Where(pb.BankAcct{InstructorID: in.InstructorID}).
			Find(i).Error
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed" +
				" to get bankAcct from DB using instructorID")
			return &resp, err
		}
	}
	resp.BankAcct = i
	log.Debug("Success read back acct from DB")
	return &resp, nil
}
