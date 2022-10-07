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
	"context"
	"io"
	"net/http"

	"github.com/google/uuid"
	cx "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type WebhookRequest struct {
	cx.WebhookRequest
	// 2022-10-07: A field for the http.Request context has been added to simplify
	// and re-use the original HTTP requests context in down stream web service calls.
	ctx context.Context
}

func NewWebhookRequest() *WebhookRequest {
	return new(WebhookRequest)
}

func NewEmptyWebhookRequest() *WebhookRequest {
	return new(WebhookRequest).emptyInit()
}

// yaquino@2022-10-08: Updated
func (req *WebhookRequest) emptyInit() *WebhookRequest {
	// Allocate SessionInfo
	req.SessionInfo = new(cx.SessionInfo)
	req.SessionInfo.Session = uuid.New().String()

	// // Allocation PageInfo
	// req.PageInfo = new(cx.PageInfo)
	// req.PageInfo.FormInfo = new(cx.PageInfo_FormInfo)
	// req.PageInfo.FormInfo.ParameterInfo = make([]*cx.PageInfo_FormInfo_ParameterInfo, 0)

	// // Allocate the Payload
	// req.Payload = new(structpb.Struct)
	// req.Payload.Fields = make(map[string]*structpb.Value)

	return req

}

// yaquino@2022-10-08: Integrated pageInfo.
func NewTestingWebhookRequest(session, payload, pageform map[string]any) (*WebhookRequest, error) {
	req := NewEmptyWebhookRequest()

	if session != nil {
		// Field structure allocation
		req.SessionInfo.Parameters = make(map[string]*structpb.Value)
		// Adding parameters
		params, err := req.GetSessionParameters()
		if err != nil {
			return nil, err
		}
		for key, val := range session {
			params[key] = val
		}
		err = req.addSessionParameters(params)
		if err != nil {
			return nil, err
		}
	}

	if payload != nil {
		// Field structure allocation
		req.Payload = new(structpb.Struct)
		req.Payload.Fields = make(map[string]*structpb.Value)
		// Adding parameters
		pay, err := req.GetPayload()
		if err != nil {
			return nil, err
		}
		for key, val := range payload {
			pay[key] = val
		}
		err = req.addPayload(pay)
		if err != nil {
			return nil, err
		}
	}

	if pageform != nil {
		// Field structure allocation
		req.PageInfo = new(cx.PageInfo)
		req.PageInfo.FormInfo = new(cx.PageInfo_FormInfo)
		req.PageInfo.FormInfo.ParameterInfo = make([]*cx.PageInfo_FormInfo_ParameterInfo, 0)
		// Adding parameters
		pageParams, err := req.GetPageFormParameters()
		if err != nil {
			return nil, err
		}
		for key, val := range pageform {
			pageParams[key] = val
		}
		err = req.addPageFormParameters(pageParams)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (req *WebhookRequest) TestCxHandler(out io.Writer, h HandlerFunc) (*WebhookResponse, error) {
	res := req.PrepareResponse()
	err := h(res, req)
	if err != nil {
		return nil, err
	}
	err = res.WriteResponse(out)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (req *WebhookRequest) addSessionParameters(params map[string]any) error {
	for key, val := range params {
		protoVal, err := structpb.NewValue(val)
		if err != nil {
			return err
		}
		req.SessionInfo.Parameters[key] = protoVal
	}
	return nil
}

func (req *WebhookRequest) addPayload(data map[string]any) error {
	if req.Payload.Fields == nil {
		req.Payload.Fields = make(map[string]*structpb.Value)
	}
	for key, val := range data {
		protoVal, err := structpb.NewValue(val)
		if err != nil {
			return err
		}
		req.Payload.Fields[key] = protoVal
	}
	return nil
}

// yaquino@2022-10-08: Complete this method!
func (req *WebhookRequest) addPageFormParameters(data map[string]any) error {
	for key, val := range data {
		var formParameter cx.PageInfo_FormInfo_ParameterInfo
		protoVal, err := structpb.NewValue(val)
		if err != nil {
			return err
		}
		formParameter.DisplayName = key
		formParameter.Value = protoVal
		formParameter.State = cx.PageInfo_FormInfo_ParameterInfo_FILLED
		req.PageInfo.FormInfo.ParameterInfo = append(req.PageInfo.FormInfo.ParameterInfo, &formParameter)
	}
	return nil
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

// yaquino@2022-10-07: Refactored to flow http.Request's context to the
// WebhookRequest instance.
func WebhookRequestFromRequest(r *http.Request) (*WebhookRequest, error) {
	req, err := WebhookRequestFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	req.ctx = r.Context()
	return req, nil
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

	resp.SessionInfo = req.SessionInfo

	// resp.SessionInfo.Session = req.SessionInfo.Session
	// if req.SessionInfo.Parameters != nil {
	// 	resp.SessionInfo.Parameters = req.SessionInfo.Parameters
	// }

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

func (req *WebhookRequest) GetPageFormParameters() (map[string]any, error) {
	params := make(map[string]any)

	// Just in case - I don't think we can iterate over a nil map.
	if req.PageInfo.FormInfo.ParameterInfo == nil {
		return params, nil
	}

	for _, paramInfo := range req.PageInfo.FormInfo.ParameterInfo {
		params[paramInfo.DisplayName] = paramInfo.Value.AsInterface()
	}

	return params, nil
}

func (req *WebhookRequest) GetSessionParameter(key string) (any, bool) {
	// Check if SessionInfo Parameters is nil.
	if req.SessionInfo.Parameters == nil {
		return nil, false
	}

	protoVal, ok := req.SessionInfo.Parameters[key]
	return protoVal.AsInterface(), ok
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

func (req *WebhookRequest) GetPayloadParameter(key string) (any, bool) {
	// Just in case - I don't think we can iterate over a nil map.

	if req.Payload == nil {
		return nil, false
	}
	if req.Payload.Fields == nil {
		return nil, false
	}

	protoVal, ok := req.Payload.Fields[key]
	return protoVal.AsInterface(), ok
}

func (req *WebhookRequest) Context() context.Context {
	return req.ctx
}
