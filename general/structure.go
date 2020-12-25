package general

//general structure -> for child packages use
type Session struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	UserID      int    `json:"userId"`
	DbCookie    string `json:"db_cookie"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

//general global variable
var API struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message"`
}
