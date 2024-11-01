package models

type Message struct {
	ChatId   string `json:"chat_id"`
	UserId   string `json:"user_id"`
	Message  string `json:"message"`
	Category string `json:"category"`
	AiModel  string `json:"ai_model"`
}

type Chat struct {
	Id     string `json:"chat_id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type SendMessagePayload struct {
	ChatId   string `json:"chat_id"`
	Message  string `json:"message"`
	Category string `json:"category"`
}

type SendFirstMessagePayload struct {
	Message  string `json:"message"`
	Category string `json:"category"`
	Name     string `json:"name"`
}
