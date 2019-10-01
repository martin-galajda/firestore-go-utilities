package googleapis

import (
	"log"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// NewFirestoreClient creates new Firestore Client.
func NewFirestoreClient(ctx context.Context, firebaseProjectID, pathToConfigFile string) (*firestore.Client, error) {
	var firebaseClientOpt = option.WithCredentialsFile(pathToConfigFile)

	firestoreClient, err := firestore.NewClient(
		ctx,
		firebaseProjectID,
		firebaseClientOpt,
	)

	if err != nil {
		log.Fatalf("Error occurred when initiating firestore client: %v.", err)
	}

	return firestoreClient, err
}
