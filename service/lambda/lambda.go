package lambda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	aws_event "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/run"
)

type APIResponse struct {
	Data   []byte `json:"data"`
	Code   int64  `json:"code,omitempty"`
	ErrMsg string `json:"err_msg,omitempty"`
}

type Handler func(request aws_event.APIGatewayProxyRequest) (aws_event.APIGatewayProxyResponse, error)

var sampleHandler Handler = func(request aws_event.APIGatewayProxyRequest) (aws_event.APIGatewayProxyResponse, error) {
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf("this is message from %s", hostname)
	return responseSuccess([]byte(msg)), nil
}

func responseSuccess(data []byte) aws_event.APIGatewayProxyResponse {
	resp := APIResponse{
		Data: data,
		Code: http.StatusOK,
	}
	bs, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		Body:       string(bs),
		StatusCode: http.StatusOK,
	}
}

func RegisterRuntime(g *run.Group) {
	g.Add(func() error {
		aws_lambda.Start(sampleHandler)
		return nil
	}, func(err error) {})
}
