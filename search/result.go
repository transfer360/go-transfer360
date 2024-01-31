package search

type Result struct {
	Sref              string              `json:"sref"`
	IsHirerVehicle    bool                `json:"is_hirer_vehicle"`
	VRM               string              `json:"vrm"`
	ContraventionDate string              `json:"contravention_date"`
	Reference         string              `json:"your_reference"`
	LeaseCompany      LeaseCompanyAddress `json:"lease_company,omitempty"`
}
