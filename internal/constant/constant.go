package constant

type TicketStatus string

const (
	Reserved TicketStatus = "Reserved"
	Free     TicketStatus = "Free"
	Onhold   TicketStatus = "Onhold"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "Pending"
)
