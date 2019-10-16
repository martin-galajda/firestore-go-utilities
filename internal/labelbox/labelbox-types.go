package labelbox

import "github.com/martin-galajda/firestore-go-utilities/internal/mathutils"

import "encoding/json"

import "fmt"

type LabelboxToolDef struct {
	Mid             string `json:"mid"`
	Name            string `json:"name"`
	Color           string `json:"color"`
	Tool            string `json:"tool"`
	FeatureSchemaID string `json:"featureSchemaId"`
	SchemaNodeID    string `json:"schemaNodeId"`
}

type LabelboxClassificationDef struct{}

type LabelboxLabelSettings struct {
	Tools           []*LabelboxToolDef           `json:"tools"`
	Classifications []*LabelboxClassificationDef `json:"classifications"`
}

var min = mathutils.Min
var max = mathutils.Max

// Structs adapting exported JSON file from Labelbox after annotating images

type LabelboxExportLabelGeometry struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type LabelboxExportLabel struct {
	Geometry []LabelboxExportLabelGeometry `json:"geometry"`
}

type LabelboxExportAnnotation struct {
	ID         string                           `json:"ID"`
	ImageURL   string                           `json:"Labeled Data"`
	Labels     map[string][]LabelboxExportLabel `json:"Label,omitempty"`
	// Labels     interface{} 											`json:"Label,omitempty"`
	ExternalID string                           `json:"External ID"`
}

type labelboxExportAnnotation LabelboxExportAnnotation

func (annotation *LabelboxExportAnnotation) UnmarshalJSON(b []byte) error {
	var lbea labelboxExportAnnotation
	if err := json.Unmarshal(b, &lbea); err != nil {
		fmt.Println(err)
		*annotation = LabelboxExportAnnotation{ID: "Skip"}
		return nil
	}

	*annotation = LabelboxExportAnnotation(lbea)
	return nil
}

// GetBoundingBoxPoints returns Geometry of topleft and bottomright point of the bounding box
func (label *LabelboxExportLabel) GetBoundingBoxPoints() (LabelboxExportLabelGeometry, LabelboxExportLabelGeometry) {
	xCoords := []int32{}
	yCoords := []int32{}

	const MaxInt = int(^uint(0) >> 1)
	const MinInt = -MaxInt - 1


	for _, point := range label.Geometry {
		xCoords = append(xCoords, point.X)
		yCoords = append(yCoords, point.Y)
	}

	leftTopPoint := LabelboxExportLabelGeometry{}
	rightBottomPoint := LabelboxExportLabelGeometry{}

	leftTopPoint.X = min(xCoords...)
	leftTopPoint.Y = min(yCoords...)

	rightBottomPoint.X = max(xCoords...)
	rightBottomPoint.Y = max(yCoords...)

	return leftTopPoint, rightBottomPoint
}
