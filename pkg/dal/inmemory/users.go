package inmemory

import (
	"github.com/brinaviv/ss-go-web-server/pkg/dal"
	"github.com/brinaviv/ss-go-web-server/pkg/model"
	"sync"
	"time"
)

type UserIDGenerator func() model.UserID

type userDAO struct {
	entities       map[model.UserID]*model.UserEntity
	mu             sync.Mutex
	generateUserID UserIDGenerator
}

func newUserDAO(generateUserID UserIDGenerator) dal.UserDAO {
	return &userDAO{entities: map[model.UserID]*model.UserEntity{}, generateUserID: generateUserID}
}

func (dao *userDAO) FetchByID(id model.UserID) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	e, ok := dao.entities[id]
	if !ok {
		return e, &dal.ErrIDNotFound{EntityName: "user", Id: id.String()}
	}
	return e, nil
}

func (dao *userDAO) Create(user *model.User) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	e := dao.createUserEntity(user)
	dao.entities[e.ID] = e
	return e, nil
}

func (dao *userDAO) Update(id model.UserID, user *model.User) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	e, ok := dao.entities[id]
	if !ok {
		dao.entities[id] = dao.createUserEntity(user)
		return dao.entities[id], nil
	}
	e.UpdatedAt = time.Now()
	e.Data = user
	return e, nil
}

func (dao *userDAO) Delete(id model.UserID) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	_, ok := dao.entities[id]
	if !ok {
		return &dal.ErrIDNotFound{EntityName: "user", Id: id.String()}
	}
	delete(dao.entities, id)
	return nil
}

func (dao *userDAO) createUserEntity(user *model.User) *model.UserEntity {
	now := time.Now()
	return &model.UserEntity{Data: user, ID: model.UserID(dao.generateUserID()), CreatedAt: now, UpdatedAt: now}
}
