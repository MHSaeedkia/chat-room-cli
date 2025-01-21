package chat

type uuId string

type Client struct {
	UserName string `json:"userName"`
	UserId   uuId   `json:"userId"`
	Message  string `json:"Message"`
	Online   bool   `json:"online"`
}
