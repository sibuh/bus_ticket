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

	if err := os.RemoveAll(ticketPath); err != nil {
		newError := model.Error{
			ErrCode: http.StatusInternalServerError,
			Message: "failed to remove pdf file after written to body",
		}
		t.log.Error("failed to remove pdf file", newError)
		c.JSON(newError.ErrCode, err)
		return
	}
}
