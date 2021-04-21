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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	CreateUserResponse struct {
		Ok string `json:"ok"`
	}

	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		User User  `json:"user"`
		Err  error `json:"err,omitempty"`
	}
)

// This function is responsible for encoding the Response body into JSON
func encodeResponse(ctx context.Context, respWriter http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(respWriter).Encode(response)
}

// Responsible for decoding the Request body for Creating a User and returning an
// unmarshaled request body in the CreateUserRequest format.
// It takes the actual HTTP Request and returns an empty interface with the underlying
// type/value being the CreateUserRequest struct, thus "casting" it from one to the other.
func decodeUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	// Declare a variable of CreateUserRequest struct type
	var userReq CreateUserRequest
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
func decodeGetUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var userReq GetUserRequest
	// Parse the Id passed in as a URL variable
	pathVars := mux.Vars(req)

	// Init the GetUserRequest struct with the Id we gleaned from the URL
	// path variable
	userReq = GetUserRequest{
		Id: pathVars["id"],
	}

	return userReq, nil

}
