package inmemory

import (
	"github.com/brinaviv/ss-go-web-server/pkg/dal"
)

func NewDAO(generateUUID UserIDGenerator) dal.DAO {
	return dal.DAO{
		Users: newUserDAO(generateUUID),
	}
}
