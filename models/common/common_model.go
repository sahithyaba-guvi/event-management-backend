package common_model

type Authtoken struct {
	AuthToken string `json:"authToken"`
}
type Admins struct {
	Admin []string `bson:"admins"`
}
