package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

type GraphQLClient interface {
	doGraphRequest(ctx context.Context, gqlQuery string, gqlVariables map[string]interface{}, resp interface{}) error
}

type graphQLClient struct {
	*graphql.Client
	authHeader string
}

func (client *graphQLClient) doGraphRequest(ctx context.Context, gqlQuery string, gqlVariables map[string]interface{}, resp interface{}) error {
	req := graphql.NewRequest(gqlQuery)

	for key, val := range gqlVariables {
		req.Var(key, val)
	}

	// set auth header
	req.Header.Set("Authorization", "Bearer "+client.authHeader)

	return client.Client.Run(ctx, req, resp)
}

func newGraphQLClient(apiURL string, apiToken string) GraphQLClient {
	client := &graphQLClient{
		graphql.NewClient(apiURL),
		apiToken,
	}

	client.Client.Log = func(s string) { fmt.Println(s) }

	return client
}
