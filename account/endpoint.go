package account

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-kit/kit/endpoint"
	u "github.com/rjjp5294/gokit-tutorial/account/user"
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
}

// Factory function that exposes this service-specific functionalities
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:        makeCreateUserEndpoint(s),
		GetUser:           makeGetUserAccountEndpoint(s),
		LoginUser:         makeLoginUserEndpoint(s),
		UpdateUserProfile: makeUpdateUserProfileEndpoint(s),
	}
}

// Function returns a function (satisfying the Endpoint spec) that will be called upon
// by the server. In this case this returns a function that exposes/calls upon the Service's
// method of creating a User
func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(u.CreateUserRequest) // Assert that the underlying type of the interface we received is of CreateUserRequest, and if it is extract the value
		// TODO: Have some validation somewhere that can prevent us from wasting effort
		// if the request isn't going to be valid.
		id, err := s.CreateUserAccount(ctx, req.Username, req.Password)
		if err != nil {
			return nil, err
		} // Call the service method
		err = s.CreateUserProfile(ctx, id, req.FirstName, req.LastName, req.Email, req.Phone) // Call service method to create profile info for the user
		if err != nil {
			s.DeleteUserAccount(ctx, id)
		}
		// TODO: Error handling to delete the user account if the user profile is unable to be made
		return u.CreateUserResponse{ID: id}, err // Return a response in the shape we specified in this struct
	}
}

func makeGetUserAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(u.GetUserRequest)
		account, err := s.GetUserAccount(ctx, req.ID)
		return u.GetUserAccountResponse{ID: account.ID, Username: account.Username, JoinedOn: account.JoinedOn}, err
	}
}

func makeLoginUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(u.LoginRequest)

		detailedUser, err := s.Login(ctx, req.Username, req.Password)

		account, profile := &detailedUser.Account, &detailedUser.Profile

		return u.LoginResponse{
			Account: u.LoginAccount{
				ID:       account.ID,
				Username: account.Username,
				JoinedOn: account.JoinedOn,
			},
			Profile: u.LoginProfile{
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
				Email:     profile.Email,
				Phone:     profile.Phone,
				LastLogin: profile.LastLogin,
			},
		}, err

	}
}

func makeUpdateUserProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(u.UpdateProfileRequest)

		// Unmarshall Updates struct into map based on json tags
		updatesMap := structToMap(req.Updates)

		err := s.UpdateUserProfile(ctx, req.AccountID, updatesMap)

		if err != nil {
			return nil, err
		}

		return u.UpdateProfileResponse{OK: "ok"}, nil
	}
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")

		// remove omitEmpty
		omitEmpty := false
		if strings.HasSuffix(tag, "omitempty") {
			omitEmpty = true
			idx := strings.Index(tag, ",")
			if idx > 0 {
				tag = tag[:idx]
			} else {
				tag = ""
			}
		}

		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				if !(omitEmpty && reflectValue.Field(i).IsZero()) {
					res[tag] = field
				}
			}
		}
	}
	return res
}
