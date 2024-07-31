package account

type service struct {
	store AccountStore
}

func NewAccountService(store AccountStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateAccount(req *CreateAccountRequest) (*CreateAccountResponse, error) {
	account := &Account{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := s.store.Create(account); err != nil {
		return nil, err
	}

	return &CreateAccountResponse{
		ID: account.ID,
	}, nil
}
