package payment

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
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
	c.JSON(http.StatusOK, p.publishableKey)
}

func (p *payment) HandleCreatePaymentIntent(c *gin.Context) {

	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var userRequest model.User
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("ShouldBindJSON: %v", err)
		return
	}

	// data := db.GetAProduct(product.Id)

	// Create a PaymentIntent with amount and currency
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
