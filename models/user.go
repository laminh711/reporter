package models

// User model
type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	IsNew    bool   `json:"first_time" bson:"is_new"`
}
