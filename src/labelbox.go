package main

import (
	"log"

	"github.com/AvraamMavridis/randomcolor"
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
	featSchemaUUID := makeUUID()
	schemaNodeUUID := makeUUID()

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
		FeatureSchemaID: featSchemaUUID,
		SchemaNodeID:    schemaNodeUUID,
	})
}
