// Netlify Functions entrypoint (Lambda-compatibility mode).
// Wraps the existing http.Handler so the same code serves Cloud Run and Netlify.
package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/ignite-hq/installer/handler"
)

const functionPrefix = "/.netlify/functions/installer"

func main() {
	c := handler.DefaultConfig
	if user := os.Getenv("INSTALLER_USER"); user != "" {
		c.User = user
	}
	c.Token = os.Getenv("GH_TOKEN")
	adapter := httpadapter.New(&handler.Handler{Config: c})
	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// direct function URLs carry the function prefix; the /* rewrite does not
		req.Path = strings.TrimPrefix(req.Path, functionPrefix)
		if req.Path == "" {
			req.Path = "/"
		}
		return adapter.ProxyWithContext(ctx, req)
	})
}
