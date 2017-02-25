package db

/*
InstructorInfo{}
UserInfo{}
SessionInfo{}
UserInstructorReview{}
UserSessionReview{}
CreditCard{}
BankAcct{}
*/

type InstructorInfo struct {
	Firstname        string `sql:"not null"`
	Age              int32  `sql:"not null"`
	Sex              string `sql:"not null"`
	Location         string `sql:"not null"`
	Email            string `sql:"not null"`
	Password         string `sql:"not null"`
	Desc             string `gorm:"size:65536"`
	Certification    string `sql:"not null"`
	FitnessType      int32  `sql:"not null"`
	Lastname         string `sql:"not null"`
	City             string `sql:"not null"`
	Country          string `sql:"not null"`
	ID               int32  `gorm:"primary_key"`
	DisplayImage     []byte `gorm:"size:65536"`
	DisplayImageDesc string
	ReviewInfoID     string `sql:"not null"`
	BankInfoID       string `sql:"not null"`
}

type UserInfo struct {
	Firstname        string `sql:"not null"`
	Age              int32  `sql:"not null"`
	Sex              string `sql:"not null"`
	City             string `sql:"not null"`
	Country          string `sql:"not null"`
	Email            string `sql:"not null"`
	Password         string `sql:"not null"`
	Lastname         string `sql:"not null"`
	ID               int32  `gorm:"primary_key,auto_increment"`
	RecordActivityID int32  `sql:"not null"`
}

type SessionInfo struct {
	SessionTime          string `sql:"not null"`
	SessionDesc          string `gorm:"size:65536"`
	InstructorInfoID     int32  `sql:"not null"`
	SessionType          int32  `sql:"not null"`
	DifficultyLevel      int32  `sql:"not null"`
	TagList              int32  `sql:"not null"`
	ID                   int32  `gorm:"primary_key"`
	InstructorName       string `sql:"not null"`
	PreviewVideoUrl      string `sql:"not null"`
	DurationInMins       int32  `sql:"not null"`
	SessionUsersEnrolled uint32 `sql:"not null"`
}

type UserInstructorReview struct {
	InstructorRating float32 `sql:"not null"`
	InstructorReview string  `sql:"not null"`
	ID               int32   `gorm:"primary_key"`
	InstructorID     int32   `sql:"not null"`
	UserID           int32   `sql:"not null"`
}

type UserSessionReview struct {
	SessionRating float32 `sql:"not null"`
	SessionReview string  `sql:"not null"`
	ID            int32   `gorm:"primary_key"`
	SessionID     int32   `sql:"not null"`
	UserID        int32   `sql:"not null"`
}

type CreditCard struct {
	Name        string `sql:"not null"`
	Number      string `sql:"not null"`
	ExpiryMonth string `sql:"not null"`
	ExpiryYear  string `sql:"not null"`
	CCV         string `sql:"not null"`
	UserID      int32  `sql:"not null"`
	ID          int32  `gorm:"primary_key"`
}

type BankAcct struct {
	RoutingNum   string `sql:"not null"`
	AcctNum      string `sql:"not null"`
	BankName     string `sql:"not null"`
	InstructorID int32
	ID           int32 `gorm:"primary_key"`
}
