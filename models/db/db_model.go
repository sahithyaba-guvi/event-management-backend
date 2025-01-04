package db_model

type UserData struct {
	UserName  string `json:"userName" bson:"userName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserHash  string `json:"userHash" bson:"userHash"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

type Event struct {
	UniqueId              string               `json:"uniqueId" bson:"uniqueId"`
	EventName             string               `json:"eventName" bson:"eventName"`
	Category              string               `json:"category" bson:"category"`
	EventDescription      string               `json:"eventDescription" bson:"eventDescription"`
	SoloOrTeam            string               `json:"soloOrTeam" bson:"soloOrTeam"`
	TeamSize              int                  `json:"teamSize,omitempty" bson:"teamSize,omitempty"`
	ParticipationCapacity int                  `json:"participationCapacity" bson:"participationCapacity"`
	EventDate             int64                `json:"eventDate" bson:"eventDate"`
	EventMode             string               `json:"eventMode" bson:"eventMode"`
	EventType             string               `json:"eventType" bson:"eventType"`
	EventLocation         string               `json:"eventLocation" bson:"eventLocation,omitempty"`
	OrganizerName         string               `json:"organizerName" bson:"organizerName"`
	AdminHash             string               `json:"adminHash" bson:"adminHash"`
	PaymentType           string               `json:"paymentType" bson:"paymentType"`
	RegistrationAmount    float64              `json:"registrationAmount,omitempty" bson:"registrationAmount,omitempty"`
	FlierImage            string               `json:"flierImage" bson:"flierImage"`
	CreatedAt             int64                `json:"createdAt" bson:"createdAt"`
	Status                string               `json:"status" bson:"status"`
	UpdatedAt             int64                `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CertificationLink     string               `json:"certificationLink" bson:"certificationLink"`
	RegistrationForm      []RegisterFormFileds `json:"registrationForm" bson:"registrationForm"`
}

type RegisterFormFileds struct {
	Label string   `json:"label"`          // Label for the field (e.g., "name", "rollno")
	Type  string   `json:"type"`           // Type of the field (e.g., "textarea", "input", "select")
	Data  []string `json:"data,omitempty"` // Optional data for dropdown/select fields

}
