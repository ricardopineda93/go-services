package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// These are akin to Controller Handlers, this is where we can "unpack" the requests
// from the Transport layer and begin to expose the Service's functionalities safely
// to service the requests

// Define the struct that will contain the endpoint functions that expose the service.
// In this particular case, they will interface with an HTTP server
type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
}

// Factory function that exposes this service-specific functionalities
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
	}
}

// Function returns a function (satisfying the Endpoint spec) that will be called upon
// by the server. In this case this returns a function that exposes/calls upon the Service's
// method of creating a User
func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)                    // Assert that the underlying type of the interface we received is of CreateUserRequest, and if it is extract the value
		ok, err := s.CreateUser(ctx, req.Email, req.Password) // Call the service method
		return CreateUserResponse{Ok: ok}, err                // Return a response in the shape we specified in this struct
	}
}

func makeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{User: user, Err: err}, err
	}
}