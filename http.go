package accountsrv

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

/*
This can be thought of as the Transport Layer of our service. It deals with the HTTP
specific functionality of directing HTTP endpoint traffic to our Endpoints (which then
safely expose the Service methods according to the Transport type, in this case HTTP.)

In other words, this maps the HTTP-world stuff to a corresponding Endpoint, where
the endpoint can then take the next steps to do some safe "middleman" work to safely "transform"
the Transport request into a Service-friendly format to then safely invoke the Service's methods
*/

// Factory function for creating an HTTP server that will specifically handle the HTTP
// traffic for our Account service. We map HTTP methods and routes to the respective
// Endpoint, where the Endpoint will then broker that HTTP request to the Service
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	// Init the router
	router := mux.NewRouter()
	// Have the router use the middleware we defined, in this case it simply
	// adds the content-type:application/json header to each of our responses.
	router.Use(commonMiddleware)

	// TODO: Subrouting

	router.Methods("GET").Path("/users/{id}").Handler(
		httptransport.NewServer(
			endpoints.GetUser,
			DecodeGetUserReq,
			EncodeResponse,
		))

	router.Methods("PUT").Path("/users/{id}/profile").Handler(
		httptransport.NewServer(
			endpoints.UpdateUserProfile,
			DecodeUpdateUserProfileReq,
			EncodeResponse,
		))

	router.Methods("POST").Path("/orgs").Handler(
		httptransport.NewServer(
			endpoints.CreateOrg,
			DecodeCreateOrgReq,
			EncodeResponse,
		))

	// Instead of passing in the Endpoint directly, we instead
	// do a bit of functional programming-esque stuff where we pass in another function to handler
	// that itself consumes ANOTHER function (the actual endpoint which is itself a function that returns a function)
	router.Methods("POST").Path("/orgs/{org_id}/login").Handler(
		httptransport.NewServer(
			endpoints.LoginUser,
			DecodeLoginReq,
			EncodeResponse,
		))

	router.Methods("POST").Path("/orgs/{org_id}/users").Handler(
		httptransport.NewServer( // This "brokers" the Transport and Endpoint layers, allowing us a way to transform the request from on layer to the next
			endpoints.CreateUser, // The endpoint itself
			DecodeCreateUserReq,  // How we want to "decode" the request, i.e. take the HTTP request and cast it into something the Endpoint/service can use
			EncodeResponse,       // How we want to "encode" the resulting response our Endpoint returns (this case as JSON)
		))

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		respWriter.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(respWriter, req)
	})
}

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
