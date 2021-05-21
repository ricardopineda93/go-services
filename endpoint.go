package accountsrv

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
	CreateUser        endpoint.Endpoint
	GetUser           endpoint.Endpoint
	LoginUser         endpoint.Endpoint
	UpdateUserProfile endpoint.Endpoint

	CreateOrg endpoint.Endpoint
}

// Factory function that exposes this service-specific functionalities
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:        makeCreateUserEndpoint(s),
		GetUser:           makeGetUserAccountEndpoint(s),
		LoginUser:         makeLoginUserEndpoint(s),
		UpdateUserProfile: makeUpdateUserProfileEndpoint(s),

		CreateOrg: makeCreateOrgEndpoint(s),
	}
}

// Function returns a function (satisfying the Endpoint spec) that will be called upon
// by the server. In this case this returns a function that exposes/calls upon the Service's
// method of creating a User
func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest) // Assert that the underlying type of the interface we received is of CreateUserRequest, and if it is extract the value
		// TODO: Have some validation somewhere that can prevent us from wasting effort
		// if the request isn't going to be valid.
		id, err := s.CreateUser(ctx, req.OrgID, req.Username, req.Password, req.OrgType, req.FirstName, req.LastName, req.Email, req.Phone)
		return CreateUserResponse{ID: id, Err: err}, nil // Return a response in the shape we specified in this struct
	}
}

func makeGetUserAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		account, err := s.GetUserAccount(ctx, req.ID)
		return GetUserAccountResponse{UserAccount: account, Err: err}, nil
	}
}

func makeLoginUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		loginDetails, err := s.Login(ctx, req.OrgID, req.Username, req.Password)

		return LoginResponse{
			LoginDetails: loginDetails,
			Err:          err,
		}, nil

	}
}

func makeUpdateUserProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateProfileRequest)

		// Unmarshall Updates struct into map based on json tags
		updatesMap := structToMapByTag(req.Updates, "json")

		err := s.UpdateUserProfile(ctx, req.AccountID, updatesMap)

		return UpdateProfileResponse{OK: "ok", Err: err}, nil
	}
}

func makeCreateOrgEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrgRequest)

		id, err := s.CreateOrg(ctx, req.Name, req.Type, req.Phone, req.Address, req.Timezone, req.Website)

		return CreateOrgResponse{ID: id, Err: err}, nil
	}
}
