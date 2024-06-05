package ticket

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ticket struct {
	log     slog.Logger
	payment module.Payment
	tkt     module.Ticket
}

func Init(log slog.Logger, pmt module.Payment, tkt module.Ticket) handler.Ticket {
	return &ticket{
		log:     log,
		payment: pmt,
		tkt:     tkt,
	}
}

// func (t *ticket) Buy(c *gin.Context) {
// 	var user model.User
// 	if err := c.ShouldBind(&user); err != nil {
// 		t.log.Error("failed to bind user input", err)
// 		c.HTML(http.StatusBadRequest, "error.html", err)
// 		return
// 	}
// 	//validation
// 	if err := user.Validate(); err != nil {
// 		fmt.Println("email:", user.Email)
// 		t.log.Error("invalid user input", err)
// 		c.HTML(http.StatusBadRequest, "error.html", err)
// 		return
// 	}
// 	err := t.module.CreateCheckoutSession(c, user)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", err)
// 		return
// 	}

// }

// func (t *ticket) Notify(c *gin.Context) {

// 	var notify model.Notification
// 	if err := c.ShouldBindJSON(&notify); err != nil {
// 		t.log.Error("failed to bind notification response body", err)
// 		c.AbortWithStatusJSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	user, err := t.module.UpdatePaymentStatus(notify.TransactionStatus, notify.SessionID)
// 	if err != nil {
// 		t.log.Error("failed to update payment status", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	pdf, err := t.module.GeneratePDFTicket(user)
// 	if err != nil {
// 		t.log.Error("Error generating PDF ticket:", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	attachmentPath := fmt.Sprintf("./public/pdfs/attachment_%s.pdf", user.ID)
// 	_, err = os.Create(attachmentPath)
// 	if err != nil {
// 		t.log.Error("error when creating attachment file for email")
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	if err := pdf.WritePdf(attachmentPath); err != nil {
// 		t.log.Error("error when copying pdf to email attachment file", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	wg := &sync.WaitGroup{}
// 	wg.Add(2)
// 	go t.sms.SendSms(user, wg)
// 	go t.email.SendEmail(user, attachmentPath, wg)
// 	wg.Wait()

// 	c.Writer.WriteHeader(http.StatusOK)
// }

func (t *ticket) GetTicket(c *gin.Context) {
	intentID := c.Param("intent_id")
	fmt.Println("going to create ticket ------>")
	pdf, err := t.tkt.GeneratePDFTicket(intentID)
	if err != nil {
		newError := err.(*model.Error)
		t.log.Error("failed to generate pdf for the given intent_id", newError)
		c.JSON(newError.ErrCode, newError)
		return
	}

	ticketPath := fmt.Sprintf("./public/pdfs/ticket_%s.pdf", intentID)
	_, err = os.Create(ticketPath)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "unable to create file to store ticket",
			RootError: err,
		}
		t.log.Error("failed to create pdf file", newError)
		c.JSON(newError.ErrCode, err)
		return
	}

	err = pdf.WritePdf(ticketPath)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to write pdf file to response body",
			RootError: err,
		}
		t.log.Error("failed to write pdf file", newError)
		c.JSON(newError.ErrCode, err)
		return
	}

	c.File(ticketPath)

	// if err := os.RemoveAll(ticketPath); err != nil {
	// 	newError := model.Error{
	// 		ErrCode: http.StatusInternalServerError,
	// 		Message: "failed to remove pdf file after written to body",
	// 	}
	// 	t.log.Error("failed to remove pdf file", newError)
	// 	c.JSON(newError.ErrCode, err)
	// 	return
	// }
}
