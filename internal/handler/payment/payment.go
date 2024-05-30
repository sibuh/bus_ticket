package payment

import (
	"encoding/json"
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"golang.org/x/exp/slog"
)

type payment struct {
	publishableKey string
	secretKey      string
	logger         slog.Logger
}

func Init(pkey, secretKey string, logger slog.Logger) handler.Payment {
	return &payment{
		publishableKey: pkey,
		secretKey:      secretKey,
		logger:         logger,
	}
}

func (p *payment) GetPublishableKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"publishableKey": p.publishableKey})
}

func (p *payment) HandleCreatePaymentIntent(c *gin.Context) {

	stripe.Key = p.secretKey
	fmt.Println("secret key:", p.secretKey)
	eventID, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 32)
	userID := c.Value("id").(int)

	//TODO: FETCH EVENT PRICE AND CREATE PAYMENT INTENT BASED ON THAT
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(2000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("pi.New: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clientSecret": pi.ClientSecret,
	})
}
func (p *payment) PaymentWebhook(c *gin.Context) {
	var stripeEvent stripe.Event
	if err := c.ShouldBindJSON(&stripeEvent); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusOK,
			Message:   "failed to bind request body",
			RootError: err,
		}
		p.logger.Info("unable to bind event request bosy", newError)
		c.JSON(newError.ErrCode, newError)
		return
	}
	switch stripeEvent.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(stripeEvent.Data.Raw, &paymentIntent)
		if err != nil {
			newError := model.Error{
				ErrCode:   http.StatusBadRequest,
				Message:   "failed to unmarshal event data to payment intent",
				RootError: err,
			}
			p.logger.Error("Error parsing webhook JSON", newError)
			c.JSON(newError.ErrCode, newError)
			return
		}
		p.logger.Info("PaymentIntent was successful!")
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(stripeEvent.Data.Raw, &paymentMethod)
		if err != nil {
			newError := model.Error{
				ErrCode:   http.StatusBadRequest,
				Message:   "failed to unmarshal event data to stripe paymentMethod object",
				RootError: err,
			}
			p.logger.Error("Error parsing webhook JSON", newError)
			c.JSON(newError.ErrCode, newError)
			return
		}
		p.logger.Info("PaymentMethod was attached to a Customer!")

	default:
		p.logger.Info("unhandled envet type", stripeEvent.Type)
	}

	c.JSON(http.StatusOK, nil)

}
