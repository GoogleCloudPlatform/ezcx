# ezcx
`ezcx` is a framework for building containerized Dialogflow CX webhook fulfillment APIs.  `ezcx` runs happiest on Google's Cloud Run service.  

`ezcx` was designed to remove most (if not all) the complexity associated with building Dialogflow CX webhook fulfillment APIs:

- `ezcx` is a convenience wrapper on top of Google Cloud's code-generated gRPC definitions.  `ezcx` exposes wrappers around the WebhookResponse, WebhookRequest, and subsequent protobuf messages used in defining the WebhookRequest and WebhookResponse.

- `ezcx` makes it easy to add WebhookResponse response messages.  The WebhookResponse object has a number of helper methods like `AddTextResponse`, `AddOutputAudioTextResponse` and `AddSessionParameters` that circumvent the need to manage sub-object allocation (forgot to make that map?  PANIC!), programming labor, and the absolute headache of keeping track of deeply nested objects.  

```go
func CxHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	// Add text
    res.AddTextResponse("Made with ezcx in 5 minutes or less...")
    
    // Add SSML
    res.AddOutputAudioTextResponse("<ssml>Made with ezcx in <prosody rate=slow>5 minutes </prosody>or less...</ssml>")    

    // Add Session Parameters
    params = make(map[string]any)
    params["made_with"] = "ezcx"
    res.AddSessionParameters(params)

	return nil
}
```

- `ezcx` is designed to be a full solution for Dialogflow CX Webhook Fulfillment APIs.  You can use ezcx.NewServer to create an http.Server that instance that's wired up to work with functional ezcx.HandlerFunc handlers.  Very much like http.HandleFunc, ezcx.HandlerFunc is an adapter that allows for the definition of CxHandlers of the form `func(*WebhookResponse, *WebhookRequest) error` to be directly used via the server's HandleCx method.  

```go
type HandlerFunc func(*WebhookResponse, *WebhookRequest) error

// Ignore for now.
func (h HandlerFunc) Handle(res *WebhookResponse, req *WebhookRequest) error {
	return h(res, req)
}

// Implements the http.Handler interface - this is relatively low level;
// ezcx "handles" this for you, instead, allowing you to focus on what really matters:
// working with the data in WebhookRequest and returning a WebhookResponse.
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req, err := WebhookRequestFromRequest(r)
	if err != nil {
		log.Println(err)
		return
	}
	res := req.PrepareResponse()
	err = h.Handle(res, req)
	if err != nil {
		log.Println(err)
		return
	}
	res.WriteResponse(w)
}
```

- Creating a web service with `ezcx` couldn't be easier.  the `ezcx.Server` object most (if not all) the features you'd need from a production http.Server instance: signal handling, logging, and graceful shutdowns.  If you see an opportuntiy for improvement, please reach out!

```go
// main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yaq-cc/ezcx"
)

func main() {
    parent := context.Background()
    lg := log.Default()

    server := ezcx.NewServer(parent, ":8082", lg)
    // HandleCx adapts ezcx.HandlerFunc into an http.Handler for you!
    server.HandleCx("/from-dfcx", CxHandler)
    server.ListenAndServe(parent)
}
```

Handlers have been moved to a separate file to show just how little effort is required. server.HandleCx adapts an `ezcx.HandlerFunc` into an `http.Handler` for you!  

```go 
// handlers.go
func CxHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	// Read parameters.
	params, err := req.GetSessionParameters()
	if err != nil {
		return err
	}
	callerName := params["caller-name"]

	// Update your response.
	res.AddTextResponse(fmt.Sprintf("Hi there %s, how are you?", callerName))

	// Update some session parameters
	params["saidHi"] = true
	err = res.AddSessionParameters(params)
	if err != nil {
		return err
	}

	return nil
}
```

# Basic Usage

## Request-scoped Web Service Calls.
`ezcx`'s WebhookRequest flows down `http.Request`'s context; this context is accessible via the WebhookRequest's Context() method. Under the hood, WebhookRequests.Context() method is just a pass-through for (*http.Request).Context().  

```go
func CxHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
  ...
  ctx := req.Context() 
  apiResult, err := makeWebServiceCall(ctx, ...callOpts)
  if err != nil {
    return err
  }
  ...
	return nil
}
```

## Testing
More on testing coming soon!


# Examples
Please visit the examples folder to check out how ezcx stacks up!  

# Dockerfile
Provided for convenience.  

```dockerfile
FROM    golang:1.18-buster as builder
WORKDIR /app
COPY    . ./
RUN     go build -o service

FROM    debian:buster-slim
RUN     set -x && \
		apt-get update && \
		DEBIAN_FRONTEND=noninteractive apt-get install -y \
			ca-certificates && \
			rm -rf /var/lib/apt/lists/*
COPY    --from=builder /app/service /app/service

CMD     ["/app/service"]
```

# Cloud Build
Provided for convenience.  Review all the parameters for deploying to Cloud Run before issuing a gcloud builds submit!

```yaml
steps:
- id: docker-build-push-ezcx-service
  waitFor: ['-']
  name: gcr.io/cloud-builders/docker
  dir: service
  entrypoint: bash
  args:
    - -c
    - |
      docker build -t gcr.io/$PROJECT_ID/${_SERVICE} . &&
      docker push gcr.io/$PROJECT_ID/${_SERVICE}

- id: gcloud-run-deploy-ezcx-service
  waitFor: ['docker-build-push-ezcx-service']
  name: gcr.io/google.com/cloudsdktool/cloud-sdk
  entrypoint: bash
  args:
    - -c
    - |
      gcloud run deploy ${_SERVICE} \
        --project $PROJECT_ID \
        --image gcr.io/$PROJECT_ID/${_SERVICE} \
        --timeout 5m \
        --region ${_REGION} \
        --no-cpu-throttling \
        --min-instances 0 \
        --max-instances 3 \
        --allow-unauthenticated

substitutions:
  _SERVICE: ezcx-service
  _REGION: us-central-1
```

# Backlog
## Testing
## EmptyWebhookRequest 
Review the initialization of emptyWebhookRequest which is used for testing webhooks. I need to add a place for "pageInfo" which is an object that's rarely used. Details on pageInfo here: https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3#PageInfo

# Updates
- 2022-10-07: WebhookRequest now has a method that returns the http.Request's context.  Adding in a Context() method was the simplest and most effective way of providing a request-scoped context to downstream web service calls.
 