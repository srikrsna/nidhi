package nidhi

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	InvalidCursor Error = "invalid cursor"
)
