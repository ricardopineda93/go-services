package organization

type Organization struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}
