package account

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internals/common"
)

func GetCreateAccountRequestFromBody(w http.ResponseWriter, r *http.Request) (*CreateAccountRequest, error) {
	var requestData *CreateAccountRequest

	if err := common.DecodeJSONBody(w, r, &requestData); err != nil {
		return nil, err
	}

	if err := validateCreateAccountRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateCreateAccountRequest(req *CreateAccountRequest) error {
	var errs []string

	if req.Email == "" {
		errs = append(errs, ErrEmailRequired)
	}

	if req.Password == "" {
		errs = append(errs, ErrPasswordRequired)
	}

	if req.Name == "" {
		errs = append(errs, ErrNameRequired)
	}

	if len(errs) > 0 {
		msg := fmt.Sprintf(ErrUnableToProcessRequest, strings.Join(errs, ", "))
		return &AccountError{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
