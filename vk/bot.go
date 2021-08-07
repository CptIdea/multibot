package vk

import (
	"context"
	"github.com/CptIdea/multibot"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"math/rand"
	"time"
)

func NewBotVK(token string, groupID int) (multibot.Bot, error) {
	vk := api.NewVK(token)
	rand.Seed(time.Now().Unix())

	lp, err := longpoll.NewLongPoll(vk, groupID)
	if err != nil {
		return nil, err
	}

	messageChan := make(chan multibot.Message)

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		messageChan <- multibot.Message{
			Text:    obj.Message.Text,
			FromID:  int64(obj.Message.FromID),
			PeerID:  int64(obj.Message.PeerID),
			Payload: obj.Message.Payload,
		}
	})

	lp.MessageEvent(func(_ context.Context, obj events.MessageEventObject) {
		messageChan <- multibot.Message{
			Text:    string(obj.Payload),
			FromID:  int64(obj.UserID),
			PeerID:  int64(obj.PeerID),
			Payload: string(obj.Payload),
		}
	})

	go func() {
		if err := lp.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	return &vkBot{vk, messageChan}, nil
}

type vkBot struct {
	client      *api.VK
	messageChan chan multibot.Message
}

func (v *vkBot) GetMessagesChan() chan multibot.Message {
	return v.messageChan
}

func (v *vkBot) SendText(peer int, text string) error {
	_, err := v.client.MessagesSend(params.NewMessagesSendBuilder().PeerID(peer).Message(text).RandomID(rand.Int()).Params)
	return err
}

func (v *vkBot) SendKeyboard(peer int, text string, keyboard *multibot.Keyboard) error {
	kb := object.NewMessagesKeyboard(object.BaseBoolInt(keyboard.GetOnce()))
	kb.Inline = object.BaseBoolInt(keyboard.GetInline())
	for _, line := range keyboard.Buttons {
		kb.AddRow()
		for _, button := range *line {
			kb.AddTextButton(button.Text, button.Payload, GetVKColor(button.Color))
		}
	}

	_, err := v.client.MessagesSend(params.NewMessagesSendBuilder().PeerID(peer).Message(text).RandomID(rand.Int()).Keyboard(kb).Params)
	return err
}

func GetVKColor(color multibot.KeyboardButtonColor) string {
	switch color {
	case multibot.ColorNegative:
		return "negative"
	case multibot.ColorPrimary:
		return "primary"
	case multibot.ColorPositive:
		return "positive"
	default:
		return "secondary"
	}
}
