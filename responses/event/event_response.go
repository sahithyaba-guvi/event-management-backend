package event_response

import db_model "em_backend/models/db"

type EventRegisterFormInfoResp struct {
	FormFields db_model.RegistrationForm `json:"formFields"`
	EventName  string                    `json:"eventName"`
}

type GetEventByIdResp struct {
	Event             db_model.Event `json:"event"`
	IsEmailRegistered bool           `json:"isEmailRegistered"`
}
