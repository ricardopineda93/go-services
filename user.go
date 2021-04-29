package accountsrv

type UserAccount struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	JoinedOn string `db:"joined_on"`
	OrgType  string `db:"org_type"`
}

type UserProfile struct {
	AccountID string `db:"account_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Phone     string `db:"phone"`
	LastLogin string `db:"last_login"`
}

type DetailedUser struct {
	Account UserAccount
	Profile UserProfile
}
