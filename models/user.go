package models

type User struct {
	ID       string   `json:"id,omitempty" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	ApiKeys  []string `json:"api_keys,omitempty" bson:"api_keys"`
}
