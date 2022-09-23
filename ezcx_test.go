package ezcx

import (
	"os"
	"strings"
	"testing"
)

func TestWebhookRequest(t *testing.T) {

	// whreq := NewWebhookRequest()
	// err := whreq.FromReader(strings.NewReader(sample))
	whreq, err := WebhookRequestFromReader(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(whreq)
}

func TestPrepareResponse(t *testing.T) {
	whreq, err := WebhookRequestFromReader(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n!!!SESSION PARAMS: ", whreq.SessionInfo.Parameters)
	reqParams, err := whreq.GetSessionParameters()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("REQ PARAMS:", reqParams)
	whresp := whreq.PrepareResponse()
	whresp.AddTextResponse("Hello", " World!")
	params := make(map[string]any)
	params["manually_added"] = "Hello from Yvan!"
	params["random_number"] = 6.3
	err = whresp.AddSessionParameters(params)
	if err != nil {
		t.Fatal(err)
	}

	err = whresp.AddPayload(params)
	if err != nil {
		t.Fatal(err)
	}

	whresp.WriteResponse(os.Stdout)
}

var sample = `{
"detectIntentResponseId": "e12be281-028f-4a6b-95c6-9850a27542f1",
"pageInfo": {
	"currentPage": "projects/oktony-cx/locations/global/agents/c5e716ba-9b90-4edc-a792-2ee7dd24b428/flows/2e387ccd-a8f4-4a0e-9cb8-17bad040d8fe/pages/b34fda0b-0769-4f42-b91c-ff38e4bc1268",
	"displayName": "get-national-benchmarks"
},
"sessionInfo": {
	"session": "projects/oktony-cx/locations/global/agents/c5e716ba-9b90-4edc-a792-2ee7dd24b428/sessions/0591c1-9e2-06b-79c-49e9affb8",
	"parameters": {
	"cohort-name": "Back surgery (Spinal fusion)",
	"is-medicare": "medicare",
	"measure-name": "30d-Readmission"
	}
},
"fulfillmentInfo": {
	"tag": "nb-cohorts"
},
"text": "65+",
"languageCode": "en",
"callerId": "+17862512635"
}`
