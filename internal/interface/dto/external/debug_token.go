package external

type DebugTokenRequest struct {
	InputToken  string `param:"input_token"`
	AccessToken string `param:"access_token"`
}

type DebugTokenResponse struct {
	Data struct {
		AppID               string   `json:"app_id"`
		Type                string   `json:"type"`
		Application         string   `json:"application"`
		ExpiresAt           int64    `json:"expires_at"`
		DataAccessExpiresAt int64    `json:"data_access_expires_at"`
		IsValid             bool     `json:"is_valid"`
		UserID              string   `json:"user_id"`
		Scopes              []string `json:"scopes"`
	} `json:"data"`
}
