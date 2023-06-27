package dal

import "github.com/brinaviv/ss-go-web-server/pkg/model"

type UserDAO interface {
	FetchByID(id model.UserID) (*model.UserEntity, error)
	Create(*model.User) (*model.UserEntity, error)
	Update(id model.UserID, user *model.User) (*model.UserEntity, error)
	Delete(id model.UserID) error
}
