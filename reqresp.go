package accountsrv

type CreateUserRequest struct {
	OrgID     string `json:"org_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	OrgType   string `json:"org_type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type CreateUserResponse struct {
	ID  string `json:"id,omitempty"`
	Err error  `json:"error,omitempty"`
}

func (r CreateUserResponse) error() error { return r.Err }

type GetUserRequest struct {
	ID string `json:"id"`
}

type GetUserAccountResponse struct {
	UserAccount UserAccount `json:"user_account,omitempty"`
	Err         error       `json:"error,omitempty"`
}

func (r GetUserAccountResponse) error() error { return r.Err }

type LoginRequest struct {
	OrgID    string
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	LoginDetails LoginUser `json:"login_details"`
	Err          error     `json:"error,omitempty"`
}

func (r LoginResponse) error() error { return r.Err }

type UpdateProfileRequest struct {
	AccountID string
	Updates   ProfileUpdates `json:"profile_updates"`
}

type UpdateProfileResponse struct {
	OK  string `json:"ok"`
	Err error  `json:"error,omitempty"`
}

func (r UpdateProfileResponse) error() error { return r.Err }

type CreateOrgRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Timezone string `json:"timezone"`
	Website  string `json:"website"`
}

type CreateOrgResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r CreateOrgResponse) error() error { return r.Err }

type ProfileUpdates struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
