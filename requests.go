package ezcx

import (
	"io"
	"net/http"

	cx "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"google.golang.org/protobuf/encoding/protojson"
)

type WebhookRequest struct {
	cx.WebhookRequest
}

func NewWebhookRequest() *WebhookRequest {
	return new(WebhookRequest)
}

func WebhookRequestFromReader(rd io.Reader) (*WebhookRequest, error) {
	var req WebhookRequest
	b, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	err = protojson.Unmarshal(b, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func  WebhookRequestFromRequest(r *http.Request) (*WebhookRequest, error) {
	return WebhookRequestFromReader(r.Body)
}

func (req *WebhookRequest) ReadReader(rd io.Reader) error {
	b, err := io.ReadAll(rd)
	if err != nil {
		return err
	}
	err = protojson.Unmarshal(b, req)
	if err != nil {
		return err
	}
	return nil
}

func (req *WebhookRequest) ReadRequest(r *http.Request) error {
	return req.ReadReader(r.Body)
}

func (req *WebhookRequest) PrepareResponse() *WebhookResponse {
	resp := NewWebhookResponse()

	// Added 2022-09-19
	if req.PageInfo != nil {
		resp.PageInfo = req.PageInfo
	}

	resp.SessionInfo.Session = req.SessionInfo.Session
	if req.SessionInfo.Parameters != nil {
		resp.SessionInfo.Parameters = req.SessionInfo.Parameters
	}

	if req.Payload != nil {
		resp.Payload = req.Payload
	}

	return resp
}

func (req *WebhookRequest) GetSessionParameters() (map[string]any, error) {
	params := make(map[string]any)

	// Just in case - I don't think we can iterate over a nil map.
	if req.SessionInfo.Parameters == nil {
		return params, nil
	}

	for key, protoVal := range req.SessionInfo.Parameters {
		params[key] = protoVal.AsInterface()
	}

	return params, nil
}

func (req *WebhookRequest) GetPayload() (map[string]any, error) {
	params := make(map[string]any)

	// Just in case - I don't think we can iterate over a nil map.
	if req.Payload == nil {
		return params, nil
	}
	if req.Payload.Fields == nil {
		return params, nil
	}

	for key, protoVal := range req.Payload.Fields {
		params[key] = protoVal.AsInterface()
	}

	return params, nil
}
