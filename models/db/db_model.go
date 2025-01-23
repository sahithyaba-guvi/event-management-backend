package db_model

type UserData struct {
	UserName  string `json:"userName" bson:"userName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserHash  string `json:"userHash" bson:"userHash"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

type Event struct {
	UniqueId  string `json:"uniqueId,omitempty" bson:"uniqueId"`
	EventName string `json:"eventName" bson:"eventName"`
	Category  string `json:"category" bson:"category"`
	// CategoryId                string                     `json:"categoryId" bson:"categoryId"`
	EventDescription        string   `json:"eventDescription" bson:"eventDescription"`
	EventType               string   `json:"eventType" bson:"eventType"`
	EventMode               string   `json:"eventMode" bson:"eventMode"`
	EventLocation           string   `json:"eventLocation,omitempty" bson:"eventLocation"`
	EventDate               int64    `json:"eventDate" bson:"eventDate"`
	FlierImage              string   `json:"flierImage" bson:"flierImage"`
	PaymentType             string   `json:"paymentType" bson:"paymentType"`
	TicketComboDetails      []string `json:"ticketComboDetails" bson:"ticketComboDetails"`
	ParticipationGuidelines string   `json:"participationGuidelines,omitempty" bson:"participationGuidelines"`
	// RegistrationLimit         int                        `json:"registrationLimit" bson:"registrationLimit"`
	RegistrationDetailsFormId string                     `json:"registrationDetailsFormId,omitempty" bson:"registrationDetailsFormId"`
	RegistrationCount         int                        `json:"registrationCount,omitempty" bson:"registrationCount"`
	RegistrationData          []RegistrationPricingCombo `json:"registrationData,omitempty" bson:"-"`
	RegistrationForm          []RegisterFormFields       `json:"registrationForm,omitempty" bson:"-"`
	CreatedAt                 int64                      `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt                 int64                      `json:"updatedAt,omitempty" bson:"updatedAt"`
	Status                    string                     `json:"status,omitempty" bson:"status"`
}
type RegistrationPricingCombo struct {
	TicketType int32   `json:"ticketType"`
	Price      float64 `json:"price"`
}

type RegisterFormFields struct {
	Label string   `json:"label"`
	Type  string   `json:"type"`
	Data  []string `json:"data,omitempty"`
}

type RegistrationForm struct {
	RegistrationFormId     string               `json:"registrationFormId" bson:"registrationFormId"`
	EventId                string               `json:"eventId" bson:"eventId"`
	RegistrationFormFields []RegisterFormFields `json:"registrationFormFields" bson:"registrationFormFields"`
}

type TeamMemberDetail struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

type Registration struct {
	RegisterId        string             `json:"registerId" bson:"registerId"`
	EventId           string             `json:"eventId" bson:"eventId"`
	TeamSize          int                `json:"teamSize" bson:"teamSize"`
	TeamMemberDetails []TeamMemberDetail `json:"teamMemberDetails" bson:"teamMemberDetails"`
	PrimaryEmailId    string             `json:"primaryEmailId" bson:"primaryEmailId"`
	CreatedAt         int64              `json:"createdAt" bson:"createdAt"`
	UpdatedAt         int64              `json:"updatedAt" bson:"updatedAt"`
}
