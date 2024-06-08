package handler

import (
	"github.com/gin-gonic/gin"
)

type Ticket interface {
	GetTicket(c *gin.Context)
}

type Payment interface {
	GetPublishableKey(c *gin.Context)
	HandleCreatePaymentIntent(c *gin.Context)
	PaymentWebhook(c *gin.Context)
}

type User interface {
	CreateUser(c *gin.Context)
	LoginUser(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type Event interface {
	PostEvent(c *gin.Context)
	FetchEvents(c *gin.Context)
	FetchEvent(c *gin.Context)
}
