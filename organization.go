package accountsrv

type OrgAccount struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	JoinedOn string `db:"joined_on" json:"joined_on"`
}

type OrgProfile struct {
	AccountID string `db:"account_id" json:"account_id"`
	Phone     string `db:"phone" json:"phone"`
	Address   string `db:"address" json:"address"`
	Timezone  string `db:"timezone" json:"timezone"`
	Website   string `db:"website" json:"website"`
}

type ProviderDetails struct {
	AccountID string `db:"account_id" json:"account_id"`
	NPI       string `db:"npi" json:"npi"`
	TaxID     string `db:"tax_id" json:"tax_id"`
}

type PayorDetails struct {
	AccountID string `db:"account_id" json:"account_id"`
	PayorID   string `db:"payor_id" json:"payor_id"`
}

type DetailedOrg struct {
	Account         OrgAccount      `json:"account"`
	Profile         OrgProfile      `json:"profile"`
	ProviderDetails ProviderDetails `json:"provider_details,omitempty"`
	PayorDetails    PayorDetails    `json:"payor_details,omitempty"`
}
