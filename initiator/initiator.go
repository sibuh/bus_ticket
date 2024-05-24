package initiator

import (
	"context"
	"event_ticket/internal/handler/payment"
	huser "event_ticket/internal/handler/user"

	muser "event_ticket/internal/module/user"
	"event_ticket/internal/routing"
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
	v1 := server.Group("v1")
	logger.Info("initiate database")
	queries := InitDB(viper.GetString("dbConn"))
	logger.Info("intiating storage layer")
	storage := NewStorage(user.Init(logger, queries))
	module := NewModule(muser.Init(logger, storage.user))

	//sms := sms.Init(logger, viper.GetString("sms.token"), viper.GetString("sms.url"), viper.GetString("sms.template"))
	//email := email.Init(logger, viper.GetString("email.host"), viper.GetString("email.user_name"), viper.GetString("email.password"), viper.GetString("email.subject"))

	//handler := ticket.Init(logger, viper.GetString("payment.error_url"), mod, sms, email)
	handler := InitHandler(huser.Init(logger, module.user), payment.Init(viper.GetString("payment.publishable_key"), viper.GetString("payment.secret_key")))
	routing.InitRouter(v1, handler.user, handler.payment)
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
