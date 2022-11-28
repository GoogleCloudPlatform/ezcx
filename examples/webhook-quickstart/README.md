# Refactored Webhook Quickstart Example

# Source Code
```go
// ezcx/examples/webhook-quickstart is a refactoring of the Google Cloud provided
// Go webhook quickstart: https://cloud.google.com/dialogflow/cx/docs/quick/webhook 
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/googlecloudplatform/ezcx"
)

var (
	PORT = os.Getenv("PORT")
)

func main() {
	ctx := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(ctx, ":"+PORT, lg)
	server.HandleCx("/confirm", cxConfirm)
	server.ListenAndServe(ctx)
}

func cxConfirm(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	params := req.GetSessionParameters()

	size := params["size"]
	color := params["color"]

	res.AddTextResponse(
		fmt.Sprintf("You can pick up your order for a %s %s shirt in 5 days.",
			size, color),
	)
	params["cancel-period"] = "2"
	res.AddSessionParameters(params)
	return nil
}
```