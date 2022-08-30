package entities

type User struct {
	ID       uint64  `gorm:"primary_key:auto_increment" json:"id"`
	Name     string  `gorm:"type:varchar(255);not null" json:"name"`
	Email    string  `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password string  `gorm:"->;<-;not null" json:"-"`
	Token    string  `gorm:"-" json:"token,omitempty"`
	Books    *[]Book `gorm:"foreignKey:UserId;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"books,omitempty"`
}
