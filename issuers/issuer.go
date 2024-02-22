package issuers

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type IssuerInformation struct {
	T360ID         string `json:"t360_id" firestore:"t360_id"`
	ClientID       string `json:"clientid"`
	IssuerID       string `json:"issuer_id"`
	SoftwareID     int    `json:"software_id"`
	Issuer         string `json:"issuer" firestore:"issuer"`
	PrivateParking bool   `json:"private_parking"`
}

func GetIssuerInformationFromT360ID(ctx context.Context, issuerID string, fsclient *firestore.Client) (IssuerInformation, error) {

	oi := IssuerInformation{}

	issuer := struct {
		T360ID           string `firestore:"t360_id"`
		SoftwareProvider int    `firestore:"software_provider"`
		Issuer           string `firestore:"issuer"`
		PrivateParking   bool   `firestore:"private_parking"`
	}{}

	itr := fsclient.Collection("registered_issuers").Where("t360_id", "==", issuerID).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Error("GetOperatorIssuerID:", err)
				return oi, err
			}
		}

		if doc.Exists() {
			err = doc.DataTo(&issuer)
			if err != nil {
				log.Error("GetOperatorIssuerID:", err)
				return oi, err
			}
		}
	}

	oi.T360ID = issuer.T360ID
	oi.PrivateParking = issuer.PrivateParking
	oi.Issuer = issuer.Issuer
	oi.IssuerID = issuer.T360ID
	oi.ClientID = issuer.T360ID
	oi.SoftwareID = issuer.SoftwareProvider

	return oi, nil

}
