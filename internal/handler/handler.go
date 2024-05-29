package handler

import (
	"github.com/gin-gonic/gin"
)

type Ticket interface {
	Buy(c *gin.Context)
	Notify(c *gin.Context)
	Error(c *gin.Context)
	Success(c *gin.Context)
	GetTicket(c *gin.Context)
}

type Payment interface {
	GetPublishableKey(c *gin.Context)
	HandleCreatePaymentIntent(c *gin.Context)
}

type User interface {
	CreateUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type Event interface {
	PostEvent(c *gin.Context)
	FetchEvents(c *gin.Context)
}
