package common

type AccountError struct {
	Status int
	Msg    string
}

func (a AccountError) Error() string {
	return a.Msg
}
