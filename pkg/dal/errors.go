package dal

import "fmt"

type ErrIDNotFound struct {
	EntityName string
	Id         string
}

func (e *ErrIDNotFound) Error() string {
	return fmt.Sprintf("could not find %s with id %s", e.EntityName, e.Id)
}
