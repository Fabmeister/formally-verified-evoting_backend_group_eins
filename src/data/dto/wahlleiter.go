package dto

type Wahlleiter struct {
	ID       int
	Username string
	Email    string
	Password string
	Salt     string
}

func (Wahlleiter) TableName() string {
	return "wahlleiter"
}
