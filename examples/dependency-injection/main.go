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
	lg.Println(PORT)
	server := ezcx.NewServer(ctx, ":"+PORT, lg)
	deps := NewDependencies()
	server.HandleCx("/confirm", deps.cxConfirm)
	server.HandleCx("/hello", cxHello(deps))
	server.ListenAndServe(ctx)
}

// Dependencies represents access to resources.  The contained resources
// should be safe for concurrent access and use by multiple goroutines.  
// An example of this would be *sql.DB which is a handle for the clients 
// underlying connection pool.
//
// In general, dependencies should provide access to state - but contained
// dependencies should be stateless - i.e. they're meant to provide access
// to state stored separately.   
type Dependencies struct {}

func NewDependencies() *Dependencies {
	return new(Dependencies)
}

// Structural approach.
func (d *Dependencies) cxConfirm(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
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

// Functional approach via closure.
func cxHello(d *Dependencies) func(*ezcx.WebhookResponse, *ezcx.WebhookRequest) error {
	return func(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
		res.AddTextResponse("It's ... really this easy.")
		return nil
	}
}