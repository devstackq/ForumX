package structure

//general model
type Session struct {
	ID          int
	UUID        string
	UserID      int
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

//general global variable
var API struct {
	Authenticated bool
	Message       string
}
