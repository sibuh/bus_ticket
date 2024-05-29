package payment

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type payment struct {
	publishableKey string
	secretKey      string
}

func Init(pkey, secretKey string) handler.Payment {
	return &payment{
		publishableKey: pkey,
		secretKey:      secretKey,
	}
}

func (p *payment) GetPublishableKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"publishableKey": p.publishableKey})
}

func (p *payment) HandleCreatePaymentIntent(c *gin.Context) {

	var userRequest model.User
	stripe.Key = p.secretKey
	fmt.Println("secret key:", p.secretKey)

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("ShouldBindJSON: %v", err)
		return
	}

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
