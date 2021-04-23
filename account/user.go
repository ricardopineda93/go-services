package account

import (
	"database/sql"
	"encoding/json"
)

/* Here we are defining the structure of our User object, defining
the fields it should have and how they are expected to be received via
JSON format. */
type Account struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	JoinedOn string `json:"joined_on"`
}

type Profile struct {
	AccountID string `json:"account_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	LastLogin string `json:"last_login,omitempty"`
}

type User struct {
	Account Account `json:"user_account"`
	Profile Profile `json:"user_profile"`
}

/*
Custom implementations for handling NULL sql values and Marshalling those types into
JSON
*/

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
