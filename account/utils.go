package account

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

func getCreateRequestFromBody(w http.ResponseWriter, r *http.Request) (*CreateAccountRequest, error) {
	var requestData *CreateAccountRequest

	if err := common.DecodeJSONBody(w, r, &requestData); err != nil {
		return nil, err
	}

	if err := validateCreateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateCreateRequest(req *CreateAccountRequest) error {
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
		return &common.AccountError{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}

func getSigninRequestFromBody(w http.ResponseWriter, r *http.Request) (*SigninAccountRequest, error) {
	var requestData *SigninAccountRequest

	if err := common.DecodeJSONBody(w, r, &requestData); err != nil {
		return nil, err
	}

	if err := validateSigninRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateSigninRequest(req *SigninAccountRequest) error {
	var errs []string

	if req.Email == "" {
		errs = append(errs, ErrEmailRequired)
	}

	if req.Password == "" {
		errs = append(errs, ErrPasswordRequired)
	}

	if len(errs) > 0 {
		msg := fmt.Sprintf(ErrUnableToProcessRequest, strings.Join(errs, ", "))
		return &common.AccountError{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
