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
