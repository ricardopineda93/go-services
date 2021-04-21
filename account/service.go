package account

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

// Service interface exposes what methods are available for the service to the Endpoint Layer.
// These are very high level methods that will be implemented and will actually be more
// control/orchastration methods that abstract away the more granular operations from the requestor,
// basically a controller function that relies on other methods/functions to carry out the operation
// as a whole while it just provides context, control and order.
type Service interface {
	CreateUser(ctx context.Context, email string, password string) (string, error)
	GetUser(ctx context.Context, id string) (User, error)
}

// The properties the service will contain
type service struct {
	repository Repository // To interface with the DB, this is an interface that will handle all the DB interaction, the service just needs to know this is the "persistance store"
	logger     log.Logger // To log and see what's going on inside the service
}

// Implement the Service interface using the service struct and methods defined for it.
// What's genius is that the repository field is itself an interface, and the methods
// defined for the service struct actually utilize the methods of the Repository interface
// to implement the methods of the Service interface... amazing.
func NewService(rep Repository, logger log.Logger) Service {
	// Return pointer to a service struct, which will be the concrete type implementing
	// the Service interface.
	return &service{
		repository: rep,
		logger:     logger,
	}
}

/*
It looks like the Service interface methods defined here for the service struct are basically
orcastrator -- the methods it implements use methods implemented by interfaces by it's
properties (logger, respository) to do granualar business logic while controlling the flow
of the work. Cool idea.
*/

// Method for creating a user for the service struct for the Service interface to implement
// Even though the Service's underlying value/type is pointer to service struct,
// we can still use the value receiver type. However, this means this method is operating on
// a COPY of the service struct rather than the "actual" underlying service struct
func (s service) CreateUser(ctx context.Context, email string, password string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	uuid, _ := uuid.NewV4()
	id := uuid.String()
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	// Use the respository interface's implementation of create user to actually do
	// the portion of the business logic, this implementation just abstracts that away
	// and provides context and orchastration.. very neat.
	if err := s.repository.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create user", id)

	return "Success!", nil

}

// Method for service struct for the Service interface to implement
func (s service) GetUser(ctx context.Context, id string) (User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	// Same thing here, using the repository property's methods to actually do the
	// fetching while this method just is kind of a control flow method.
	user, err := s.repository.GetUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return user, err
	}

	logger.Log("Get user", id)

	return user, nil
}
