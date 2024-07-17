package ticket

import (
	"context"
	"event_ticket/internal/constant"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"net/http"

	"golang.org/x/exp/slog"
)

func Scheduler(storage storage.Session, url, sid string, logger *slog.Logger) {
	status, err := storage.GetTicketStatus(context.Background(), sid)
	if err != nil {
		return
	}
	if status == string(constant.Onhold) {
		_, err = http.Get(url + "/" + sid)
		if err != nil {
			newError := model.Error{
				ErrCode:   http.StatusInternalServerError,
				Message:   "payment status check request not successfull",
				RootError: err,
			}
			logger.Error(newError.Message, newError)

		}

	}

}
