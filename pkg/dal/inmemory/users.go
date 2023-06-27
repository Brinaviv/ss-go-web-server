package inmemory

import (
	"github.com/brinaviv/ss-go-web-server/pkg/dal"
	"github.com/brinaviv/ss-go-web-server/pkg/model"
	"golang.org/x/exp/maps"
	"sync"
	"time"
)

type UserIDGenerator func() model.UserID

type userDAO struct {
	users          map[model.UserID]*model.UserEntity
	following      map[model.UserID]map[model.UserID]bool // following[id] = set of users the user follows
	followers      map[model.UserID]map[model.UserID]bool // followers[id] = set of users who follow user
	mu             sync.Mutex
	generateUserID UserIDGenerator
}

func newUserDAO(generateUserID UserIDGenerator) dal.UserDAO {
	return &userDAO{
		users:          make(map[model.UserID]*model.UserEntity),
		following:      make(map[model.UserID]map[model.UserID]bool),
		followers:      make(map[model.UserID]map[model.UserID]bool),
		generateUserID: generateUserID}
}

func (dao *userDAO) FetchByID(id model.UserID) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	e, ok := dao.users[id]
	if !ok {
		return e, dal.NewErrUserNotFound(id)
	}
	return e, nil
}

func (dao *userDAO) Create(user *model.User) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	return dao.createNewEntityWithLock(user), nil
}

func (dao *userDAO) Update(id model.UserID, user *model.User) (*model.UserEntity, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	e, ok := dao.users[id]
	if !ok {
		dao.users[id] = dao.createNewEntityWithLock(user, id)
		return dao.users[id], nil
	}
	e.Data = user
	e.UpdatedAt = time.Now()
	return e, nil
}

// Delete
// O(f) time, f is the number of followers
func (dao *userDAO) Delete(id model.UserID) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	_, ok := dao.users[id]
	if !ok {
		return dal.NewErrUserNotFound(id)
	}
	for fid := range dao.following[id] {
		delete(dao.followers[fid], id)
	}
	delete(dao.following, id)
	delete(dao.users, id)
	return nil
}

func (dao *userDAO) Follow(userID model.UserID, targetID model.UserID) (*model.Following, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	userFollowing, ok := dao.following[userID] // user does not exist
	if !ok {
		return nil, dal.NewErrUserNotFound(userID)
	}

	if userFollowing[targetID] {
		return nil, dal.NewErrAlreadyFollowing(userID, targetID)
	}

	targetFollowers, ok := dao.followers[targetID]
	if !ok { // target does not exist
		return nil, dal.NewErrUserNotFound(targetID)
	}

	userFollowing[targetID], targetFollowers[userID] = true, true

	return &model.Following{Following: maps.Keys(userFollowing)}, nil
}

func (dao *userDAO) Unfollow(userID model.UserID, targetID model.UserID) (*model.Following, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	userFollowing, ok := dao.following[userID] // user does not exist
	if !ok {
		return nil, dal.NewErrUserNotFound(userID)
	}

	targetFollowers, ok := dao.followers[targetID]
	if !ok { // target does not exist
		return nil, dal.NewErrUserNotFound(targetID)
	}

	if !userFollowing[targetID] {
		return nil, dal.NewErrNotFollowing(userID, targetID)
	}

	delete(userFollowing, targetID)
	delete(targetFollowers, userID)

	return &model.Following{Following: maps.Keys(userFollowing)}, nil
}

func (dao *userDAO) createNewEntityWithLock(user *model.User, optionalID ...model.UserID) *model.UserEntity {
	e := dao.createUserEntity(user, optionalID...)
	dao.users[e.ID] = e
	dao.following[e.ID] = make(map[model.UserID]bool)
	dao.followers[e.ID] = make(map[model.UserID]bool)
	return e
}

func (dao *userDAO) createUserEntity(user *model.User, optionalID ...model.UserID) *model.UserEntity {
	var id model.UserID
	if len(optionalID) > 0 {
		id = optionalID[0]
	} else {
		id = dao.generateUserID()
	}
	now := time.Now()
	return &model.UserEntity{Data: user, ID: id, CreatedAt: now, UpdatedAt: now}
}
