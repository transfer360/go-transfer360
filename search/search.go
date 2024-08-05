// Package search is a package that can search the Transfer360 database for lease vehicles. docs: https://transfer360.dev/
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
)

func SendEnquiry(ctx context.Context, n Request, apiKey string) (scanReturn Result, err error) {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	if len(apiKey) == 0 {
		return scanReturn, fmt.Errorf("missing API Key")
	}

	err = n.Validate()

	if err != nil {
		return scanReturn, err
	}

	url := "https://api.transfer360.io/search"

	jsonStr, err := json.Marshal(n)
	if err != nil {
		log.Error("SendEnquiry:1:", err)
		return scanReturn, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error("SendEnquiry:1.1:", err)
		return scanReturn, err
	}
	req.Header.Set("api_key", apiKey) // Not added as an Environment variables as some software providers have different Keys per client
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		client.Timeout = time.Second * 60
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") { // dont log this error out.
			return scanReturn, err
		} else {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				log.Debugln("SendEnquiry:2: [TimeOut]", err)
				return scanReturn, err
			} else {
				log.Error("SendEnquiry:2:", err)
				return scanReturn, err
			}
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		sr := Result{}

		err = json.NewDecoder(resp.Body).Decode(&sr)
		if err != nil {
			log.Error("SendEnquiry:3:", err)
			return scanReturn, fmt.Errorf("%w %s", ErrInvalidSearchResultBody, err.Error())
		}

		scanReturn.VRM = n.VRM
		scanReturn.Reference = n.Reference
		scanReturn.ContraventionDate = n.DateTime
		scanReturn.Sref = sr.Sref
		scanReturn.IsHirerVehicle = sr.IsHirerVehicle
		scanReturn.LeaseCompany = sr.LeaseCompany

		return scanReturn, nil

	} else if resp.StatusCode == http.StatusGatewayTimeout {

		return scanReturn, ErrTimeOutStatusCode

	} else if resp.StatusCode == http.StatusServiceUnavailable {

		return scanReturn, ErrUnableToHandleStatusCode

	} else {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("SendEnquiry:4:", err)
			return scanReturn, fmt.Errorf("%w %s", ErrInvalidSearchResultBody, err.Error())
		}

		return scanReturn, fmt.Errorf("%w invalid StatusCode returned (%d) [%s]", ErrInvalidSearchResultCodeReturned, resp.StatusCode, string(body))

	}

}
