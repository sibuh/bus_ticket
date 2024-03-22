package initiator

import (
	"context"
	"event_ticket/internal/handler/ticket"
	"event_ticket/internal/module/email"
	"event_ticket/internal/module/sms"
	module "event_ticket/internal/module/ticket"
	"event_ticket/internal/routing"
	store "event_ticket/internal/storage/ticket"
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
	server.LoadHTMLGlob("public/html/*.html")
	server.Static("/public", "./public")
	storage := store.Init(logger, fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")))
	mod := module.Init(
		logger,
		viper.GetString("payment.cancel_url"),
		viper.GetString("payment.error_url"),
		viper.GetString("payment.notify_url"),
		viper.GetString("payment.account_number"),
		viper.GetString("payment.bank"),
		viper.GetString("payment.session_url"),
		viper.GetString("payment.success_url"),
		viper.GetString("payment.api_key"),
		viper.GetFloat64("payment.item_price"),
		viper.GetFloat64("payment.amount"),
		storage,
		viper.GetDuration("payment.expire_date"),
	)
	sms := sms.Init(logger, viper.GetString("sms.token"), viper.GetString("sms.url"), viper.GetString("sms.template"))
	email := email.Init(logger, viper.GetString("email.host"), viper.GetString("email.user_name"), viper.GetString("email.password"), viper.GetString("email.subject"))
	handler := ticket.Init(logger, viper.GetString("payment.error_url"), mod, sms, email)
	routing.InitRouter(v1, handler)
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
