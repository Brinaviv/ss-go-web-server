package dal

import (
	"errors"
	"fmt"
	"github.com/brinaviv/ss-go-web-server/pkg/model"
)

var (
	ErrUserNotFound     = errors.New("entity not found")
	ErrAlreadyFollowing = errors.New("already following user")
	ErrNotFollowing     = errors.New("already not following user")
)

func NewErrUserNotFound(id model.UserID) error {
	return fmt.Errorf("user with id %s: %w", id.String(), ErrUserNotFound)
}

func NewErrAlreadyFollowing(user, toFollow model.UserID) error {
	return fmt.Errorf("user with id %v %w user with id %v", user, ErrAlreadyFollowing, toFollow)
}

func NewErrNotFollowing(user, toFollow model.UserID) error {
	return fmt.Errorf("user with id %v %w user with id %v", user, ErrNotFollowing, toFollow)
}
