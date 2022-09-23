package main

import (
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
	var res ezcx.WebhookResponse
	err = cxHedgeTrimmer(&res, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
