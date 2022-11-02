package auth

type TokenInfoResponse struct {
	TokenInfo  *TokenInfo
	ClientInfo *ClientInfo
	Scope      string
}

type TokenInfo struct {
	UserID         string
	UserName       string
	UserEmail      string
	UserSerial     string
	UserType       string
	CompanyID      string
	CompanySerial  string
	CompanyName    string
	Permissions    []string
	IsInternalCall bool
}

type ClientInfo struct {
	ClientID   string
	ClientName string
}
