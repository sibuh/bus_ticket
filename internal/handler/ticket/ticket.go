package ticket

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/module"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ticket struct {
	log *slog.Logger
	tkt module.Ticket
}

func Init(log *slog.Logger, tkt module.Ticket) handler.Ticket {
	return &ticket{
		log: log,
		tkt: tkt,
	}
}
func (t *ticket) GetTicket(c *gin.Context) {

}
