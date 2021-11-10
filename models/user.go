package models

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is used to login and keep authentication information
type User struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`

	Groups []uuid.UUID `json:"groups"`

	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// CheckPasswordHash is used to validate that a given password matches the stored hash
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsMemberOfGroup checks if a user is in a given set of groups
func (u *User) IsMemberOfGroup(groups ...uuid.UUID) bool {
	return util.ContainsOneOfUUIDs(&u.Groups, &groups)
}

// SetPassword updates the user object's password hash
func (u *User) SetPassword(password string) error {
	password, err := hashPassword(password)
	if err != nil {
		return err
	}

	u.Password = password
	return nil
}

// Private Functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
