package parking_charge_notice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator/v10"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	pcn "github.com/transfer360/sys360/notices/parking_charge_notice"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Information struct {
	pcn.Data
}

var ErrNoticeAlreadyExists = errors.New("notice already exists")
var ErrIssuerNotSetup = errors.New("issuer is not setup")

// Validate ----------------------------------------------------------------------------------------------------------
func (notice *Information) Validate() error {

	dateNow := time.Now()

	if !govalidator.IsRFC3339(notice.ContraventionDateTime) {

		notice.ContraventionDateTime = strings.ReplaceAll(notice.ContraventionDateTime, "T", " ")
		notice.ContraventionDateTime = strings.ReplaceAll(notice.ContraventionDateTime, "Z", "")

		if tm, err := time.Parse("2006-01-02 15:04:05", notice.ContraventionDateTime); err == nil {
			notice.ContraventionDateTime = tm.Format(time.RFC3339)
		} else {
			return fmt.Errorf("invalid ContraventionDatetime format, should be RFC3339 - please see documentation")
		}

	}

	cDateTime, _ := time.Parse(time.RFC3339, notice.ContraventionDateTime)

	if cDateTime.After(dateNow) {
		return fmt.Errorf("invalid ContraventionDatetime is a future date - please see documentation")
	}

	if len(notice.EntryExit.Exit) > 0 {
		if !govalidator.IsRFC3339(notice.EntryExit.Exit) {
			return fmt.Errorf("invalid ExitDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.EntryExit.Entry) > 0 {
		if !govalidator.IsRFC3339(notice.EntryExit.Entry) {
			return fmt.Errorf("invalid EntryDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.Observation.To) > 0 {
		if !govalidator.IsRFC3339(notice.Observation.To) {
			return fmt.Errorf("invalid ObservationToDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.Observation.From) > 0 {
		if !govalidator.IsRFC3339(notice.Observation.From) {
			return fmt.Errorf("invalid ObservationFromDatetime format, should be RFC3339 - please see documentation")
		}
	}

	validate := validator.New()
	return validate.Struct(notice)

}

// Send ----------------------------------------------------------------------------------------------------------
func (notice *Information) Send(apiKey string) error {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	err := notice.Validate()

	if err != nil {
		return fmt.Errorf("go-transfer360 information invalid: %w", err)
	}

	sendURL := "https://api.transfer360.io/notice/parking_charge"
	noticeData, err := json.Marshal(notice)
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Debugln(string(noticeData))

	req, err := http.NewRequest("POST", sendURL, bytes.NewBuffer(noticeData))
	if err != nil {
		log.Errorln(err)
		return err
	}
	req.Header.Set("api_key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		client.Timeout = time.Second * 20
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorln(err)
		return fmt.Errorf("sending go-transfer360 to api server %w", err)
	} else {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		log.Debugf("response Body: %s", string(body))

		if resp.StatusCode != http.StatusOK {

			if resp.StatusCode == http.StatusConflict {
				return ErrNoticeAlreadyExists
			}
			if resp.StatusCode == http.StatusTooEarly {
				return ErrIssuerNotSetup
			}

			log.Warnf("Non-200: %d %s", resp.StatusCode, resp.Status)
			log.Warnf("%s | %s", notice.SearchReference, string(body))
			return fmt.Errorf("error code returned from api server (%d)[%s]", resp.StatusCode, string(body))

		} else {
			log.Debugln("OK")
			log.Debugln(string(body))
		}
	}
	return nil

}

// ----------------------------------------------------------------------------------------------------------
