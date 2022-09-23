package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/yaq-cc/ezcx"
)

var (
	PORT = os.Getenv("PORT")
)

func main() {
	ctx := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(ctx, ":"+PORT, lg)
	server.HandleCx("/trimmer", cxHedgeTrimmer)
	server.ListenAndServe(ctx)
}

func cxHedgeTrimmer(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {

	trimmer := strings.NewReplacer(".", "", ",", "", " ", "")

	params, err := req.GetSessionParameters()
	if err != nil {
		return err
	}

	for key, val := range params {
		strVal, ok := val.(string)
		if !ok {
			continue
		}
		params[key] = trimmer.Replace(strVal)
	}

	res.AddSessionParameters(params)
	return nil
}
