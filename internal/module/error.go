package module

import (
	"fmt"
)

type UnknownServerError struct {
	Message string
	Cause   error
}

func (e *UnknownServerError) Error() string {
	return e.Message + "cause: " + e.Cause.Error()
}

type NotFoundError struct {
	ResourceName string
	ResourceId   string
	Context      string
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find resource %s with id %s while %s", t.ResourceName, t.ResourceId, t.Context)
}
