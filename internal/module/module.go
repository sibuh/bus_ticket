package module

import (
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
)

type Ticket interface {
	CreateCheckoutSession(c *gin.Context, user model.User) error
	UpdatePaymentStatus(status, sid string) (db.User, error)
	GeneratePDFTicket(userData db.User) (*gopdf.GoPdf, error)
	GetUser(nonce string) (db.User, error)
}

type Sms interface {
	SendSms(user db.User, wg *sync.WaitGroup) error
}

type Email interface {
	SendEmail(user db.User, attachmentPath string, wg *sync.WaitGroup) error
}
