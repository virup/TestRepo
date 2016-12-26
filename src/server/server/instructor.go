package main

import (
	pb "server/rpcdef"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"

	"golang.org/x/net/context"
)

// Given a instructKey, return the UserInfo
func GetInstructorFromDB(iKey string) (error, *pb.InstructorInfo) {
	var err error
	var buf []byte

	var i *pb.InstructorInfo = new(pb.InstructorInfo)

	v, err := rdb.GetCF(ro, instructorsCF, []byte(iKey))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get instructor from DB")
		return err, i
	}
	log.WithFields(log.Fields{"value": v}).Debug("Read" +
		"instructor value from DB")

	if v.Size() > 0 {
		buf = make([]byte, v.Size())
		copy(buf, v.Data())
		v.Free()
	} else {
		log.WithFields(log.Fields{"error": err}).Error("invalid key/" +
			"instructor from DB")
	}
	err = proto.Unmarshal(buf, i)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to unmarshal proto from DB")
		return err, nil
	}
	log.WithFields(log.Fields{"instructorInfo": i, "key": iKey}).
		Debug("Read from DB")
	return err, i
}

func (s *server) GetInstructors(ctx context.Context,
	in *pb.GetInstructorsReq) (*pb.GetInstructorsReply, error) {

	var resp pb.GetInstructorsReply
	var err error
	//err = rdb.GetCF(wo, sessionsCF, []byte(sessionKey), binBuf.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get instructors from DB")
		return &resp, err
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
	log.WithFields(log.Fields{"instructor": in.Instructor}).
		Debug("Added to DB")
	return &resp, nil
}

func postInstructorDB(in pb.InstructorInfo) (err error, iKey string) {

	log.WithFields(log.Fields{"instructorInfo": in}).Debug("Adding to DB")
	iKey = GetRandomID()

	byteBuf, err := proto.Marshal(&in)
	if err != nil {
		log.WithFields(log.Fields{"instructorInfo": in, "error": err}).
			Error("Failed to convert to binary")
		return err, ""
	}

	err = rdb.PutCF(wo, instructorsCF, []byte(iKey), byteBuf)
	if err != nil {
		log.WithFields(log.Fields{"instructorInfo": in, "error": err}).
			Error("Failed to write to DB")
		return err, ""
	}
	log.WithFields(log.Fields{"instructorInfo": in, "key": iKey}).
		Debug("Added to DB")
	return nil, iKey
}
