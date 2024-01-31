package search

type LeaseCompanyAddress struct {
	Companyname  string `json:"companyname,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	AddressLine3 string `json:"address_line3,omitempty"`
	AddressLine4 string `json:"address_line4,omitempty"`
	Postcode     string `json:"postcode,omitempty"`
}
