package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/epithet-ssh/epithet-oidc-policy/pkg/authenticator"
	"github.com/epithet-ssh/epithet-oidc-policy/pkg/authorizer"
	"github.com/epithet-ssh/epithet-oidc-policy/pkg/policyserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	jwksURL := os.Getenv("JWKS_URL")
	issuer := os.Getenv("ISSUER")
	audienceString := os.Getenv("AUDIENCE")
	audience := strings.Split(audienceString, ",")
	authorizerCommand := os.Getenv("AUTHORIZER_COMMAND")

	authenticator, err := authenticator.New(jwksURL, issuer, audience)
	if err != nil {
		panic(err)
	}

	authorizer, err := authorizer.New(authorizerCommand)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/*", policyserver.New(authenticator, authorizer))

	adapter := httpadapter.New(r)
	h := handler{
		mux: adapter,
	}
	lambda.Start(h.Handle)
}

type handler struct {
	mux *httpadapter.HandlerAdapter
}

// Handle handles lambda invocations :-)
func (h *handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return h.mux.ProxyWithContext(ctx, req)
}
