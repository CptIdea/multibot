package multibot

type Bot interface {
	SendText(peer int,text string)error
	SendKeyboard(peer int,text string, keyboard Keyboard)error
	GetMessagesChan() chan Message
}
