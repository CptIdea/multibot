package telegram

import (
	"github.com/CptIdea/multibot"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func NewBotTG(token string) (multibot.Bot, error) {
	tg, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	messageChan := make(chan multibot.Message)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tg.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil { // ignore any non-Message Updates
				messageChan <- multibot.Message{
					Text:   update.Message.Text,
					FromID: int64(update.Message.From.ID),
					PeerID: update.Message.Chat.ID,
				}
			}

			if update.InlineQuery != nil {
				messageChan <- multibot.Message{
					Text:   update.InlineQuery.Query,
					FromID: int64(update.InlineQuery.From.ID),
				}
			}

		}
	}()

	return &tgBot{tg, messageChan}, nil
}

type tgBot struct {
	client      *tgbotapi.BotAPI
	messageChan chan multibot.Message
}

func (t *tgBot) GetMessagesChan() chan multibot.Message {
	return t.messageChan
}

func (t *tgBot) SendText(peer int, text string) error {
	_, err := t.client.Send(tgbotapi.NewMessage(int64(peer), text))
	return err
}

func (t *tgBot) SendKeyboard(peer int, text string, keyboard multibot.Keyboard) error {
	message := tgbotapi.NewMessage(int64(peer), text)
	if keyboard.GetInline() {
		rows := [][]tgbotapi.InlineKeyboardButton{}
		for _, line := range keyboard.Buttons {
			row := []tgbotapi.InlineKeyboardButton{}
			for _, button := range *line {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Payload))
			}
			rows = append(rows, row)
		}
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	} else {
		rows := [][]tgbotapi.KeyboardButton{}
		for _, line := range keyboard.Buttons {
			row := []tgbotapi.KeyboardButton{}
			for _, button := range *line {
				row = append(row, tgbotapi.NewKeyboardButton(button.Text))
			}
			rows = append(rows, row)
		}
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)
	}
	_, err := t.client.Send(message)
	return err
}
