package initiator

import (
	"context"
	"event_ticket/internal/handler/payment"
	huser "event_ticket/internal/handler/user"
	"event_ticket/internal/middleware"
	mpayment "event_ticket/internal/module/payment"
	"event_ticket/internal/utils/token/paseto"

	hevnt "event_ticket/internal/handler/event"
	muser "event_ticket/internal/module/user"
	"event_ticket/internal/routing"
	"event_ticket/internal/storage/event"

	"event_ticket/internal/storage/user"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Initiate() {
	logger := InitLogger()
	InitConfig("config", logger)
	server := gin.Default()
	server.Use(middleware.Cors())
	v1 := server.Group("v1")
	logger.Info("initiate database")
	queries := InitDB(viper.GetString("dbConn"))
	logger.Info("intiating storage layer")
	storage := NewStorage(user.Init(logger, queries), event.Init(&logger, queries))
	maker := paseto.NewPasetoMaker(viper.GetString("token.key"), viper.GetDuration("token.duration")*time.Second)
	mware := middleware.NewMiddleware(logger, maker, storage.user)
	module := NewModule(
		muser.Init(
			logger,
			storage.user,
			maker,
		),
		storage.event,
		mpayment.Init(&logger, storage.event),
	)

	handler := InitHandler(
		huser.Init(logger, module.user),
		payment.Init(
			viper.GetString("payment.publishable_key"),
			viper.GetString("payment.secret_key"),
			logger, module.payment),
		hevnt.Init(&logger, module.event),
	)
	routing.InitRouter(v1, handler.user, handler.payment, handler.event, mware)
	srv := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		ReadTimeout: viper.GetDuration("server.read_time_out") * time.Second,
		Handler:     server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		fmt.Println("server starting at ", viper.GetString("server.port"))
		srv.ListenAndServe()
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	logger.Warn("sever is going to shut down %+V", srv.Shutdown(ctx))

}
