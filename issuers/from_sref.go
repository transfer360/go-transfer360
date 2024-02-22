package issuers

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func FromSearchReference(ctx context.Context, sref string, fs *firestore.Client) (IssuerInformation, error) {

	iInfo := IssuerInformation{}
	issuerid := struct {
		Client struct {
			IssuerID  string `firestore:"issuerid"`
			Issuer_ID string `firestore:"issuer_id"`
		} `firestore:"client"`
	}{}

	itr := fs.Collection("searches").Where("sref", "==", sref).Where("result.is_hirer_vehicle", "==", true).Limit(1).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				log.Errorf("FromSearchReference:[%s]:%v", sref, err)
				return iInfo, err
			}
		} else {
			if doc.Exists() {
				err = doc.DataTo(&issuerid)
				if err != nil {
					log.Errorf("FromSearchReference:[%s]:%v", sref, err)
					return iInfo, err
				}
			}
		}
	}

	issuerID := issuerid.Client.Issuer_ID
	if len(issuerID) == 0 {
		issuerID = issuerid.Client.IssuerID
	}

	if len(issuerID) == 0 {
		return iInfo, fmt.Errorf("issuer not found with search ref [%s]", sref)
	}

	iInfo, err := GetIssuerInformationFromT360ID(ctx, issuerID, fs)
	if err != nil {
		log.Errorf("FromSearchReference:GetIssuerInformationFromT360ID:[%s]:%v", sref, err)
		return iInfo, err
	}

	return iInfo, nil
}
