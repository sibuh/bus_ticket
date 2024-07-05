package ticket

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/module"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ticket struct {
	log     *slog.Logger
	payment module.Payment
	tkt     module.Ticket
}

func Init(log *slog.Logger, pmt module.Payment, tkt module.Ticket) handler.Ticket {
	return &ticket{
		log:     log,
		payment: pmt,
		tkt:     tkt,
	}
}
func (t *ticket) GetTicket(c *gin.Context) {

}
