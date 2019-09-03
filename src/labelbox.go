package main

import (
	"log"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/google/uuid"
)

func NewLabelboxLabelSettings() *LabelboxLabelSettings {
	tools := []*LabelboxToolDef{}
	classifications := []*LabelboxClassificationDef{}
	settings := &LabelboxLabelSettings{
		Tools:           tools,
		Classifications: classifications,
	}

	return settings
}

func (s *LabelboxLabelSettings) AddToolDefinition(mid, labelName string) {
	featSchemaUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	schemaNodeUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	translatedLabel, err := translateLabel(labelName)

	if err != nil {
		log.Fatalf("Error translating label %q: %v", labelName, err)
	}
	name := labelName + "(" + *translatedLabel + ")"

	s.Tools = append(s.Tools, &LabelboxToolDef{
		Mid:             mid,
		Name:            name,
		Color:           randomcolor.GetRandomColorInHex(),
		Tool:            LabelboxToolRectangle,
		FeatureSchemaID: featSchemaUUID.String(),
		SchemaNodeID:    schemaNodeUUID.String(),
	})
}
