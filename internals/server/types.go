package server

// ContentType is a struct that holds a valid value for the "Content-Type" header
type ContentType struct {
	value string
}

func (c ContentType) String() string {
	return c.value
}

// PlainText has the value of "text/plain"
var PlainText = ContentType{"text/plain"}
