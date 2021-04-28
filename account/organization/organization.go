package account

type Account struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Type     string `db:"type"`
	JoinedOn string `db:"joined_on"`
}

type Profile struct {
	AccountID string `db:"account_id"`
	Phone     string `db:"phone"`
	Address   string `db:"address"`
	Timezone  string `db:"timezone"`
	Website   string `db:"website"`
}

type ProviderDetails struct {
	AccountID string `db:"account_id"`
	NPI       string `db:"npi"`
	TaxID     string `db:"tax_id"`
}

type PayorDetails struct {
	AccountID string `db:"account_id"`
	PayorID   string `db:"payor_id"`
}

type DetailedOrg struct {
	Account         Account
	Profile         Profile
	ProviderDetails ProviderDetails
	PayorDetails    PayorDetails
}
