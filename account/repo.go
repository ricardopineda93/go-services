package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	u "github.com/rjjp5294/gokit-tutorial/account/user"
)

// A custom error we can pass back in place of the SQL error in the event
// that our DB functions return any failures.
var RepoErr = errors.New("repository error")

// Exposes similar methods to the Service interface, but these will specifically
// help us deal with the DB whereas the Service interface is the methods we expose
// for the service as a whole.
type Repository interface {
	CreateUserAccount(ctx context.Context, account u.Account) error
	DeleteUserAccount(ctx context.Context, id string) error
	CreateUserProfile(ctx context.Context, profile u.Profile) error
	GetUserProfile(ctx context.Context, accountID string) (u.Profile, error)
	UpdateUserProfile(ctx context.Context, accountID string, updates map[string]interface{}) error
	GetUserAccount(ctx context.Context, id string) (u.Account, error)
	GetAccountByLoginCredentials(ctx context.Context, username string, password string) (u.Account, error)
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
func (repo *repo) CreateUserAccount(ctx context.Context, account u.Account) error {
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

func (repo *repo) DeleteUserAccount(ctx context.Context, id string) error {
	sqlCmd := `DELETE FROM user_accounts WHERE id=$1`

	_, err := repo.db.ExecContext(ctx, sqlCmd, id)
	if err != nil {
		return errors.New("error deleting user account")
	}
	return nil
}

func (repo *repo) CreateUserProfile(ctx context.Context, profile u.Profile) error {
	sqlCmd := `
		INSERT INTO user_profiles (account_id, first_name, last_name, email, phone)
		VALUES ($1, $2, $3, $4, $5)`

	// Validation framework?
	if profile.FirstName == "" || profile.LastName == "" {
		return errors.New("first name and last name are required")
	}

	_, err := repo.db.ExecContext(ctx, sqlCmd, profile.AccountID, profile.FirstName, profile.LastName, profile.Email, profile.Phone)
	if err != nil {
		return errors.New("error saving user profile")
	}
	return nil
}

func (repo *repo) GetUserProfile(ctx context.Context, accountID string) (u.Profile, error) {
	var profile u.Profile

	err := repo.db.QueryRowContext(ctx,
		`SELECT first_name, last_name, email, phone, last_login
	FROM user_profiles
	WHERE account_id=$1`,
		accountID).Scan(&profile.FirstName, &profile.LastName, &profile.Email, &profile.Phone, &profile.LastLogin)

	if err != nil {
		return profile, errors.New("error getting user profile")
	}

	return profile, nil
}

func (repo *repo) UpdateUserProfile(ctx context.Context, accountID string, updates map[string]interface{}) error {

	var updateStr string = ``

	fmt.Println(accountID)

	keysLeft := len(updates)
	for k, v := range updates {
		updateStr += genUpdateLine(k, v)
		if keysLeft -= 1; keysLeft > 0 {
			updateStr += ","
		}
	}
	sqlCmd := `UPDATE user_profiles SET ` + updateStr + `WHERE account_id = $1`

	fmt.Println(sqlCmd)

	_, err := repo.db.ExecContext(ctx,
		sqlCmd,
		accountID)

	if err != nil {
		fmt.Println(err)
		return errors.New("unable to update user profile")
	}
	return nil
}

// Defining the method that will handle finding a user in the DB for the
// Repository interface to use
func (repo *repo) GetUserAccount(ctx context.Context, id string) (u.Account, error) {
	var account u.Account
	err := repo.db.QueryRow(`
	SELECT acct.id, acct.username, acct.joined_on
	FROM user_accounts AS acct
	WHERE id=$1`,
		id).Scan(&account.ID, &account.Username, &account.JoinedOn)
	if err != nil {
		fmt.Println(err)
		return account, errors.New("no user found")
	}
	return account, nil
}

func (repo *repo) GetAccountByLoginCredentials(ctx context.Context, username string, password string) (u.Account, error) {
	var account u.Account

	// First find user account by username
	err := repo.db.QueryRowContext(ctx,
		`SELECT id, username, password,joined_on
	FROM user_accounts
	WHERE username=$1`,
		username).Scan(&account.ID, &account.Username, &account.Password, &account.JoinedOn)

	// Check to see if the user account's password matches input password
	if err != nil {
		fmt.Println(err)
		return u.Account{}, errors.New("invalid credentials")
	}

	if password != account.Password {
		return u.Account{}, errors.New("incorrect password")
	}

	return account, nil
}

func genUpdateLine(k string, v interface{}) string {
	update := k + " = "

	switch t := v.(type) {
	case string:
		if v == "DEFAULT" {
			update += t
		} else {
			// update += fmt.Sprintf("%q", t)
			update += "'" + t + "'"
		}
	case nil:
		return ""
	default:
		update += fmt.Sprintf("%v", t)

	}

	update += " "
	return update
}

// TODO: I should have a wrapper function over fields/values that could potentially be NULL
// from the SQL query that can then convert to the receiver's datatype's zero-value OR
// omit that from the response entirely.

// This should be handled here so as to not taint any of the other code with DB
// related workarounds or anything like that.
