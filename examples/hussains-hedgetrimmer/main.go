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
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/ezcx"
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

	params := req.GetSessionParameters()

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
