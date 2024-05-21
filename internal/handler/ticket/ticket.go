package ticket

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ticket struct {
	log      slog.Logger
	errorUrl string
	module   module.Ticket
	sms      module.Sms
	email    module.Email
}

func Init(log slog.Logger, url string, module module.Ticket, sms module.Sms, email module.Email) handler.Ticket {
	return &ticket{
		log:      log,
		errorUrl: url,
		module:   module,
		sms:      sms,
		email:    email,
	}
}
func (t *ticket) Buy(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		t.log.Error("failed to bind user input", err)
		c.HTML(http.StatusBadRequest, "error.html", err)
		return
	}
	//validation
	if err := user.Validate(); err != nil {
		fmt.Println("email:", user.Email)
		t.log.Error("invalid user input", err)
		c.HTML(http.StatusBadRequest, "error.html", err)
		return
	}
	err := t.module.CreateCheckoutSession(c, user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", err)
		return
	}

}

func (t *ticket) Notify(c *gin.Context) {

	var notify model.Notification
	if err := c.ShouldBindJSON(&notify); err != nil {
		t.log.Error("failed to bind notification response body", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := t.module.UpdatePaymentStatus(notify.TransactionStatus, notify.SessionID)
	if err != nil {
		t.log.Error("failed to update payment status", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	pdf, err := t.module.GeneratePDFTicket(user)
	if err != nil {
		t.log.Error("Error generating PDF ticket:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	attachmentPath := fmt.Sprintf("./public/pdfs/attachment_%s.pdf", user.SessionID)
	_, err = os.Create(attachmentPath)
	if err != nil {
		t.log.Error("error when creating attachment file for email")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err := pdf.WritePdf(attachmentPath); err != nil {
		t.log.Error("error when copying pdf to email attachment file", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go t.sms.SendSms(user, wg)
	go t.email.SendEmail(user, attachmentPath, wg)
	wg.Wait()

	c.Writer.WriteHeader(http.StatusOK)
}

func (t *ticket) Error(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "error.html", "Your request failed!Please try again")

}

func (t *ticket) Success(c *gin.Context) {

	c.HTML(http.StatusOK, "success.html", nil)
}

func (t *ticket) GetTicket(c *gin.Context) {
	nonce := c.Param("nonce")
	user, err := t.module.GetUser(nonce)
	if err != nil {
		t.log.Error("failed to get user by nonce", err)
		c.HTML(http.StatusInternalServerError, "error.err", err)
		return
	}
	if user.PaymentStatus == "pending" {
		user, err := t.module.UpdatePaymentStatus("SUCCESS", user.SessionID)
		if err != nil {
			t.log.Error("failed to updated payment status", err)
			return
		}
		pdf, err := t.module.GeneratePDFTicket(user)
		if err != nil {
			t.log.Error("Error generating PDF ticket:", err)
			return
		}
		attachmentPath := fmt.Sprintf("./public/pdfs/attachment_%s.pdf", user.SessionID)
		_, err = os.Create(attachmentPath)
		if err != nil {
			t.log.Error("error when creating attachment file for email")
			return
		}

		if err := pdf.WritePdf(attachmentPath); err != nil {
			t.log.Error("error when copying pdf to email attachment file", err)
			return
		}
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go t.sms.SendSms(user, wg)
		go t.email.SendEmail(user, attachmentPath, wg)
	}
	pdf, err := t.module.GeneratePDFTicket(user)
	if err != nil {
		t.log.Error("failed to generate pdf by the given nonce", err)
		c.HTML(http.StatusInternalServerError, "error.html", err)
		return
	}

	ticketPath := fmt.Sprintf("./public/pdfs/ticket_%s.pdf", user.SessionID)
	_, err = os.Create(ticketPath)
	if err != nil {
		t.log.Error("failed to create pdf file", err)
		c.HTML(http.StatusInternalServerError, "error.html", err)
		return
	}

	err = pdf.WritePdf(ticketPath)
	if err != nil {
		t.log.Error("failed to write pdf file", err)
		c.HTML(http.StatusInternalServerError, "error.html", err)
		return
	}

	c.File(ticketPath)

	if err := os.RemoveAll(ticketPath); err != nil {
		t.log.Error("failed to remove pdf file", err)
		c.HTML(http.StatusInternalServerError, "error.html", err)
		return
	}

}
