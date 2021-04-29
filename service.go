package accountsrv

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
	CreateUser(ctx context.Context, orgID string, username string, password string, orgType string, firstName string, lastName string, email string, phone string) (string, error)
	DeleteUserAccount(ctx context.Context, id string) error
	UpdateUserProfile(ctx context.Context, accountID string, updates map[string]interface{}) error
	GetUserAccount(ctx context.Context, id string) (UserAccount, error)
	Login(ctx context.Context, orgID string, username string, password string) (DetailedUser, error)

	CreateOrg(ctx context.Context, name string, orgType string, phone string, address string, timezone string, website string) (string, error)
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
func (s service) CreateUser(ctx context.Context, orgID string, username string, password string, orgType string, firstName string, lastName string, email string, phone string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	uuid, _ := uuid.NewV4()
	id := uuid.String()
	user := UserAccount{
		ID:       id,
		Username: username,
		Password: password,
		OrgType:  orgType,
	}

	// Use the respository interface's implementation of create user to actually do
	// the portion of the business logic, this implementation just abstracts that away
	// and provides context and orchastration.. very neat.
	if err := s.repository.CreateUserAccount(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	profile := UserProfile{
		AccountID: id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	if err := s.repository.CreateUserProfile(ctx, profile); err != nil {
		s.repository.DeleteUserAccount(ctx, id)
		level.Error(logger).Log("err", err)
		return "", err
	}

	if err := s.repository.AssociateUserToOrg(ctx, id, orgID); err != nil {
		s.repository.DeleteUserAccount(ctx, id)
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("created user", id)

	return id, nil
}

func (s service) DeleteUserAccount(ctx context.Context, id string) error {
	logger := log.With(s.logger, "method", "DeleteUserAccount")

	// Use the respository interface's implementation of create user to actually do
	// the portion of the business logic, this implementation just abstracts that away
	// and provides context and orchastration.. very neat.
	if err := s.repository.DeleteUserAccount(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("deleted user account", id)

	return nil
}

// Method for service struct for the Service interface to implement
func (s service) GetUserAccount(ctx context.Context, id string) (UserAccount, error) {
	logger := log.With(s.logger, "method", "GetUser")

	// Same thing here, using the repository property's methods to actually do the
	// fetching while this method just is kind of a control flow method.
	account, err := s.repository.GetUserAccount(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return account, err
	}

	logger.Log("Get user account", id)

	return account, nil
}

func (s service) Login(ctx context.Context, orgID string, username string, password string) (DetailedUser, error) {
	logger := log.With(s.logger, "method", "Login")

	account, err := s.repository.GetAccountByLoginCredentials(ctx, username, password)

	if err != nil {
		level.Error(logger).Log("err", err)
		return DetailedUser{}, err
	}

	if err := s.repository.ConfirmUserToOrgAssociation(ctx, account.ID, orgID); err != nil {
		level.Error(logger).Log("err", err)
		return DetailedUser{}, err
	}

	if err := s.repository.UpdateUserProfile(ctx, account.ID, map[string]interface{}{
		"last_login": "DEFAULT",
	}); err != nil {
		level.Error(logger).Log("err", err)
		return DetailedUser{}, err
	}

	profile, err := s.repository.GetUserProfile(ctx, account.ID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return DetailedUser{}, err
	}

	logger.Log("Login user", account.ID)

	return DetailedUser{
		Account: account,
		Profile: profile,
	}, nil
}

func (s service) UpdateUserProfile(ctx context.Context, accountID string, updates map[string]interface{}) error {
	logger := log.With(s.logger, "method", "UpdateUserProfile")

	err := s.repository.UpdateUserProfile(ctx, accountID, updates)

	if err != nil {
		return err
	}

	logger.Log("updated user profile", accountID)

	return nil
}

func (s service) CreateOrg(ctx context.Context, name string, orgType string, phone string, address string, timezone string, website string) (string, error) {
	logger := log.With(s.logger, "method", "CreateOrg")

	uuid, _ := uuid.NewV4()
	id := uuid.String()

	orgAccount := OrgAccount{
		ID:   id,
		Name: name,
		Type: orgType,
	}

	err := s.repository.CreateOrgAccount(ctx, orgAccount)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	orgProfile := OrgProfile{
		AccountID: id,
		Phone:     phone,
		Address:   address,
		Timezone:  timezone,
		Website:   website,
	}

	err = s.repository.CreateOrgProfile(ctx, orgProfile)

	if err != nil {
		s.repository.DeleteOrgAccount(ctx, id)
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("created organization", id)

	return id, nil
}
