package event_response

import db_model "em_backend/models/db"

type EventRegisterFormInfoResp struct {
	FormFields db_model.RegistrationForm `json:"formFields"`
	EventName  string                    `json:"eventName"`
}
