package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
)

// A custom error we can pass back in place of the SQL error in the event
// that our DB functions return any failures.
var RepoErr = errors.New("repository error")

// Exposes similar methods to the Service interface, but these will specifically
// help us deal with the DB whereas the Service interface is the methods we expose
// for the service as a whole.
type Repository interface {
	CreateUserAccount(ctx context.Context, account Account) error
	CreateUserProfile(ctx context.Context, profile Profile) error
	GetUser(ctx context.Context, id string) (User, error)
}

// Defining a struct we will create methods for to implement the Repository interface
type repo struct {
	db     *sql.DB    // Pointer to DB
	logger log.Logger // The Logger to use
}

// Factory func to initiate a new Repository interace with the underlying DB and Logger
// of choice
func NewRepo(db *sql.DB, logger log.Logger) Repository {
	// Apparently if you implement the methods defined on an interface on the
	// underlying value (this case a struct), you can return an instance of that
	// interface like this (here returning a pointer to the struct) where the return
	// value specified on the func is the Interface but the function returns a pointer
	// to the underlying implementation of the struct... cool!
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

// Defining the method that will handle creating the user in the DB for the
// Repository interface to use.
// This is using a pointer receiver type so this method can mutate the parent
// repo struct directly if it wanted to...
func (repo *repo) CreateUserAccount(ctx context.Context, account Account) error {
	sqlCmd := `
		INSERT INTO user_accounts (id, username, password)
		VALUES ($1, $2, $3)`

	// TODO: Add new tables for things like user details for the user's name, etc

	// Validation framework?
	if account.Username == "" || account.Password == "" {
		return errors.New("username and password are required")
	}

	_, err := repo.db.ExecContext(ctx, sqlCmd, account.ID, account.Username, account.Password)
	if err != nil {
		return errors.New("error saving user account")
	}
	return nil
}

func (repo *repo) CreateUserProfile(ctx context.Context, profile Profile) error {
	sqlCmd := `
		INSERT INTO user_profiles (account_id, first_name, last_name, email, phone)
		VALUES ($1, $2, $3, $4, $5)`

	// TODO: Add new tables for things like user details for the user's name, etc

	// Validation framework?
	if profile.FirstName == "" || profile.LastName == "" {
		return errors.New("first name and last name are required")
	}

	_, err := repo.db.ExecContext(ctx, sqlCmd, profile.AccountID, profile.FirstName, profile.LastName, profile.Email, profile.Phone)
	if err != nil {
		return errors.New("error saving user account")
	}
	return nil
}

// Defining the method that will handle finding a user in the DB for the
// Repository interface to use
func (repo *repo) GetUser(ctx context.Context, id string) (User, error) {
	var user User
	err := repo.db.QueryRow(`
	SELECT acct.id, acct.username, acct.joined_on, prof.first_name,prof.last_name, prof.email, prof.phone, prof.last_login
	FROM user_accounts AS acct
	INNER JOIN user_profiles as prof
	ON acct.id = prof.account_id
	WHERE id=$1`,
		id).Scan(&user.Account.ID, &user.Account.Username, &user.Account.JoinedOn, &user.Profile.FirstName, &user.Profile.LastName, &user.Profile.Email, &user.Profile.Phone, &user.Profile.LastLogin)
	if err != nil {
		fmt.Println(err)
		return user, errors.New("no user found")
	}
	return user, nil
}

// TODO: I should have a wrapper function over fields/values that could potentially be NULL
// from the SQL query that can then convert to the receiver's datatype's zero-value OR
// omit that from the response entirely.

// This should be handled here so as to not taint any of the other code with DB
// related workarounds or anything like that.
