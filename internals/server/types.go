package server

// ContentType is a struct that holds a valid value for the "Content-Type" header
type ContentType struct {
	value string
}

func (c ContentType) String() string {
	return c.value
}

// ContentTypePlainText has the value of "text/plain"
var (
	ContentTypePlainText = ContentType{"text/plain"}
	ContentTypeJSON      = ContentType{"application/json"}
)
