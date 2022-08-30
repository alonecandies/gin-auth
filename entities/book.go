package entities

type Book struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text;not null" json:"description"`
	UserId      uint64 `gorm:"not null" json:"-"`
	User        User   `gorm:"foreignKey:UserId;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user,omitempty"`
}
