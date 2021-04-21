package account

/* Here we are defining the structure of our User object, defining
the fields it should have and how they are expected to be received via
JSON format. */
type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
