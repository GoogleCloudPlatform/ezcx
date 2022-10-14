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
	parent := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(parent, ":"+PORT, lg)
	server.HandleCx("/tell-a-joke", CxJokeHandler)
	server.ListenAndServe(parent)
}

// Sends a joke upon invocation.. 
func CxJokeHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	lg := req.Logger()   // Access the logger via req.Logger (it's passed as a context value)
	ctx := req.Context() // Access the context, which is a proxy for (*http.Request).Context

	joke, err := defaultJokesClient.get(ctx)
	if err != nil {
		lg.Println(err)
		return err
	}
	lg.Println(joke.Joke) // added for testing purposes!
	res.AddTextResponse(joke.Joke)
	return nil
}

func CxHelloWorldHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	params := req.GetSessionParameters()
	color, ok := params["color"]
	if !ok {
		res.AddTextResponse("I couldn't find the provided color.")
		return fmt.Errorf("missing session parameter: color")
	}
	// add a parameter
	params["color-processed"] = true
	// delete a parameter
	delete(params, "color")
	err := res.SetSessionParameters(params)
	if err != nil {
		return err
	}
	res.AddTextResponse(fmt.Sprintf("The provided color was %s", color))
	return nil
}