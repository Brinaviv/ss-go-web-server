package model

import (
	"fmt"
	"time"
)

type Entity[T any, ID fmt.Stringer] struct {
	ID        ID        `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      T         `json:"data"`
}
