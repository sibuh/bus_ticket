package module

import (
	"context"
	"event_ticket/internal/model"
)

/*
	type Ticket interface {
		CreateCheckoutSession(c *gin.Context, user model.User) error
		UpdatePaymentStatus(status, sid string) (db.User, error)
		GeneratePDFTicket(userData db.User) (*gopdf.GoPdf, error)
	}

	type Sms interface {
		SendSms(user db.User, wg *sync.WaitGroup) error
	}

	type Email interface {
		SendEmail(user db.User, attachmentPath string, wg *sync.WaitGroup) error
	}
*/
type User interface {
	CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error)
	GetUser(ctx context.Context, id int32) (model.User, error)
	LoginUser(ctx context.Context) (string, error)
}
