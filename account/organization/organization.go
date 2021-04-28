package account

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"org_name"`
	Type     string `json:"org_type"`
	JoinedOn string `json:"joined_on"`
}

type Profile struct {
	AccountID string `json:"account_id"`
	Website   string `json:"website"`
	Phone     string `json:"phone"`
	Timezone  string `json:"timezone"`
	Address   struct {
		Street string `json:"street"`
		City   string `json:"city"`
		State  string `json:"state"`
		Zip    string `json:"zip"`
	} `json:"address"`
	PrimaryContact struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
	} `json:"primary_contact"`
}

type ProviderDetails struct {
	AccountID string `json:"account_id"`
	NPI       string `json:"npi"`
	TaxID     string `json:"tax_id"`
}

type PayorDetails struct {
	AccountID string `json:"account_id"`
	PayorID   string `json:"payor_id"`
}
