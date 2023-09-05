// Package search is a package that can search the Transfer360 database for lease vehicles. docs: https://transfer360.dev/
package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator/v10"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

// ErrInvalidSearchResultCodeReturned - error raised when a search request to the Transfer360 API server returns a non 200 result code
var ErrInvalidSearchResultCodeReturned = errors.New("unexpected search result code returned")

// ErrInvalidSearchResultBody - error raised when a JSON body is expected on sending or return
var ErrInvalidSearchResultBody = errors.New("missing or invalid search result body returned")

// Request is search request struct, which will be converted to JSON and sent the Transfer360 API server
type Request struct {
	// vehicle registration - required field
	VRM string `json:"vrm" validate:"required"`
	// date/time of vehicle contravention - required field and must be in R3339 format
	DateTime string `json:"contravention_date" validate:"required"`
	// a reference number which goes with the vehicle registration - required field
	Reference string `json:"your_reference" validate:"required"`
	InitalSref        string `json:"inital_sref,omitempty"` // if a search reference is generated and passed through
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

type Result struct {
	Sref              string              `json:"sref"`
	IsHirerVehicle    bool                `json:"is_hirer_vehicle"`
	VRM               string              `json:"vrm"`
	ContraventionDate string              `json:"contravention_date"`
	Reference         string              `json:"your_reference"`
	LeaseCompany      LeaseCompanyAddress `json:"lease_company,omitempty"`
}

type LeaseCompanyAddress struct {
	Companyname  string `json:"companyname,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	AddressLine3 string `json:"address_line3,omitempty"`
	AddressLine4 string `json:"address_line4,omitempty"`
	Postcode     string `json:"postcode,omitempty"`
}

func SendEnquiry(n Request, apiKey string) (scanReturn Result, err error) {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	if len(apiKey)==0 {
		return scanReturn,fmt.Errorf("missing API Key")
	}

	log.SetLevel(log.DebugLevel)

	err = n.Validate()

	if err != nil {
		return scanReturn, err
	}

	url := "https://api.transfer360.io/search"

	jsonStr, err := json.Marshal(n)
	if err != nil {
		log.Error(err)
		return scanReturn, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("api_key", apiKey) // Not added as an Environment variables as some software providers have different Keys per client
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		client.Timeout = time.Second * 20
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return scanReturn, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		sr := Result{}

		err = json.NewDecoder(resp.Body).Decode(&sr)
		if err != nil {
			log.Error(err)
			return scanReturn, fmt.Errorf("%w %s", ErrInvalidSearchResultBody, err.Error())
		}

		scanReturn.Sref = sr.Sref
		scanReturn.IsHirerVehicle = sr.IsHirerVehicle
		scanReturn.LeaseCompany = sr.LeaseCompany

		return scanReturn, nil

	} else {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return scanReturn, fmt.Errorf("%w %s", ErrInvalidSearchResultBody, err.Error())
		}

		return scanReturn, fmt.Errorf("%w invalid StatusCode returned (%d) [%s]", ErrInvalidSearchResultCodeReturned, resp.StatusCode, string(body))

	}

}
