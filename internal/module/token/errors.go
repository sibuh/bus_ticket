package token

import (
	"event_ticket/internal/constant"
	"fmt"
)

type UnknownServerError struct {
	message string
	cause   error
}

func (e *UnknownServerError) Error() string {
	return e.message + "cause: " + e.cause.Error()
}

type NotFoundError struct {
	resourceName string
	resourceId   string
	context      string
}
type ErrStatusNotUpdated struct {
	ID     string
	Status string
	Retry  bool
}
type ErrInvalidTicketStatus struct {
	ID      string
	Status  string
	Message string
}

func (e *ErrInvalidTicketStatus) Error() string {
	return e.Message
}

func (e *ErrStatusNotUpdated) Error() string {
	return fmt.Sprintf("status of ticket with id %s is not yet updated.Expected status %s but got %s", e.ID, constant.Reserved, e.Status)
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find resource %s with id %s while %s", t.resourceName, t.resourceId, t.context)
}

type ErrInvalidPayload struct {
	payload interface{}
	message string
}

func (e ErrInvalidPayload) Error() string {
	return fmt.Sprintf("invalid ticket token payload got %v", e.payload)
}
