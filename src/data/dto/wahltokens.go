package dto

type Wahltoken struct {
	ID         int    `gorm:"primaryKey;not null"`
	ElectionID int    `gorm:"column:election_id"`
	Token      string `gorm:"type:varchar(32)"`
	Voted      bool   `gorm:"column:voted"`
}

func (Wahltoken) TableName() string {
	return "wahltoken"
}

type WahltokenNotValidError struct {
	Message string
}

func (e WahltokenNotValidError) Error() string {
	return e.Message
}
