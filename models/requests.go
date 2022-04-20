package models

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	ID string `json:"id"`
}

type ErrorModel struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type GetChatHistoryRequest struct {
	UserId string `json:"user_id"`
	PeerId string `json:"peer_id"`
}

type LoginTokenResp struct {
	Code    string      `json:"code"`
	AccessToken string `json:"access_token"`
	UserId string `json:"user_id"`
}

type GetChatUsersRequest struct {
	ID string `json:"id"`
}

type GetChatUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	Username string `json:"username"`
	ID string `json:"id"`
	UnreadMessages int `json:"unread_messages"`
	LastMessage Message `json:"last_message"`
}

type Message struct {
	From string `json:"from"`
	To string `json:"to"`
	Text string `json:"text"`
	CreatedAt string `json:"created_at"`
	ID string `json:"id"`
	Read bool `json:"read"`
}

type CreateMessageResponse struct {
	ID string `json:"id"`
}

type DateMessage struct {
	Date string `json:"date"`
	Messages []Message `json:"messages"`
}