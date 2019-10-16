package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/martin-galajda/firestore-go-utilities/internal/labelbox"
)

// LabelboxAPI is interface for communicating with Labelbox API
type LabelboxAPI interface {
	ExportProjectDatasetRows(ctx context.Context, projectID string, datasetIds []string) ([]GraphQLDatasetRow, error)
	CreatePredictionModelForProject(ctx context.Context, projectID, modelName string) (string, error)
	CreatePrediction(ctx context.Context, label map[string][]labelbox.LabelboxExportLabel, projectID, predictionModelID, datarowID string) error
	AttachPredictionModel(ctx context.Context, projectID, predictionModelID string) error
}

type labelboxAPI struct {
	client GraphQLClient
}

// NewLabelboxAPI creates new client implementing LabelboxAPI interface
func NewLabelboxAPI(apiURL, apiToken string) LabelboxAPI {
	client := newGraphQLClient(apiURL, apiToken)

	return &labelboxAPI{client: client}
}

// ExportProjectDatasetRows exports project dataset rows
func (api *labelboxAPI) ExportProjectDatasetRows(ctx context.Context, projectID string, datasetIds []string) ([]GraphQLDatasetRow, error) {

	client := api.client

	gqlQuery := `
		query datasetRows(
			$whereProject: WhereUniqueIdInput!,
			$whereDatasets: DatasetWhereInput,
			$skip: Int,
			$pageSizeDatarows: PageSize
		) {
			project(where: $whereProject) {
				datasets(where: $whereDatasets) {
					id
					dataRows(skip: $skip, first:$pageSizeDatarows) {
						id
						externalId
					}
				}
			}
		}
	`

	lastPageFetched := false
	currPage := 0
	pageSize := 100

	result := []GraphQLDatasetRow{}

	for !lastPageFetched {
		responseStruct := struct {
			Project struct {
				Datasets []struct {
					Id       string
					DataRows []GraphQLDatasetRow
				}
			}
		}{}

		gqlVariables := map[string]interface{}{
			"whereProject": map[string]interface{}{
				"id": projectID,
			},
			"whereDatasets": map[string]interface{}{
				"id_in": datasetIds,
			},
			"skip":             currPage * pageSize,
			"pageSizeDatarows": pageSize,
		}

		err := client.doGraphRequest(ctx, gqlQuery, gqlVariables, &responseStruct)

		if err != nil {
			return nil, err
		}

		dataset := responseStruct.Project.Datasets[0]

		result = append(result, dataset.DataRows...)
		lastPageFetched = len(dataset.DataRows) < pageSize
		currPage++
	}

	return result, nil
}

func (api *labelboxAPI) GetActivePredictionModelID(ctx context.Context, projectID string) (string, error) {
	client := api.client

	gqlQuery := `
	query activePredictionModel($projectID: ID!) {
	  project(where: { id: $projectID}) {
			activePredictionModel {
				id
			}
	  }
	}
	`

	responseStruct := struct {
		Project struct {
			ActivePredictionModel struct {
				Id string
			}
		}
	}{}

	gqlVariables := map[string]interface{}{
		"projectID": projectID,
	}

	err := client.doGraphRequest(ctx, gqlQuery, gqlVariables, &responseStruct)

	if err != nil {
		return "", err
	}

	if responseStruct.Project.ActivePredictionModel.Id == "" {
		return "", notFoundError{msg: fmt.Sprintf("Missing prediction model for project with ID %v.", projectID)}
	}

	return responseStruct.Project.ActivePredictionModel.Id, nil
}

func (api *labelboxAPI) createPredictionModel(ctx context.Context, modelName string) (string, error) {

	gqlQuery := `
		mutation ($modelName: String!) {
			createPredictionModel(data:{
				name: $modelName,
				version:1
			}){
				id
			}
		}
	`

	gqlVariables := map[string]interface{}{
		"modelName": modelName,
	}

	responseStruct := struct{
		CreatePredictionModel struct {
			Id string
		}
	}{}

	err := api.client.doGraphRequest(ctx, gqlQuery, gqlVariables, &responseStruct)

	if err != nil {
		return "", err
	}

	return responseStruct.CreatePredictionModel.Id, nil
}

func (api *labelboxAPI) AttachPredictionModel(ctx context.Context, projectID, predictionModelID string) error {
	gqlQuery := `
		mutation attachPredictionModel($projectID: ID!, $predictionModelID: ID!) {
			updateProject(where:{
				id:$projectID
			}, data:{
				activePredictionModel:{
					connect:{
						id: $predictionModelID
					}
				}
			}){
				id
			}
		}	
	`

	gqlVariables := map[string]interface{}{
		"projectID": projectID,
		"predictionModelID": predictionModelID,
	}

	responseStruct := struct{
		Id string
	}{}

	err := api.client.doGraphRequest(ctx, gqlQuery, gqlVariables, &responseStruct)

	if err != nil {
		return err
	}

	return nil
}

func (api *labelboxAPI) CreatePredictionModelForProject(ctx context.Context, projectID, modelName string) (string, error) {
	newPredictionModelID, err := api.createPredictionModel(ctx, modelName)

	if err != nil {
		return "", err
	}

	err = api.AttachPredictionModel(ctx, projectID, newPredictionModelID)

	if err != nil {
		return "", err
	}

	fmt.Printf("Successfully created and attached new prediction model for project with ID %s. Prediction model ID: %v", projectID, newPredictionModelID)

	return newPredictionModelID, nil
}

func (api *labelboxAPI) CreatePrediction(ctx context.Context, label map[string][]labelbox.LabelboxExportLabel, projectID, predictionModelID, datarowID string) error {

	serializedLabel, err := serializeValueToJSONStr(label)

	if err != nil {
		return err
	}

	gqlQuery := `
	mutation createPrediction($label: String!, $predictionModelId: ID!, $projectId: ID!, $datarowId: ID!) {
		createPrediction(data:{
		  label: $label,
		  predictionModelId: $predictionModelId,
		  projectId: $projectId,
		  dataRowId: $datarowId,
		}) {
		  id
		}
	}
	`


	gqlVariables := map[string]interface{}{
		"projectId": projectID,
		"predictionModelId": predictionModelID,
		"datarowId": datarowID,
		"label": serializedLabel,
	}

	var responseStruct struct {
		CreatePrediction struct {
			Id string
		}
	}

	err = api.client.doGraphRequest(ctx, gqlQuery, gqlVariables, &responseStruct)

	if err != nil {
		return err
	}

	return nil

}

func serializeValueToJSONStr(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ExportProjectDatasetRows exports project dataset rows
