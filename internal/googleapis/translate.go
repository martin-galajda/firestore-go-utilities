package googleapis

import (
	"log"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// Translator is interface for translating strings with Google API
type Translator interface {
	Translate(label string) (string, error)
}

// Translate translate given string. It returns error in case anything goes wrong.
func (translatorWrapper *translator) Translate(label string) (string, error) {
	translations, err := translatorWrapper.client.Translate(translatorWrapper.ctx, []string{label}, language.Czech, &translate.Options{
		Format: translate.Text,
		Source: language.English,
	})

	if err != nil {
		return "", err
	}

	translatedWords := []string{}
	for _, translation := range translations {
		translatedWords = append(translatedWords, translation.Text)
	}

	translatedLabel := strings.Join(translatedWords, " ")

	return translatedLabel, nil
}

type translator struct {
	client *translate.Client
	ctx    context.Context
}

// NewTranslator creates new Google Translate API Translator.
func NewTranslator(ctx context.Context, pathToConfigFile string) Translator {
	var clientOpt = option.WithCredentialsFile(pathToConfigFile)

	// ctx := context.Background()
	translateClient, err := translate.NewClient(
		ctx,
		clientOpt,
	)

	if err != nil {
		log.Fatalf("Error occurred when initiating translate client: %v.", err)
	}

	return &translator{
		translateClient,
		ctx,
	}
}
