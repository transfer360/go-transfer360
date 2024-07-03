package parking_charge_notice

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type HirerInformation struct {
	CompanyName  string `json:"company_name,omitempty" bigquery:"CompanyName"`
	Name         string `json:"name,omitempty" bigquery:"Name"`
	Surname      string `json:"surname,omitempty" bigquery:"Surname"`
	AddressLine1 string `json:"address_line_1,omitempty" bigquery:"AddressLine1"`
	AddressLine2 string `json:"address_line_2,omitempty" bigquery:"AddressLine2"`
	AddressLine3 string `json:"address_line_3,omitempty" bigquery:"AddressLine3"`
	AddressLine4 string `json:"address_line_4,omitempty" bigquery:"AddressLine4"`
	PostCode     string `json:"post_code,omitempty" bigquery:"PostCode"`
	Country      string `json:"country,omitempty" bigquery:"Country"`
}

var ERRHirerNotFound = errors.New("hirer not found")

func GetHirer(ctx context.Context, sref string, client *firestore.Client) (HirerInformation, error) {

	hirerData := struct {
		LeaseReturn struct {
			ContactInfo HirerInformation `bigquery:"ContactInfo"`
		} `bigquery:"LeaseReturn"`
	}{}

	itr := client.Collection("parking_charge_notices_status_update").Where("Sref", "==", sref).Documents(ctx)

	hirerFound := false
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorln(err)
				return HirerInformation{}, err
			}
		} else {
			err = doc.DataTo(&hirerData)
			if err != nil {
				log.Errorln(err)
				return HirerInformation{}, err
			} else {
				hirerFound = true
				break
			}
		}

	}

	if !hirerFound {
		return HirerInformation{}, ERRHirerNotFound
	}

	return hirerData.LeaseReturn.ContactInfo, nil
}
