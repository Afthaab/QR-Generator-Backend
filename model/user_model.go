package model

// Student struct represents a student entity
type Student struct {
	Id         string            `bson:"_id,omitempty" json:"id"`
	Name       string            `bson:"name" json:"name"`
	Branch     string            `bson:"branch" json:"branch"`
	Batch      string            `bson:"batch" json:"batch"`
	Phoneno    string            `bson:"phoneo" json:"phoneno"`
	Email      string            `bson:"email" json:"email"`
	Image      string            `json:"image"`
	Password   string            `bson:"Password" json:"Password"`
	Attandance map[string]string `bson:"attandance" json:"attandance"`
}

type RegisterAttendance struct {
	Id string `json:"id"`
}

// Attandance struct represents a date for attendance
type GetAttendance struct {
	Date string `json:"date"`
}

// TeacherData struct represents teacher registration data
type TeacherData struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Repeatpassword string `json:"repeatpassword"`
	Securitykey    string `json:"securitykey"`
}

// TeacherSignIn struct represents teacher sign-in data
type TeacherSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TimeSheet struct {
	Date      string   `bson:"date" json:"date"`
	Attendees []string `bson:"Attendees" json:"Attendees"`
}

type StudentSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
