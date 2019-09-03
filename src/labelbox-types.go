package main

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
