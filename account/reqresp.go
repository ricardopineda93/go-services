package account

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Defining the "shapes" of the requests and responses to cast and
// decode/encode to these respective shapes
type (
	CreateUserRequest struct {
		OrgID     string `json:"org_id"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		OrgType   string `json:"org_type"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
	}
	CreateUserResponse struct {
		ID string `json:"id"`
	}

	GetUserRequest struct {
		ID string `json:"id"`
	}
	GetUserAccountResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		JoinedOn string `json:"joined_on"`
	}

	LoginRequest struct {
		OrgID    string
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Account LoginAccount `json:"user_account"`
		Profile LoginProfile `json:"user_profile"`
	}
	LoginAccount struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		JoinedOn string `json:"joined_on"`
	}
	LoginProfile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		LastLogin string `json:"last_login"`
	}

	UpdateProfileRequest struct {
		AccountID string
		Updates   ProfileUpdates `json:"profile_updates"`
	}
	UpdateProfileResponse struct {
		OK string `json:"ok"`
	}
	ProfileUpdates struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
	}

	CreateOrgRequest struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
		Timezone string `json:"timezone"`
		Website  string `json:"website"`
	}
	CreateOrgResponse struct {
		ID string `json:"id"`
	}
)

// This function is responsible for encoding the Response body into JSON
func EncodeResponse(ctx context.Context, respWriter http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(respWriter).Encode(response)
}

// Responsible for decoding the Request body for Creating a User and returning an
// unmarshaled request body in the CreateUserRequest format.
// It takes the actual HTTP Request and returns an empty interface with the underlying
// type/value being the CreateUserRequest struct, thus "casting" it from one to the other.
func DecodeCreateUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	pathVars := mux.Vars(req)
	// Declare a variable of CreateUserRequest struct type
	var userReq CreateUserRequest

	userReq.OrgID = pathVars["org_id"]
	// Decode and "cast" the values of the request body to the variable declared above
	err := json.NewDecoder(req.Body).Decode(&userReq)
	if err != nil {
		return nil, err
	}
	// When returning here, because this function returns an empty interface, the interface's
	// underlying "concrete type & value" is this struct -- in other words we are NOT
	// returning this struct as is but instead an interface that contains this struct
	// A bit roundabout since in the endpoint function we then cast the CreateUserRequest struct
	// from this empty interface we return here to get the underlying CreateUserRequest struct
	// but, yknow, whatever.
	return userReq, nil
}

// Responsible for taking the Get User request and returning a GetUserRequest struct
// formatted object
func DecodeGetUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var userReq GetUserRequest
	// Parse the Id passed in as a URL variable
	pathVars := mux.Vars(req)

	// Init the GetUserRequest struct with the Id we gleaned from the URL
	// path variable
	userReq = GetUserRequest{
		ID: pathVars["id"],
	}

	return userReq, nil

}

func DecodeLoginReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var loginReq LoginRequest

	pathVars := mux.Vars(req)

	loginReq.OrgID = pathVars["org_id"]
	// Decode and "cast" the values of the request body to the variable declared above
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		return nil, err
	}
	return loginReq, nil
}

func DecodeUpdateUserProfileReq(ctx context.Context, req *http.Request) (interface{}, error) {
	pathVars := mux.Vars(req)
	var updatesReq UpdateProfileRequest

	err := json.NewDecoder(req.Body).Decode(&updatesReq.Updates)

	if err != nil {
		return updatesReq, err
	}

	updatesReq.AccountID = pathVars["id"]

	return updatesReq, nil
}

func DecodeCreateOrgReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var orgReq CreateOrgRequest

	err := json.NewDecoder(req.Body).Decode(&orgReq)
	if err != nil {
		return nil, err
	}

	return orgReq, nil
}
