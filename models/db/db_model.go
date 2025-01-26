package db_model

type UserData struct {
	UserName  string `json:"userName" bson:"userName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserHash  string `json:"userHash" bson:"userHash"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

type Event struct {
	UniqueId                  string                   `json:"uniqueId,omitempty" bson:"uniqueId"`
	EventName                 string                   `json:"eventName" bson:"eventName"`
	Category                  string                   `json:"category" bson:"category"`
	EventDescription          string                   `json:"eventDescription" bson:"eventDescription"`
	EventType                 string                   `json:"eventType" bson:"eventType"`
	EventMode                 string                   `json:"eventMode" bson:"eventMode"`
	EventLocation             string                   `json:"eventLocation,omitempty" bson:"eventLocation"`
	EventDate                 int64                    `json:"eventDate" bson:"eventDate"`
	FlierImage                string                   `json:"flierImage" bson:"flierImage"`
	PaymentType               string                   `json:"paymentType" bson:"paymentType"`
	ComboPrices               RegistrationPricingCombo `json:"comboPrices,omitempty" bson:"comboPrices"`
	ParticipationGuidelines   string                   `json:"participationGuidelines,omitempty" bson:"participationGuidelines"`
	RegistrationDetailsFormId string                   `json:"registrationDetailsFormId,omitempty" bson:"registrationDetailsFormId"`
	RegistrationCount         int                      `json:"registrationCount,omitempty" bson:"registrationCount"`
	TicketCombos              []int                    `json:"ticketCombos,omitempty" bson:"ticketCombos"`
	CreatedAt                 int64                    `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt                 int64                    `json:"updatedAt,omitempty" bson:"updatedAt"`
	Status                    string                   `json:"status,omitempty" bson:"status"`
}

type RegistrationPricingCombo struct {
	RegistrationAmount float64 `json:"registrationAmount" bson:"registrationAmount"`
	Combo5Price        string  `json:"combo5Price" bson:"combo5Price"`
	Combo10Price       string  `json:"combo10Price" bson:"combo10Price"`
}

type RegisterFormFields struct {
	Label string   `json:"label" bson:"label"`
	Type  string   `json:"type" bson:"type"`
	Data  []string `json:"data,omitempty" bson:"data"`
}

type RegistrationForm struct {
	RegistrationFormId string `json:"registrationFormId" bson:"registrationFormId"`
	EventId            string `json:"eventId" bson:"eventId"`
	// RegistrationFormFields RegisterFormFields   `json:"registrationFormFields" bson:"registrationFormFields"`
	PrimaryMemberForm []RegisterFormFields `json:"primaryMemberForm,omitempty" bson:"primaryMemberForm"`
	TeamDetailsForm   []RegisterFormFields `json:"teamDetailsForm,omitempty" bson:"teamDetailsForm"`
}

type EventPayload struct {
	EventName               string                   `json:"eventName"`
	Category                string                   `json:"category"`
	EventDescription        string                   `json:"eventDescription"`
	EventType               string                   `json:"eventType"`
	EventMode               string                   `json:"eventMode"`
	EventLocation           string                   `json:"eventLocation"`
	EventDate               int64                    `json:"eventDate"`
	FlierImage              string                   `json:"flierImage"`
	PaymentType             string                   `json:"paymentType"`
	ComboPrices             RegistrationPricingCombo `json:"comboPrices"`
	ParticipationGuidelines string                   `json:"participationGuidelines"`
	TicketCombos            []int                    `json:"ticketCombos"`
	PrimaryMemberForm       []RegisterFormFields     `json:"primaryMemberForm"`
	TeamDetailsForm         []RegisterFormFields     `json:"teamDetailsForm"`
	RegistrationForm        RegistrationForm         `json:"registrationForm"`
}

type RegistrationData struct {
	UniqueId                     string               `json:"uniqueId" bson:"uniqueId"`
	PrimaryMemberForm            []RegisterFormFields `json:"primaryMemberForm,omitempty" bson:"primaryMemberForm"`
	TeamDetailsForm              []RegisterFormFields `json:"teamDetailsForm,omitempty" bson:"teamDetailsForm"`
	RegistrationId               string               `json:"registrationId"`
	PrimaryEmailId               string               `json:"primaryEmailId"`
	RegisteredAt                 int64                `json:"registeredAt"`
	IsTicketVerified             bool                 `json:"isTicketVerified"`
	TicketVerificationStatusTeam []string             `json:"ticketVerificationStatusTeam"`
	QrCode                       string               `json:"qrCode"`
}

type RegisterReq struct {
	PrimaryMemberForm []RegisterFormFields `json:"primaryMemberForm,omitempty" bson:"primaryMemberForm"`
	TeamDetailsForm   []RegisterFormFields `json:"teamDetailsForm,omitempty" bson:"teamDetailsForm"`
	UniqueId          string               `json:"uniqueId" bson:"uniqueId"`
}

type RegistrationRequestData struct {
	RegistrationId   string `json:"registrationId"`
	PrimaryEmailId   string `json:"primaryEmailId"`
	QrCode           string `json:"qrCode"`
	IsTicketVerified bool   `json:"isTicketVerified"`
	CreatedAt        int64  `json:"createdAt"`
	UpdatedAt        int64  `json:"updatedAt"`
}
