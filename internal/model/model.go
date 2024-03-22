package model

import (
	"fmt"

	"net/mail"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
	Phone     string `form:"phone" json:"phone"`
	Email     string `form:"email" json:"email"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required.Error("first_name is required")),
		validation.Field(&u.LastName, validation.Required.Error("last_name is required")),
		validation.Field(&u.Phone, validation.Required.Error("phone is required"), validation.By(ValidatePhone)),
		validation.Field(&u.Email, validation.Required.Error("email is required"), validation.By(validateEmail)))
}
func ValidatePhone(phone any) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}

	return nil
}
func validateEmail(address any) error {
	email, ok := address.(string)
	if !ok {
		return fmt.Errorf("invalid email")
	}
	m, err := mail.ParseAddress(email)

	if err != nil || m.Address == "" {
		return fmt.Errorf("invalid email")
	}
	return nil
}

type Beneficiary struct {
	AccountNumber string  `json:"accountNumber"`
	Bank          string  `json:"bank"`
	Amount        float64 `json:"amount"`
}

type PaymentRequest struct {
	CancelURL      string        `json:"cancelUrl"`
	Nonce          string        `json:"nonce"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email"`
	ErrorURL       string        `json:"errorUrl"`
	NotifyURL      string        `json:"notifyUrl"`
	SuccessURL     string        `json:"successUrl"`
	PaymentMethods []string      `json:"paymentMethods"`
	ExpireDate     string        `json:"expireDate"`
	Items          []interface{} `json:"items"`
	Beneficiaries  []Beneficiary `json:"beneficiaries"`
	Lang           string        `json:"lang"`
}

type Session struct {
	SessionId   string  `json:"sessionId"`
	PaymentUrl  string  `json:"paymentUrl"`
	CancelUrl   string  `json:"cancelUrl"`
	TotalAmount float64 `json:"totalAmount"`
}
type CheckoutResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Data    Session `json:"data"`
}

type Notification struct {
	TransactionStatus string `json:"transactionStatus"`
	SessionID         string `json:"sessionId"`
	NotificationURL   string `json:"notificationUrl"`
}

type Sms struct {
	Token string `json:"token"`
	Phone string `json:"phone"`
	Msg   string `json:"msg"`
}
