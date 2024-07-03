package search

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type Result struct {
	Sref              string              `json:"sref"`
	IsHirerVehicle    bool                `json:"is_hirer_vehicle"`
	VRM               string              `json:"vrm"`
	ContraventionDate string              `json:"contravention_date"`
	Reference         string              `json:"your_reference"`
	LeaseCompany      LeaseCompanyAddress `json:"lease_company,omitempty"`
}

func (r *Result) FromSREF(ctx context.Context, sref string, client *firestore.Client) error {

	data := struct {
		Result struct {
			IsHirerVehicle    bool                `firestore:"is_hirer_vehicle"`
			VRM               string              `firestore:"vrm"`
			ContraventionDate string              `firestore:"contravention_date"`
			Reference         string              `firestore:"your_reference"`
			LeaseCompany      LeaseCompanyAddress `firestore:"lease_company"`
		} `firestore:"result"`
	}{}

	itr := client.Collection("searches").Where("sref", "==", sref).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			log.Errorln(err)
			return err
		} else {
			if err := doc.DataTo(&data); err != nil {
				log.Errorln(err)
				return err
			}
			break
		}

	}

	r.VRM = data.Result.VRM
	r.IsHirerVehicle = data.Result.IsHirerVehicle
	r.ContraventionDate = data.Result.ContraventionDate
	r.Reference = data.Result.Reference
	r.LeaseCompany = data.Result.LeaseCompany

	return nil

}
