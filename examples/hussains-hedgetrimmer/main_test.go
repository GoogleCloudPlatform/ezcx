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
		t.Log(err)
	}
	var res ezcx.WebhookResponse
	err := cxHedgeTrimmer(*res, req)
	if er != nil {
		t.Log(err)
	}
	res.Write(os.Stdout)
}
