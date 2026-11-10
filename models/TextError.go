package m

type TextError struct {
	Text string
}

func (te TextError) Error() string {
	return te.Text
}

func NewTextError(text string) error {
	return TextError{Text: text}
}
