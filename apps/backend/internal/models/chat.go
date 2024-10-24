package models

type Message struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	Category string `json:"category"`
}

type SendMessagePayload struct {
	Message  string `json:"message"`
	Category string `json:"category"`
}
