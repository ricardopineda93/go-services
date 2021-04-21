package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
)

// A custom error we can pass back in place of the SQL error in the event
// that our DB functions return any failures.
var RepoErr = errors.New("Unable to handle Repo Request")

// Exposes similar methods to the Service interface, but these will specifically
// help us deal with the DB whereas the Service interface is the methods we expose
// for the service as a whole.
type Repository interface {
	CreateUser(ctx context.Context, user User) error
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
func (repo *repo) CreateUser(ctx context.Context, user User) error {
	sqlCmd := `
		INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)`

	if user.Email == "" || user.Password == "" {
		return RepoErr
	}

	_, err := repo.db.ExecContext(ctx, sqlCmd, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Defining the method that will handle finding a user in the DB for the
// Repository interface to use
func (repo *repo) GetUser(ctx context.Context, id string) (User, error) {
	var user User
	err := repo.db.QueryRow("SELECT id, email FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, RepoErr
	}
	return user, nil
}
