package accountsrv

type UserAccount struct {
	ID       string `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password,omitempty"`
	OrgType  string `db:"org_type" json:"org_type"`
	JoinedOn string `db:"joined_on" json:"joined_on"`
}

type UserProfile struct {
	AccountID string `db:"account_id" json:"account_id,omitempty"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email,omitempty"`
	Phone     string `db:"phone" json:"phone,omitempty"`
	LastLogin string `db:"last_login" json:"last_login,omitempty"`
}

type DetailedUser struct {
	Account UserAccount `json:"account"`
	Profile UserProfile `json:"profile"`
}

type LoginUser struct {
	User DetailedUser `json:"user"`
	Org  DetailedOrg  `json:"org"`
}
