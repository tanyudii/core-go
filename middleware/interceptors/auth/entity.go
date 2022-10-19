package auth

type TokenInfoResponse struct {
	TokenInfo  *TokenInfo
	ClientInfo *ClientInfo
	Scope      string
}

type TokenInfo struct {
	UserID      string
	UserName    string
	UserEmail   string
	UserType    string
	Permissions []string
}

type ClientInfo struct {
	ClientID   string
	ClientName string
}
