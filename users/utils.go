package users

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

func getCreateRequestFromBody(w http.ResponseWriter, r *http.Request) (*SignUpUserRequest, error) {
	var requestData *SignUpUserRequest

	if err := common.DecodeJSONStrict(w, r, &requestData); err != nil {
		return nil, common.InvalidJSON(err)
	}

	if err := validateCreateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateCreateRequest(req *SignUpUserRequest) error {
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

func getSigninRequestFromBody(w http.ResponseWriter, r *http.Request) (*SignInUserRequest, error) {
	var requestData *SignInUserRequest

	if err := common.DecodeJSONStrict(w, r, &requestData); err != nil {
		return nil, common.InvalidJSON(err)
	}

	if err := validateSigninRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func validateSigninRequest(req *SignInUserRequest) error {
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
