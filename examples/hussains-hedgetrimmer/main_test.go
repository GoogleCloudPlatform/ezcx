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
	"os"
	"testing"

	"github.com/yaq-cc/ezcx"
)

func TestHussainsHedgeTrimmer(t *testing.T) {
	session := make(map[string]any)
	session["name"] = "Hussain"
	session["id"] = 5
	session["isUser"] = true
	session["trimmable"] = "Ugh, I wish I knew what was really... like, really going on!"

	payload := make(map[string]any)
	payload["callerId"] = "+14242556256"

	req, err := ezcx.NewTestWebhookRequest(session, payload)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(req)
	res := req.PrepareResponse()
	err = cxHedgeTrimmer(res, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
	res.Write(os.Stdout)
}
