package account

import (
	"context"
	"reflect"
	"strings"

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
		if err != nil {
			return nil, err
		}
		return CreateUserResponse{ID: id}, err // Return a response in the shape we specified in this struct
	}
}

func makeGetUserAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		account, err := s.GetUserAccount(ctx, req.ID)
		return GetUserAccountResponse{ID: account.ID, Username: account.Username, JoinedOn: account.JoinedOn}, err
	}
}

func makeLoginUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		detailedUser, err := s.Login(ctx, req.OrgID, req.Username, req.Password)

		account, profile := &detailedUser.Account, &detailedUser.Profile

		return LoginResponse{
			Account: LoginAccount{
				ID:       account.ID,
				Username: account.Username,
				JoinedOn: account.JoinedOn,
			},
			Profile: LoginProfile{
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
		req := request.(UpdateProfileRequest)

		// Unmarshall Updates struct into map based on json tags
		updatesMap := structToMapByTag(req.Updates, "json")

		err := s.UpdateUserProfile(ctx, req.AccountID, updatesMap)

		if err != nil {
			return nil, err
		}

		return UpdateProfileResponse{OK: "ok"}, nil
	}
}

func makeCreateOrgEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrgRequest)

		id, err := s.CreateOrg(ctx, req.Name, req.Type, req.Phone, req.Address, req.Timezone, req.Website)

		if err != nil {
			return nil, err
		}

		return CreateOrgResponse{ID: id}, nil
	}
}

// Function takes in a Struct and returns a Map of that struct using the tags associated
// to the struct fields to create the keys of the Map
func structToMapByTag(item interface{}, tagName string) map[string]interface{} {

	// Declare the map that will be returned
	res := map[string]interface{}{}
	// If input item is null return the map
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
		tag := v.Field(i).Tag.Get(tagName)

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
				res[tag] = structToMapByTag(field, tagName)
			} else {
				if !(omitEmpty && reflectValue.Field(i).IsZero()) {
					res[tag] = field
				}
			}
		}
	}
	return res
}
