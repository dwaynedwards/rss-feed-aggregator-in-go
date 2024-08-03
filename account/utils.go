package account

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

func getCreateRequestFromBody(w http.ResponseWriter, r *http.Request) (*CreateAccountRequest, error) {
	var requestData *CreateAccountRequest

	if err := common.DecodeJSONStrict(w, r, &requestData); err != nil {
		return nil, common.InvalidJSON(err)
	}

	if err := validateCreateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateCreateRequest(req *CreateAccountRequest) error {
	errs := map[string]string{}

	if req.Email == "" {
		errs["email"] = ErrEmailRequired
	}

	if req.Password == "" {
		errs["password"] = ErrPasswordRequired
	}

	if req.Name == "" {
		errs["name"] = ErrNameRequired
	}

	if len(errs) > 0 {
		return common.InvalidRequestData(errs)
	}

	return nil
}

func getSigninRequestFromBody(w http.ResponseWriter, r *http.Request) (*SigninAccountRequest, error) {
	var requestData *SigninAccountRequest

	if err := common.DecodeJSONStrict(w, r, &requestData); err != nil {
		return nil, common.InvalidJSON(err)
	}

	if err := validateSigninRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateSigninRequest(req *SigninAccountRequest) error {
	errs := map[string]string{}

	if req.Email == "" {
		errs["email"] = ErrEmailRequired
	}

	if req.Password == "" {
		errs["password"] = ErrPasswordRequired
	}

	if len(errs) > 0 {
		return common.InvalidRequestData(errs)
	}

	return nil
}
