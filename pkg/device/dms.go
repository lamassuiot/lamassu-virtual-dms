package device

type CSR struct {
	Id                     int    `json:"id"`
	Name                   string `json:"dms_name"`
	CountryName            string `json:"country"`
	StateOrProvinceName    string `json:"state"`
	LocalityName           string `json:"locality"`
	OrganizationName       string `json:"organization"`
	OrganizationalUnitName string `json:"organization_unit,omitempty"`
	CommonName             string `json:"common_name"`
	EmailAddress           string `json:"mail,omitempty"`
	Status                 string `json:"status"`
	CsrFilePath            string `json:"csrpath,omitempty"`
	Url                    string `json:"url"`
}

type CSRs struct {
	CSRs []CSR `json:"-"`
}
type ProvisionForm struct {
	CountryName            string `json:"c"`
	StateOrProvinceName    string `json:"st"`
	OrganizationName       string `json:"o"`
	OrganizationalUnitName string `json:"ou"`
	CommonName             string `json:"cn"`
	KeyType                string `json:"key_type"`
	KeyBits                int    `json:"key_bits"`
}

const (
	PendingStatus  = "NEW"
	ApprovedStatus = "APPROVED"
	DeniedStatus   = "DENIED"
	RevokedStatus  = "REVOKED"
)
