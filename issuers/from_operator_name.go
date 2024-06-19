package issuers

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

var ErrIssuerNotFound = errors.New("issuer information not found")

func FromOperatorsName(ctx context.Context, operatorName string, fs *firestore.Client) (IssuerInformation, error) {

	iInfo := IssuerInformation{}
	isserInfo := struct {
		Issuer           string `firestore:"issuer"`
		T360ID           string `firestore:"t360_id"`
		SoftwareProvider int    `firestore:"software_provider"`
	}{}

	itr := fs.Collection("registered_issuers").Where("operator_name", "==", operatorName).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorf("FromOperatorsName:[%s]:%v", operatorName, err)
				return iInfo, err
			}
		} else {
			if doc.Exists() {
				err = doc.DataTo(&isserInfo)
				if err != nil {
					log.Errorf("FromOperatorsName:[%s]:%v", operatorName, err)
					return iInfo, err
				} else {
					break
				}
			}
		}
	}

	if len(isserInfo.Issuer) == 0 {
		return FromIssuerName(ctx, operatorName, fs)
	}

	iInfo.T360ID = isserInfo.T360ID
	iInfo.Issuer = isserInfo.Issuer
	iInfo.SoftwareID = isserInfo.SoftwareProvider
	iInfo.IssuerID = isserInfo.T360ID
	iInfo.ClientID = isserInfo.T360ID

	return iInfo, nil

}

func FromIssuerName(ctx context.Context, operatorName string, fs *firestore.Client) (IssuerInformation, error) {
	iInfo := IssuerInformation{}
	isserInfo := struct {
		Issuer           string `firestore:"issuer"`
		T360ID           string `firestore:"t360_id"`
		SoftwareProvider int    `firestore:"software_provider"`
	}{}

	itr := fs.Collection("registered_issuers").Where("issuer", "==", operatorName).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorf("FromOperatorsName:[%s]:%v", operatorName, err)
				return iInfo, err
			}
		} else {
			if doc.Exists() {
				err = doc.DataTo(&isserInfo)
				if err != nil {
					log.Errorf("FromOperatorsName:[%s]:%v", operatorName, err)
					return iInfo, err
				} else {
					break
				}
			}
		}
	}

	if len(isserInfo.Issuer) == 0 {

		return iInfo, fmt.Errorf("%w with name %s", ErrIssuerNotFound, operatorName)
	}

	iInfo.T360ID = isserInfo.T360ID
	iInfo.Issuer = isserInfo.Issuer
	iInfo.SoftwareID = isserInfo.SoftwareProvider
	iInfo.IssuerID = isserInfo.T360ID
	iInfo.ClientID = isserInfo.T360ID

	return iInfo, nil
}
