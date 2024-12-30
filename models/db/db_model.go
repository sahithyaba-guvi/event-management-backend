package db_model

type UserData struct {
	UserName  string `json:"userName" bson:"userName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserHash  string `json:"userHash" bson:"userHash"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

type Event struct {
	UniqueId         string  `json:"uniqueId" bson:"uniqueId"`
	EventName        string  `json:"eventName" bson:"eventName"`
	Category         string  `json:"category" bson:"category"`
	EventDescription string  `json:"eventDescription" bson:"eventDescription"`
	SoloOrTeam       string  `json:"soloOrTeam" bson:"soloOrTeam"`
	TeamSize         int     `json:"teamSize,omitempty" bson:"teamSize,omitempty"`
	ParticipantLimit int     `json:"participantLimit" bson:"participantLimit"`
	DateOfEvent      int64   `json:"dateOfEvent" bson:"dateOfEvent"`
	OnlineOffline    string  `json:"onlineOffline" bson:"onlineOffline"`
	Location         string  `json:"location" bson:"location,omitempty"`
	OrganizerName    string  `json:"organizerName" bson:"organizerName"`
	AdminHash        string  `json:"adminHash" bson:"adminHash"`
	FreePaid         string  `json:"freePaid" bson:"freePaid"`
	Amount           float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	// EventThumbnail   string  `json:"eventThumbnail" bson:"eventThumbnail"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	Status    string `json:"status" bson:"status"`
	UpdatedAt int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
