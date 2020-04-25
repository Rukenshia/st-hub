package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	resty "github.com/go-resty/resty/v2"
	"github.com/sahilm/fuzzy"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type MinimalShip struct {
	Name   string `json:"name"`
	ShipID uint   `json:"id"`
}

type MinimalShips []MinimalShip

func (s MinimalShips) String(i int) string {
	return s[i].Name
}

func (s MinimalShips) Len() int {
	return len(s)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	api := &WowsAPI{
		applicationID: "2ecce5b4b0ffcffc5e7bc04131fb5c8e",
		realm:         "eu",
		client:        resty.New(),
	}

	ships, err := api.GetWarships()
	if err != nil {
		log.Panicf("Could not get all warships: %v", err)
	}

	shipsList := MinimalShips{}

	for _, ship := range ships {
		shipsList = append(shipsList, MinimalShip{
			Name:   ship.Name,
			ShipID: ship.ShipID,
		})
	}

	names := strings.Split(request.Body, ",")

	filteredShips := MinimalShips{}
	for _, name := range names {
		matches := fuzzy.FindFrom(name, shipsList)

		if len(matches) == 0 {
			return Response{StatusCode: 500}, fmt.Errorf("Ship %s not found", name)
		}

		filteredShips = append(filteredShips, shipsList[matches[0].Index])
	}

	body, err := json.Marshal(filteredShips)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {

	lambda.Start(Handler)
}
