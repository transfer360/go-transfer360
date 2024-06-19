package search

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"time"
)

type ClientInfo struct {
	ClientID           string `firestore:"clientid"`
	IssuerID           string `json:"issuer_id" firestore:"issuer_id,omitempty"`
	SoftwareProviderID int    `firestore:"software_id,omitempty"`
}

type SearchResult struct {
	Sref                     string              `json:"sref"  firestore:"sref"`
	IsHirerVehicle           bool                `json:"is_hirer_vehicle" firestore:"is_hirer_vehicle" `
	VRM                      string              `json:"vrm"  firestore:"vrm"`
	ContraventionDate        string              `json:"contravention_date"  firestore:"contravention_date"`
	Reference                string              `json:"your_reference"  firestore:"your_reference"`
	LeaseCompany             LeaseCompanyAddress `json:"lease_company,omitempty" firestore:"lease_company,omitempty"`
	SearchPartnerID          int                 `json:"-" firestore:"search_partner_id,omitempty"`
	SearchPartnerDescription string              `json:"-" firestore:"search_partner_description,omitempty"`
}

type NoticeInformation struct {
	NoticeType          int    `json:"-" firestore:"notice_type,omitempty"`
	NoticeNumber        string `json:"notice_number,omitempty" firestore:"notice_number,omitempty"`
	VehicleRegistration string `json:"vehicle_registration,omitempty" firestore:"vehicle_registration,omitempty"`
	NoticeReceived      string `json:"-" firestore:"notice_received,omitempty"`
	DocumentReference   string `json:"-" firestore:"document_id,omitempty"`
}

type CreateSearchRecord struct {
	Sref              string       `firestore:"sref"`
	Client            ClientInfo   `firestore:"client"`
	Status            int          `firestore:"status"`
	StatusDescription string       `firestore:"status_description"`
	SearchDate        time.Time    `firestore:"search_date"`
	Result            SearchResult `firestore:"result"`
	StatusChanged     time.Time    `firestore:"status_changed,omitempty"`
}

const SEARCHES_COLLECTION = "searches"

func (s CreateSearchRecord) Save(ctx context.Context, client *firestore.Client) (docref string, err error) {

	collection := SEARCHES_COLLECTION

	docref = ""
	_, err = client.Collection(collection).NewDoc().Set(ctx, s)
	if err != nil {
		log.Error(err)
		return docref, fmt.Errorf("unexpected error updating search")
	} else {

		itr := client.Collection(collection).Where("sref", "==", s.Sref).Documents(ctx)
		for {
			doc, err := itr.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				}
			} else {
				if doc.Exists() {
					docref = doc.Ref.ID
				}
			}

		}
	}

	return docref, nil
}

func (c CreateSearchRecord) Delete(ctx context.Context, client *firestore.Client) error {

	collection := SEARCHES_COLLECTION

	itr := client.Collection(collection).Where("sref", "==", c.Sref).Documents(ctx)

	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorln(err)
				return err
			}
		} else {
			if doc.Exists() {
				_, err = doc.Ref.Delete(ctx)
				if err != nil {
					log.Errorf("Error deleting search record %s", c.Sref)
					return err
				}
			}
		}
	}

	return nil

}
