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

package ezcx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCxHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(sample))
	w := httptest.NewRecorder()
	handler := HandlerFunc(CxHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	io.Copy(os.Stdout, resp.Body)
	t.Log(resp)

}

func TestCxHandlerWithWebhookRequestTester(t *testing.T) {
	params := make(map[string]any)
	params["session-parameter-string"] = "My first session parameter"
	params["session-parameter-integer"] = 42
	params["session-parameter-bool"] = true
	req, err := NewTestingWebhookRequest(params, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := req.PrepareResponse()
	err = CxHandler(res, req)
	if err != nil {
		t.Fatal(err)
	}
	err = res.WriteResponse(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func CxHandler(res *WebhookResponse, req *WebhookRequest) error {
	params, err := req.GetSessionParameters()
	if err != nil {
		return err
	}
	_, ok := params["session-parameter-bool"]
	if ok {
		delete(params, "session-parameter-bool")
	}
	res.AddTextResponse("Hello from CxHandler!")
	res.SetSessionParameters(params)
	return nil
}
