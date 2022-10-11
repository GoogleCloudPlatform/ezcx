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
	"encoding/json"
	"io"
	"net/http"
)

type joke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func jokeFromReader(r io.Reader) (*joke, error) {
	var j joke
	err := json.NewDecoder(r).Decode(&j)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

type jokesClient struct {
	*http.Client
}

func newJokesClient() *jokesClient {
	return &jokesClient{
		Client: &http.Client{
			Transport: &jokesTransport{},
		},
	}
}

func (c *jokesClient) get(ctx context.Context) (*joke, error) {
	URL := "https://icanhazdadjoke.com/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	joke, err := jokeFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return joke, nil
}

type jokesTransport struct{}

func (t *jokesTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	return http.DefaultTransport.RoundTrip(req)
}

var defaultJokesClient = newJokesClient()