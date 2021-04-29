package accountsrv

type (
	CreateUserRequest struct {
		OrgID     string `json:"org_id"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		OrgType   string `json:"org_type"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
	}
	CreateUserResponse struct {
		ID string `json:"id"`
	}

	GetUserRequest struct {
		ID string `json:"id"`
	}
	GetUserAccountResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		JoinedOn string `json:"joined_on"`
	}

	LoginRequest struct {
		OrgID    string
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		UserAccount LoginUserAccount `json:"user_account"`
		UserProfile LoginUserProfile `json:"user_profile"`
		OrgAccount  LoginOrgAccount  `json:"org_account"`
	}
	LoginUserAccount struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	LoginUserProfile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}
	LoginOrgAccount struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	UpdateProfileRequest struct {
		AccountID string
		Updates   ProfileUpdates `json:"profile_updates"`
	}
	UpdateProfileResponse struct {
		OK string `json:"ok"`
	}
	ProfileUpdates struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
	}

	CreateOrgRequest struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
		Timezone string `json:"timezone"`
		Website  string `json:"website"`
	}
	CreateOrgResponse struct {
		ID string `json:"id"`
	}
)
