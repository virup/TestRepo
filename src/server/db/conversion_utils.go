package db

import (
	"reflect"
	pb "server/rpcdefsql"
	"fmt"
)

func typeof(in interface{}) reflect.Type {
	return reflect.TypeOf(in)
}

func Convert(in interface{}) interface{} {
	typeName := reflect.TypeOf(in)

	switch typeName {
	case typeof(&pb.InstructorInfo{}):
		return ConvertInstructorInfoToDb(in.(*pb.InstructorInfo))
	case typeof(&pb.UserInfo{}):
		return ConvertUserInfoToDb(in.(*pb.UserInfo))
	case typeof(SessionInfo{}):
		return ConvertSessionInfoToDb(in.(*pb.SessionInfo))
	case typeof(&pb.UserInstructorReview{}):
		return ConvertUserSessionReviewToDb(in.(*pb.UserSessionReview))
	case typeof(&pb.UserInstructorReview{}):
		return ConvertUserInstructorReviewToDb(in.(*pb.UserInstructorReview))
	case typeof(CreditCard{}):
		return ConvertCreditCardToDb(in.(*pb.CreditCard))
	case typeof(&pb.BankAcct{}):
		return ConvertBankAcctToDb(in.(*pb.BankAcct))
	case typeof(&InstructorInfo{}):
		return ConvertInstructorInfoToRPC(in.(*InstructorInfo))
	case typeof(&UserInfo{}):
		return ConvertUserInfoToRPC(in.(*UserInfo))
	case typeof(&SessionInfo{}):
		return ConvertSessionInfoToRPC(in.(*SessionInfo))
	case typeof(&UserInstructorReview{}):
		return ConvertUserInstructorReviewToRPC(in.(*UserInstructorReview))
	case typeof(&UserSessionReview{}):
		return ConvertUserSessionReviewToRPC(in.(*UserSessionReview))
	case typeof(&CreditCard{}):
		return ConvertCreditCardToRPC(in.(*CreditCard))
	case typeof(&BankAcct{}):
		return ConvertBankAcctToRPC(in.(*BankAcct))
	}
	fmt.Println("NIL")
	return nil
}

func ConvertInstructorInfoToDb(in *pb.InstructorInfo) *InstructorInfo {
	return &InstructorInfo{
		Firstname:        in.GetFirstname(),
		Age:              in.GetAge(),
		Sex:              in.GetSex(),
		Location:         in.GetLocation(),
		Email:            in.GetEmail(),
		Password:         in.GetPassword(),
		Desc:             in.GetDesc(),
		Certification:    in.GetCertification(),
		FitnessType:      int32(in.GetFitnessType()),
		Lastname:         in.GetLastname(),
		City:             in.GetCity(),
		Country:          in.GetCountry(),
		DisplayImage:     in.GetDisplayImage(),
		DisplayImageDesc: in.GetDisplayImageDesc(),
		ReviewInfoID:     in.GetReviewInfoID(),
		BankInfoID:       in.GetBankInfoID(),
	}
}

func ConvertInstructorInfoToRPC(in *InstructorInfo) *pb.InstructorInfo {
	return &pb.InstructorInfo{
		Firstname: in.Firstname,
		Age: in.Age,
		Sex: in.Sex,
		Location: in.Location,
		Email: in.Email,
		Password: in.Password,
		Desc: in.Desc,
		Certification: in.Certification,
		FitnessType: pb.FitnessCategory(in.FitnessType),
		Lastname: in.Lastname,
		City: in.City,
		Country: in.Country,
		DisplayImage: in.DisplayImage,
		DisplayImageDesc: in.DisplayImageDesc,
		ReviewInfoID: in.ReviewInfoID,
		BankInfoID: in.BankInfoID,
	}
}

func ConvertUserInfoToDb(in *pb.UserInfo) *UserInfo {
	return &UserInfo{
		Firstname:        in.GetFirstname(),
		Age:              in.GetAge(),
		Sex:              in.GetSex(),
		City:             in.GetCity(),
		Country:          in.GetCountry(),
		Email:            in.GetEmail(),
		Password:         in.GetPassword(),
		Lastname:         in.GetLastname(),
		RecordActivityID: in.GetRecordActivityID(),
	}
}

func ConvertUserInfoToRPC(in *UserInfo) *pb.UserInfo {
	return &pb.UserInfo{
		Firstname:        in.Firstname,
		Age:              in.Age,
		Sex:              in.Sex,
		City:             in.City,
		Country:          in.Country,
		Email:            in.Email,
		Password:         in.Password,
		Lastname:         in.Lastname,
		RecordActivityID: in.RecordActivityID,
	}
}

func ConvertSessionInfoToDb(in *pb.SessionInfo) *SessionInfo {
	return &SessionInfo{
		SessionTime:          in.GetSessionTime(),
		SessionDesc:          in.GetSessionDesc(),
		InstructorInfoID:     in.GetInstructorInfoID(),
		SessionType:          int32(in.GetSessionType()),
		DifficultyLevel:      int32(in.GetDifficultyLevel()),
		TagList:              int32(in.GetTagList()),
		InstructorName:       in.GetInstructorName(),
		PreviewVideoUrl:      in.GetPreviewVideoUrl(),
		DurationInMins:       in.GetDurationInMins(),
		SessionUsersEnrolled: in.GetSessionUsersEnrolled(),
	}
}

func ConvertSessionInfoToRPC(in *SessionInfo) *pb.SessionInfo {
	return &pb.SessionInfo{
		SessionTime:          in.SessionTime,
		SessionDesc:          in.SessionDesc,
		InstructorInfoID:     in.InstructorInfoID,
		SessionType:          pb.FitnessCategory(in.SessionType),
		DifficultyLevel:      pb.SessionDifficulty(in.DifficultyLevel),
		TagList:              pb.SessionTag(in.TagList),
		InstructorName:       in.InstructorName,
		PreviewVideoUrl:      in.PreviewVideoUrl,
		DurationInMins:       in.DurationInMins,
		SessionUsersEnrolled: in.SessionUsersEnrolled,
	}
}

func ConvertUserInstructorReviewToDb(in *pb.UserInstructorReview) *UserInstructorReview {
	return &UserInstructorReview{
		InstructorRating: in.GetInstructorRating(),
		InstructorReview: in.GetInstructorReview(),
		InstructorID:     in.GetInstructorID(),
		UserID:           in.GetUserID(),
	}
}

func ConvertUserInstructorReviewToRPC(in *UserInstructorReview) *pb.UserInstructorReview {
	return &pb.UserInstructorReview{
		InstructorRating: in.InstructorRating,
		InstructorReview: in.InstructorReview,
		InstructorID:     in.InstructorID,
		UserID:           in.UserID,
	}
}

func ConvertUserSessionReviewToDb(in *pb.UserSessionReview) *UserSessionReview {
	return &UserSessionReview{
		SessionRating: in.GetSessionRating(),
		SessionReview: in.GetSessionReview(),
		SessionID:     in.GetSessionID(),
		UserID:        in.GetUserID(),
	}
}

func ConvertUserSessionReviewToRPC(in *UserSessionReview) *pb.UserSessionReview {
	return &pb.UserSessionReview{
		SessionRating: in.SessionRating,
		SessionReview: in.SessionReview,
		SessionID:     in.SessionID,
		UserID:        in.UserID,
	}
}

func ConvertCreditCardToDb(in *pb.CreditCard) *CreditCard {
	return &CreditCard{
		Name:        in.GetName(),
		Number:      in.GetNumber(),
		ExpiryMonth: in.GetExpiryMonth(),
		ExpiryYear:  in.GetExpiryYear(),
		CCV:         in.GetCCV(),
		UserID:      in.GetUserID(),
	}
}

func ConvertCreditCardToRPC(in *CreditCard) *pb.CreditCard {
	return &pb.CreditCard{
		Name:        in.Name,
		Number:      in.Number,
		ExpiryMonth: in.ExpiryMonth,
		ExpiryYear:  in.ExpiryYear,
		CCV:         in.CCV,
		UserID:      in.UserID,
	}
}

func ConvertBankAcctToDb(in *pb.BankAcct) *BankAcct {
	return &BankAcct{
		RoutingNum:   in.GetRoutingNum(),
		AcctNum:      in.GetAcctNum(),
		BankName:     in.GetBankName(),
		InstructorID: int32(in.GetInstructorID()),
	}
}

func ConvertBankAcctToRPC(in *BankAcct) *pb.BankAcct {
	return &pb.BankAcct{
		RoutingNum:   in.RoutingNum,
		AcctNum:      in.AcctNum,
		BankName:     in.BankName,
		InstructorID: in.InstructorID,
	}
}