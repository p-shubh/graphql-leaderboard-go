package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Struct for GraphQL response parsing
type TokenData struct {
	TokenName              string `json:"token_name"`
	Description            string `json:"description"`
	CurrentTokenOwnerships []struct {
		OwnerAddress             string `json:"owner_address"`
		LastTransactionTimestamp string `json:"last_transaction_timestamp"`
	} `json:"current_token_ownerships"`
}

type GraphQLResponse struct {
	Data struct {
		CurrentTokenDatasV2 []TokenData `json:"current_token_datas_v2"`
	} `json:"data"`
}

// Function to perform GraphQL query
func performGraphQLRequest(url string, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GraphQLResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Setup() {
	// First query - Beta Test Query
	betaTestQuery := `
		query MyQuery($_lt: timestamp = "2024-01-16T00:00:00.000000") {
			current_token_datas_v2(
				where: {
					collection_id: {_eq: "0x212ee7ca88024f75e20c79dfee04898048fb9de15cb2da27d793151a6d58db25"}, 
					current_token_ownerships: {last_transaction_timestamp: {_lt: $_lt}}
				}
			) {
				token_name
				description
				current_token_ownerships {
					owner_address
					last_transaction_timestamp
				}
			}
		}
	`
	betaTestVariables := map[string]interface{}{
		"_lt": "2024-01-16T00:00:00.000000",
	}
	betaTestURL := "https://indexer-testnet.staging.gcp.aptosdev.com/v1/graphql"

	betaTestResponse, err := performGraphQLRequest(betaTestURL, betaTestQuery, betaTestVariables)
	if err != nil {
		log.Fatalf("Error performing Beta Test query: %v", err)
	}

	fmt.Println("Beta Test Query Results:")
	for _, token := range betaTestResponse.Data.CurrentTokenDatasV2 {
		fmt.Printf("Token Name: %s, Description: %s\n", token.TokenName, token.Description)
		for _, ownership := range token.CurrentTokenOwnerships {
			fmt.Printf("Owner Address: %s, Last Transaction: %s\n", ownership.OwnerAddress, ownership.LastTransactionTimestamp)
		}
	}

	// Second query - Erebrus NFT Query
	erebrusQuery := `
		query MyQuery {
			current_token_datas_v2(
				where: {collection_id: {_eq: "0x465ebc4eb9718e1555976796a4456fa1a2df8126b4e01ff5df7f9d14fb3eba19"}}
			) {
				token_name
				description
				current_token_ownerships {
					owner_address
					last_transaction_timestamp
				}
			}
		}
	`
	erebrusURL := "https://indexer.mainnet.aptoslabs.com/v1/graphql"

	erebrusResponse, err := performGraphQLRequest(erebrusURL, erebrusQuery, nil)
	if err != nil {
		log.Fatalf("Error performing Erebrus NFT query: %v", err)
	}

	fmt.Println("\nErebrus NFT Query Results:")
	for _, token := range erebrusResponse.Data.CurrentTokenDatasV2 {
		fmt.Printf("Token Name: %s, Description: %s\n", token.TokenName, token.Description)
		for _, ownership := range token.CurrentTokenOwnerships {
			fmt.Printf("Owner Address: %s, Last Transaction: %s\n", ownership.OwnerAddress, ownership.LastTransactionTimestamp)
		}
	}
}
