package parking_charge_notice

import (
	"bytes"
	"encoding/json"
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

type Information struct {
	Sref                  string `json:"sref" validate:"required"`
	NoticeNumber          string `json:"notice_number" validate:"required"`
	VehicleRegistration   string `json:"vehicle_registration" validate:"required,min=2"`
	Contravention         string `json:"contravention" validate:"required"`
	ContraventionDatetime string `json:"contravention_datetime" validate:"required"`
	EntryExitDatetime     struct {
		Entry string `json:"entry,omitempty"`
		Exit  string `json:"exit,omitempty"`
	} `json:"entry_exit_datetime,omitempty"`
	ObservationDatetime struct {
		From string `json:"from,omitempty"`
		To   string `json:"to,omitempty"`
	} `json:"observation_datetime,omitempty"`
	Location       string `json:"location" validate:"required"`
	NoticeToKeeper struct {
		File string `json:"file,omitempty" validate:"omitempty,base64"`
		URL  string `json:"url,omitempty"  validate:"omitempty,http_url"`
	} `json:"notice_to_keeper"`
	Pofa             bool     `json:"pofa"`
	TotalDue         int      `json:"total_due" validate:"required"`
	ReducedAmount    float64  `json:"reduced_amount" validate:"required"`
	ReducePeriodEnds string   `json:"reduce_period_ends" validate:"required"`
	Photos           []string `json:"photos,omitempty"`
	PaymentURL       string   `json:"payment_url"`
	AppealURL        string   `json:"appeal_url"`
}

func (notice *Information) Validate() error {

	dateNow := time.Now()

	if !govalidator.IsRFC3339(notice.ContraventionDatetime) {
		return fmt.Errorf("invalid ContraventionDatetime format, should be RFC3339 - please see documentation")
	}

	cDateTime, _ := time.Parse(time.RFC3339, notice.ContraventionDatetime)

	if cDateTime.After(dateNow) {
		return fmt.Errorf("invalid ContraventionDatetime is a future date - please see documentation")
	}

	if len(notice.EntryExitDatetime.Exit) > 0 {
		if !govalidator.IsRFC3339(notice.EntryExitDatetime.Exit) {
			return fmt.Errorf("invalid ExitDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.EntryExitDatetime.Entry) > 0 {
		if !govalidator.IsRFC3339(notice.EntryExitDatetime.Entry) {
			return fmt.Errorf("invalid EntryDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.ObservationDatetime.To) > 0 {
		if !govalidator.IsRFC3339(notice.ObservationDatetime.To) {
			return fmt.Errorf("invalid ObservationToDatetime format, should be RFC3339 - please see documentation")
		}
	}

	if len(notice.ObservationDatetime.From) > 0 {
		if !govalidator.IsRFC3339(notice.ObservationDatetime.From) {
			return fmt.Errorf("invalid ObservationFromDatetime format, should be RFC3339 - please see documentation")
		}
	}

	validate := validator.New()
	return validate.Struct(notice)

}

func (notice *Information) Send(apiKey string) error {

	if len(os.Getenv("DEVELOPMENT")) == 0 {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(log.DebugLevel)

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

	req, err := http.NewRequest("POST", sendURL, bytes.NewBuffer(noticeData))
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
			log.Warnf("Non-200: %d %s", resp.StatusCode, resp.Status)
			log.Warnf("%s | %s", notice.Sref, string(body))
			return fmt.Errorf("Invalid code returned from api server (%d)", resp.StatusCode)

		} else {
			log.Debugln("OK")
			log.Debugln(string(body))
		}
	}
	return nil

}
