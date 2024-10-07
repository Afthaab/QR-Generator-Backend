package model

type User struct {
	Id      string `bson:"_id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Branch  string `bson:"branch" JSON:"branch"`
	Batch   string `bson:"batch" JSON:"batch"`
	Phoneno string `bson:"phoneo" JSON:"phoneno"`
}

type TeacherData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Newpassword string `json:"newpassword"`
}

type TeacherSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
