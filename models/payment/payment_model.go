package paymentModel

type OrderRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
	RegID    string `json:"reg_id,omitempty"` // Custom field for internal use
}

type OrderResponse struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
}

type PaymentVerification struct {
	OrderID   string `json:"razorpay_order_id"`
	PaymentID string `json:"razorpay_payment_id"`
	Signature string `json:"razorpay_signature"`
}

type OrderDetails struct {
	OrderID  string            `gorm:"primaryKey"`
	Amount   int64             `gorm:"not null"`
	Currency string            `gorm:"not null"`
	Receipt  string            `gorm:"not null"`
	Status   string            `gorm:"not null"`
	Notes    map[string]string `gorm:"type:jsonb"`
	RegID    string            `gorm:"not null"`
}

type PaymentDetails struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   string `gorm:"not null"`
	PaymentID string `gorm:"not null"`
	Status    string `gorm:"not null"`
}
