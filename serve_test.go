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

// func TestServer(t *testing.T) {
// 	parent := context.Background()
// 	lg := log.Default()

// 	server := NewServer(parent, ":8082", lg, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
// 	mux, err := server.ServeMux()
// 	if err != nil {
// 		t.Log(err)
// 	}
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "Hello World!")
// 	})
// 	server.ListenAndServe(parent)
// }

func TestCxHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(sample))
	w := httptest.NewRecorder()
	handler := HandlerFunc(CxHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	io.Copy(os.Stdout, resp.Body)
	t.Log(resp)

}

func CxHandler(res *WebhookResponse, req *WebhookRequest) error {
	res.AddTextResponse("With much technolove from Yvan J. Aquino - I wrote this!")
	return nil
}
