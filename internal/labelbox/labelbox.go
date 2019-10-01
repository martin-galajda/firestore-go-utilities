package labelbox

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/martin-galajda/firestore-go-utilities/internal/uuid"
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

func (s *LabelboxLabelSettings) AddToolDefinition(mid, labelName, translatedLabel string) {
	featSchemaUUID := uuid.MakeUUID()
	schemaNodeUUID := uuid.MakeUUID()

	name := labelName + "(" + translatedLabel + ")"

	s.Tools = append(s.Tools, &LabelboxToolDef{
		Mid:             mid,
		Name:            name,
		Color:           randomcolor.GetRandomColorInHex(),
		Tool:            LabelboxToolRectangle,
		FeatureSchemaID: featSchemaUUID,
		SchemaNodeID:    schemaNodeUUID,
	})
}
