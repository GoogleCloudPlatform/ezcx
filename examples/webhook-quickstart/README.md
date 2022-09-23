# Refactored Webhook Quickstart Example

# Source Code
```go
// ezcx/examples/webhook-quickstart is a refactoring of the Google Cloud provided
// Go webhook quickstart: https://cloud.google.com/dialogflow/cx/docs/quick/webhook 
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/yaq-cc/ezcx"
)

var (
	PORT = *flag.String("PORT", "8080", "container port to listen to - default is 8080")
)

func main() {
	ctx := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(ctx, ":"+PORT, lg)
	server.HandleCx("/confirm", cxConfirm)
	server.ListenAndServe(ctx)
}

func cxConfirm(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	params, err := req.GetSessionParameters()
	if err != nil {
		return err
	}
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
