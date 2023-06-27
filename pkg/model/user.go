package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserID uuid.UUID

func (uid UserID) String() string {
	return uuid.UUID(uid).String()
}

func ParseUserID(id string) (UserID, error) {
	parse, err := uuid.Parse(id)
	return UserID(parse), err
}

type UserEntity Entity[*User, UserID]

func (u *UserEntity) UpdateTime() {
	u.UpdatedAt = time.Now()
}

func (u *UserEntity) MarshalJSON() ([]byte, error) {
	type Alias UserEntity
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    u.ID.String(),
		Alias: (*Alias)(u),
	})
}

func ValidateUser(_ *User) error {
	return nil
}
