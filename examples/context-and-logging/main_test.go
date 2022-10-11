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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/yaq-cc/ezcx"
)

// Unit (logical) testing for CxJokeHandler
func TestCxJokeHandler(t *testing.T) {
	req, err := ezcx.NewTestingWebhookRequest(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := req.InitializeResponse()
	err = CxJokeHandler(res, req)
	if err != nil {
		t.Fatal(err)
	}
	res.WriteResponse(os.Stdout)
}

// Unit (HTTP) testing for CxJokeHandler
func TestServeHTTPJokeHandler(t *testing.T) {
	var buf bytes.Buffer
	req, err := ezcx.NewTestingWebhookRequest(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = req.WriteRequest(&buf)
	if err != nil {
		t.Fatal(err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/tell-a-joke", &buf)
	w := httptest.NewRecorder()
	hf := ezcx.HandlerFunc(CxJokeHandler)
	hf.ServeHTTP(w, httpReq)
	res := w.Result()
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		t.Fatal(err)
	}
}
