package models

// User owns documents, concepts, and their relationships.
type User struct {
	Base

	Username string  `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Name     string  `json:"name" gorm:"type:varchar(100);not null"`
	Email    *string `json:"email,omitempty" gorm:"type:varchar(320);uniqueIndex"`
	GoogleID *string `json:"-" gorm:"type:varchar(255);uniqueIndex"`
	GithubID *string `json:"-" gorm:"type:varchar(255);uniqueIndex"`
}
