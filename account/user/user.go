package account

type Account struct {
	ID       string
	Username string
	Password string
	JoinedOn string
}

type Profile struct {
	AccountID string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	LastLogin string
}

type DetailedUser struct {
	Account Account
	Profile Profile
}
