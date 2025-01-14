package common_responses

type SuccessResponse struct {
	Access  bool        `json:"access"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // Dynamic key for additional data
}

type FailureResponse struct {
	Access  bool   `json:"access"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Access  bool   `json:"access"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type LoginDetails struct {
	UserName  string `json:"userName" bson:"userName"`
	Email     string `json:"email"`
	Hash      string `json:"hash"`
	AuthToken string `json:"authToken" bson:"authToken"`
	Access    string `json:"access"`
	IsAdmin   bool   `json:"isAdmin"`
}
