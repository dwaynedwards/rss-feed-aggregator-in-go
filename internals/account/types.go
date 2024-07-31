package account

type Account struct {
	ID       int
	Email    string
	Password string
	Name     string
}

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateAccountResponse struct {
	ID int `json:"id"`
}

type AccountError struct {
	Status int
	Msg    string
}

func (a AccountError) Error() string {
	return a.Msg
}

type AccountService interface {
	CreateAccount(*CreateAccountRequest) (*CreateAccountResponse, error)
}

type AccountStore interface {
	Create(*Account) error
}

type inMemoryAccountDB map[string]Account

const (
	ErrEmailRequired          = "email is a required field"
	ErrPasswordRequired       = "password is a required field"
	ErrNameRequired           = "name is a required field"
	ErrUnableToProcessRequest = "unable to process request body: %s"
)
