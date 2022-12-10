package telegram

import (
	"TelegramBot/clients/telegram"
	"TelegramBot/events"
	"TelegramBot/libs/e"
	"TelegramBot/storage"
	"errors"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't fetch", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	result := make([]events.Event, 0, len(updates))

	for _, update := range updates {

		result = append(result, toEvent(update))
	}
	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
	return nil
}

func (p *Processor) processMessage(event events.Event) error {

	meta, err := p.meta(event)
	if err != nil {
		return e.Wrap("can't process the message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process the message", err)
	}

	return nil
}

func (p *Processor) meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func toEvent(update telegram.Update) events.Event {
	updateType := fetchType(update)
	res := events.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}
	if updateType == events.Message {
		res.Meta = Meta{
			ChatId:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}
	return res
}

type Meta struct {
	ChatId   int
	Username string
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
