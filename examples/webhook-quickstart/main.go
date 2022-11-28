// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
