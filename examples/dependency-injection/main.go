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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yaq-cc/ezcx"
)

var (
	PORT = os.Getenv("PORT")
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
type Dependencies struct{}

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
func cxHello(d *Dependencies) ezcx.HandlerFunc {
	return func(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
		res.AddTextResponse("It's ... really this easy.")
		return nil
	}
}
