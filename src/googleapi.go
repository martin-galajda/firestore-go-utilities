package main

import (
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

func initGoogleApiClients() (context.Context, *firestore.Client, *translate.Client) {
	var firebaseClientOpt = option.WithCredentialsFile(*pathToFirebaseConfigFile)

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(
		ctx,
		*firebaseProjectID,
		firebaseClientOpt,
	)

	if err != nil {
		log.Fatalf("Error occurred when initiating firestore client: %v.", err)
	}

	translateClient, err := translate.NewClient(
		ctx,
		firebaseClientOpt,
	)

	if err != nil {
		log.Fatalf("Error occurred when initiating translate client: %v.", err)
	}

	return ctx, firestoreClient, translateClient
}
