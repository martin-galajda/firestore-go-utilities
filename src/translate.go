package main

import (
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func translateLabel(label string) (*string, error) {
	translations, err := translateClient.Translate(ctx, []string{label}, language.Czech, &translate.Options{
		Format: translate.Text,
		Source: language.English,
	})

	if err != nil {
		return nil, err
	}

	translatedWords := []string{}
	for _, translation := range translations {
		translatedWords = append(translatedWords, translation.Text)
	}

	translatedLabel := strings.Join(translatedWords, " ")

	return &translatedLabel, nil
}
