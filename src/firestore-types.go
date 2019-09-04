package main

type ProcessedUrlDocumentAnnotatedElementsData struct {
	Url              string `firebase:"url"`
	DataAnnotationID string `firebase:"dataAnnotationId"`
}

type ProcessedUrlDocumentData struct {
	AnnotatedElementsData map[string]ProcessedUrlDocumentAnnotatedElementsData `firebase:"annotatedElementsData"`
}

type ProcessedUrlDocument struct {
	Data ProcessedUrlDocumentData `firebase:"data"`
}
