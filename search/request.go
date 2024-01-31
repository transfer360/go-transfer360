package search

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator/v10"
	"time"
)

type Request struct {
	// vehicle registration - required field
	VRM string `json:"vrm" validate:"required"`
	// date/time of vehicle contravention - required field and must be in R3339 format
	DateTime string `json:"contravention_date" validate:"required"`
	// a reference number which goes with the vehicle registration - required field
	Reference  string `json:"your_reference" validate:"required"`
	InitalSref string `json:"inital_sref,omitempty"` // if a search reference is generated and passed through
}

func (sr *Request) Validate() error {

	dateNow := time.Now()

	if !govalidator.IsRFC3339(sr.DateTime) {
		return fmt.Errorf("invalid datetime format, should be RFC3339 - please see documentation")
	}

	cDateTime, _ := time.Parse(time.RFC3339, sr.DateTime)

	if cDateTime.After(dateNow) {
		return fmt.Errorf("invalid ContraventionDatetime is a future date - please see documentation")
	}

	validate := validator.New()
	return validate.Struct(sr)

}
