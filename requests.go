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
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	cx "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// 2022-11-28
type CxParameterType int

const (
	String CxParameterType = iota
	Integer
	Float
)

type WebhookRequest struct {
	cx.WebhookRequest
	// 2022-10-08: Replaced context.Context with func () context.Context.
	ctx func() context.Context
}

func NewWebhookRequest() *WebhookRequest {
	return new(WebhookRequest)
}

// WebhookRequest Initializations

// Initialize the PageInfo field
func (req *WebhookRequest) initPageInfo() {
	if req.PageInfo == nil {
		req.PageInfo = new(cx.PageInfo)
	}
	if req.PageInfo.FormInfo == nil {
		req.PageInfo.FormInfo = new(cx.PageInfo_FormInfo)
	}
	if req.PageInfo.FormInfo.ParameterInfo == nil {
		req.PageInfo.FormInfo.ParameterInfo = make([]*cx.PageInfo_FormInfo_ParameterInfo, 0)
	}
}

// Initialize the SessionInfo field
func (req *WebhookRequest) initSessionInfo() {
	if req.SessionInfo == nil {
		req.SessionInfo = new(cx.SessionInfo)
	}
}

// Initialize the Payload field
func (req *WebhookRequest) initPayload() {
	if req.Payload == nil {
		req.Payload = new(structpb.Struct)
	}
	if req.Payload.Fields == nil {
		req.Payload.Fields = make(map[string]*structpb.Value)
	}
}

func (req *WebhookRequest) Context() context.Context {
	return req.ctx()
}

// .
func (req *WebhookRequest) Logger() *log.Logger {
	ctx := req.Context()
	ctxLg := ctx.Value(Logger)
	if ctxLg == nil {
		// During testing, it's possible the user defined logger was not
		// flowed down.  This is provided for convenience.
		return log.Default()
	}
	lg, ok := ctxLg.(*log.Logger)
	if !ok {
		return log.Default()
	}
	return lg
}

// Sets (overrides) the PageInfo.ParameterInfos to match the provided map m
func (req *WebhookRequest) setPageFormParameters(m map[string]any) error {
	params := make([]*cx.PageInfo_FormInfo_ParameterInfo, 0)
	for k, v := range m {
		var formParameter cx.PageInfo_FormInfo_ParameterInfo
		pv, err := anyToProto(v)
		if err != nil {
			return err
		}
		formParameter.DisplayName = k
		formParameter.Value = pv
		formParameter.State = cx.PageInfo_FormInfo_ParameterInfo_FILLED
		params = append(params, &formParameter)
	}
	req.PageInfo.FormInfo.ParameterInfo = params
	return nil
}

// Sets (overrides) the SessionInfo.Parameters to match the provided map m
func (req *WebhookRequest) setSessionParameters(m map[string]any) error {
	pm, err := anyToProtoMap(m)
	if err != nil {
		return err
	}
	req.SessionInfo.Parameters = pm
	return nil
}

// Sets (overrides the SessionInfo.Parameters) to match the provided map m
func (req *WebhookRequest) setPayload(m map[string]any) error {
	pm, err := anyToProtoMap(m)
	if err != nil {
		return err
	}
	req.Payload.Fields = pm
	return nil
}

// yaquino@2022-10-11: Dialogflow CX API May include "extra" fields that may
// throw errors and interface with protojson.Unmarshal.  As per the documentation,
// these fields may be ignored. Now also pointing at req.WebhookRequest for unmarshalling..
func WebhookRequestFromReader(rd io.Reader) (*WebhookRequest, error) {
	var req WebhookRequest
	b, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	unmarshaler := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	err = unmarshaler.Unmarshal(b, &req.WebhookRequest)
	if err != nil {
		return nil, ErrUnmarshalWrapper("WebhookRequestFromReader", err)
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
	return req, nil
}

func (req *WebhookRequest) ReadReader(rd io.Reader) error {
	b, err := io.ReadAll(rd)
	if err != nil {
		return err
	}
	unmarshaler := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	err = unmarshaler.Unmarshal(b, &req.WebhookRequest)
	if err != nil {
		return err
	}
	return nil
}

func (req *WebhookRequest) ReadRequest(r *http.Request) error {
	return req.ReadReader(r.Body)
}

// Is this the right format??
func (req *WebhookRequest) WriteRequest(w io.Writer) error {
	m := protojson.MarshalOptions{Indent: "\t"}
	b, err := m.Marshal(&req.WebhookRequest)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (req *WebhookRequest) InitializeResponse() *WebhookResponse {
	return req.initializeResponse()
}

func (req *WebhookRequest) initializeResponse() *WebhookResponse {
	resp := NewWebhookResponse()
	return req.copySession(resp)
}

func (req *WebhookRequest) copySession(res *WebhookResponse) *WebhookResponse {
	if res.SessionInfo == nil {
		res.SessionInfo = new(cx.SessionInfo)
	}
	res.SessionInfo.Session = req.SessionInfo.Session
	return res
}

func (req *WebhookRequest) CopyPageInfo(res *WebhookResponse) {
	if req.PageInfo != nil {
		res.PageInfo = req.PageInfo
	}
}

func (req *WebhookRequest) CopySessionInfo(res *WebhookResponse) *WebhookResponse {
	if req.SessionInfo != nil {
		res.SessionInfo = req.SessionInfo
	}
	return res
}

func (req *WebhookRequest) CopyPayload(res *WebhookResponse) *WebhookResponse {
	if req.Payload != nil {
		res.Payload = req.Payload
	}
	return res
}

func (req *WebhookRequest) GetPageFormParameters() map[string]any {
	params := make(map[string]any)

	// Just in case - I don't think we can iterate over a nil map.
	if req.PageInfo == nil {
		return nil
	}
	if req.PageInfo.FormInfo == nil {
		return nil
	}
	if req.PageInfo.FormInfo.ParameterInfo == nil {
		return nil
	}

	for _, paramInfo := range req.PageInfo.FormInfo.ParameterInfo {
		params[paramInfo.DisplayName] = protoToAny(paramInfo.Value)
	}

	return params
}

func (req *WebhookRequest) GetSessionParameters() map[string]any {
	if req.SessionInfo == nil {
		return nil
	}
	if req.SessionInfo.Parameters == nil {
		return nil
	}
	return protoToAnyMap(req.SessionInfo.Parameters)
}

func (req *WebhookRequest) GetSessionParameter(key string) (any, bool) {
	// Check if SessionInfo Parameters is nil.
	if req.SessionInfo == nil {
		return nil, false
	}
	if req.SessionInfo.Parameters == nil {
		return nil, false
	}
	pv, ok := req.SessionInfo.Parameters[key]
	return protoToAny(pv), ok
}

func (req *WebhookRequest) GetPayload() map[string]any {
	if req.Payload == nil {
		return nil
	}
	if req.Payload.Fields == nil {
		return nil
	}
	return protoToAnyMap(req.Payload.Fields)
}

func (req *WebhookRequest) GetPayloadParameter(key string) (any, bool) {
	// Just in case - I don't think we can iterate over a nil map.

	if req.Payload == nil {
		return nil, false
	}
	if req.Payload.Fields == nil {
		return nil, false
	}

	pv, ok := req.Payload.Fields[key]
	return protoToAny(pv), ok
}

// Testing
func NewTestingWebhookRequest(session, payload, pageform map[string]any) (*WebhookRequest, error) {
	return NewWebhookRequest().initTestingWebhookRequest(session, payload, pageform)
}

func (req *WebhookRequest) initTestingWebhookRequest(session, payload, pageform map[string]any) (*WebhookRequest, error) {
	// Provided for testing, normally http.Request.Context is flowed down.
	req.ctx = context.Background

	// All incoming WebhookRequests should have a session.
	req.initSessionInfo()
	req.SessionInfo.Session = uuid.New().String()

	// if session parameters are provided...
	if session != nil {
		err := req.setSessionParameters(session)
		if err != nil {
			return nil, err
		}
	}

	// if payload parameters are provided...
	if payload != nil {
		req.initPayload()
		err := req.setPayload(payload)
		if err != nil {
			return nil, err
		}
	}

	// if pageForm parameters are provided...
	if pageform != nil {
		req.initPageInfo()
		err := req.setPageFormParameters(pageform)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

// yaquino: 2022-10-08Review this...!
func (req *WebhookRequest) TestCxHandler(out io.Writer, h HandlerFunc) (*WebhookResponse, error) {
	if req.ctx == nil {
		req.ctx = context.Background
	}
	res := req.initializeResponse()
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
