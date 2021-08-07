package multibot

const (
	ColorNegative  KeyboardButtonColor = iota
	ColorPositive  KeyboardButtonColor = iota
	ColorPrimary   KeyboardButtonColor = iota
	ColorSecondary KeyboardButtonColor = iota
)

type Keyboard struct {
	Buttons     []*KeyboardLine
	inline      bool
	once        bool
}

func (k Keyboard) GetInline() bool {
	return k.inline
}
func (k Keyboard) GetOnce() bool {
	return k.once
}

type KeyboardLine []*KeyboardButton

type KeyboardButton struct {
	Text    string
	Color   KeyboardButtonColor
	Payload string
}

type KeyboardButtonColor int

func NewKeyboard() *Keyboard {
	return &Keyboard{Buttons: []*KeyboardLine{}}
}

func (k *Keyboard) Inline() *Keyboard {
	k.inline = true
	return k
}

func (k *Keyboard) Once() *Keyboard {
	k.once = true
	return k
}

func (k *Keyboard) AddRow(line *KeyboardLine) *Keyboard {
	k.Buttons = append(k.Buttons, line)
	return k
}

func NewKeyboardLine() *KeyboardLine {
	return &KeyboardLine{}
}

func (l *KeyboardLine) AddButton(button *KeyboardButton) *KeyboardLine {
	*l = append(*l, button)
	return l
}

func NewKeyboardButton() *KeyboardButton {
	return &KeyboardButton{}
}

func (b *KeyboardButton) SetText(text string) *KeyboardButton {
	b.Text = text
	return b
}

func (b *KeyboardButton) SetPayload(payload string) *KeyboardButton {
	b.Payload = payload
	return b
}

func (b *KeyboardButton) SetColor(color KeyboardButtonColor) *KeyboardButton {
	b.Color = color
	return b
}