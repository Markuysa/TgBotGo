package telegram

import "TelegramBot/events"

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int                     `json:"update_id"`
	Message *events.IncomingMessage `json:"message"`
}
