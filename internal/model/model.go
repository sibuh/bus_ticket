package model

import (
	"fmt"
	"time"

	"net/mail"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID        int32     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (cur CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&cur,
		validation.Field(&cur.FirstName, validation.Required.Error("first_name is required")),
		validation.Field(&cur.LastName, validation.Required.Error("last_name is required")),
		validation.Field(&cur.Phone, validation.Required.Error("phone is required"), validation.By(ValidatePhone)),
		validation.Field(&cur.Email, validation.Required.Error("email is required"), validation.By(validateEmail)),
		validation.Field(&cur.Password, validation.Required.Error("password is required"),
			validation.Length(6, 8).Error("password legth should be b/n 6 and 8.")),
		validation.Field(&cur.Username, validation.Required.Error("username is required")),
	)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr LoginRequest) Validate() error {
	return validation.ValidateStruct(&lr,
		validation.Field(&lr.Password, validation.Required.Error("password is required"),
			validation.Length(6, 8).Error("password legth should be b/n 6 and 8.")),
		validation.Field(&lr.Username, validation.Required.Error("username is required")))
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

//	type Session struct {
//		SessionId   string  `json:"sessionId"`
//		PaymentUrl  string  `json:"paymentUrl"`
//		CancelUrl   string  `json:"cancelUrl"`
//		TotalAmount float64 `json:"totalAmount"`
//	}
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

type Error struct {
	ErrCode   int    `json:"err_code"`
	Message   string `json:"message"`
	RootError error  `json:"root_error"`
}

func (e *Error) Error() string {
	return e.Message
}

type Event struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int32     `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Price       float64   `json:"price"`
	CreateAt    time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type CreateIntentParam struct {
	UserID   int32  `json:"user_id"`
	EventID  int32  `json:"event_id"`
	IntentID string `json:"intent_id"`
}

type Payment struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"user_id"`
	EventID       int32     `json:"event_id"`
	PaymentStatus string    `json:"payment_status"`
	IntentID      string    `json:"intent_id"`
	CheckInStatus string    `json:"check_in_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type Ticket struct {
	ID       string `mapstructure:"id" json:"id"`
	TripID   int32  `mapstructure:"trip_id" json:"trip_id"`
	TicketNo int32  `mapstructure:"ticket_no" json:"ticket_no"`
	BusNo    int32  `mapstructure:"bus_no" json:"bus_no"`
	Status   string `mapstructure:"status" json:"status"`
}
type ReserveTicketRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type Session struct {
	ID            string    `json:"id"`
	Tkt           Ticket    `json:"tkt"`
	PaymentStatus string    `json:"payment_status"`
	PaymentUrl    string    `json:"paymentUrl"`
	CancelUrl     string    `json:"cancelUrl"`
	TotalAmount   float64   `json:"totalAmount"`
	CreatedAt     time.Time `json:"created_at"`
}
