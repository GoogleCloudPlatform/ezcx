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
	"io"

	cx "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type WebhookResponse struct {
	cx.WebhookResponse
}

func NewWebhookResponse() *WebhookResponse {
	res := new(WebhookResponse)
	return res
}

func (res *WebhookResponse) initializeFulfillments() {
	if res.FulfillmentResponse == nil {
		res.FulfillmentResponse = new(cx.WebhookResponse_FulfillmentResponse)
	}
	if res.FulfillmentResponse.Messages == nil {
		res.FulfillmentResponse.Messages = make([]*cx.ResponseMessage, 0)
	}
}

func (res *WebhookResponse) initializePageInfo() {
	if res.PageInfo == nil {
		res.PageInfo = new(cx.PageInfo)
	}
	if res.PageInfo.FormInfo == nil {
		res.PageInfo.FormInfo = new(cx.PageInfo_FormInfo)
	}
	if res.PageInfo.FormInfo.ParameterInfo == nil {
		res.PageInfo.FormInfo.ParameterInfo = make([]*cx.PageInfo_FormInfo_ParameterInfo, 0)
	}
}

func (res *WebhookResponse) initializeSessionInfo() {
	if res.SessionInfo == nil {
		res.SessionInfo = new(cx.SessionInfo)
	}
	if res.SessionInfo.Parameters == nil {
		res.SessionInfo.Parameters = make(map[string]*structpb.Value)
	}
}

func (res *WebhookResponse) initializePayload() {
	if res.Payload == nil {
		res.Payload = new(structpb.Struct)
	}
	if res.Payload.Fields == nil {
		res.Payload.Fields = make(map[string]*structpb.Value)
	}
}

func (res *WebhookResponse) SetSessionParameters(m map[string]any) error {
	res.initializeSessionInfo()
	pm, err := anyToProtoMap(m)
	if err != nil {
		return err
	}
	res.SessionInfo.Parameters = pm
	return nil
}

func (res *WebhookResponse) AddSessionParameters(m map[string]any) error {
	res.initializeSessionInfo()
	for k, v := range m {
		pv, err := anyToProto(v)
		if err != nil {
			return err
		}
		res.SessionInfo.Parameters[k] = pv
	}
	return nil
}

func (res *WebhookResponse) AddTextResponse(txts ...string) {
	res.initializeFulfillments()
	respMessage := &cx.ResponseMessage{}
	respMessage.Message = &cx.ResponseMessage_Text_{
		Text: &cx.ResponseMessage_Text{
			Text: txts,
		},
	}
	res.FulfillmentResponse.Messages = append(res.FulfillmentResponse.Messages, respMessage)
}

func (res *WebhookResponse) AddOutputAudioTextResponse(ssml string) {
	res.initializeFulfillments()
	respMessage := &cx.ResponseMessage{}
	respMessage.Message = &cx.ResponseMessage_OutputAudioText_{
		OutputAudioText: &cx.ResponseMessage_OutputAudioText{
			Source: &cx.ResponseMessage_OutputAudioText_Ssml{
				Ssml: ssml,
			},
		},
	}
	res.FulfillmentResponse.Messages = append(res.FulfillmentResponse.Messages, respMessage)
}

func (res *WebhookResponse) AddTelephonyTransferResponse(phnum string) {
	res.initializeFulfillments()
	respMessage := &cx.ResponseMessage{}
	respMessage.Message = &cx.ResponseMessage_TelephonyTransferCall_{
		TelephonyTransferCall: &cx.ResponseMessage_TelephonyTransferCall{
			Endpoint: &cx.ResponseMessage_TelephonyTransferCall_PhoneNumber{
				PhoneNumber: phnum,
			},
		},
	}
	res.FulfillmentResponse.Messages = append(res.FulfillmentResponse.Messages, respMessage)
}

func (res *WebhookResponse) SetPayload(m map[string]any) error {
	res.initializePayload()
	pm, err := anyToProtoMap(m)
	if err != nil {
		return err
	}
	res.Payload.Fields = pm
	return nil
}

func (res *WebhookResponse) AddPayload(m map[string]any) error {
	res.initializePayload()
	for k, v := range m {
		pv, err := anyToProto(v)
		if err != nil {
			return err
		}
		res.Payload.Fields[k] = pv
	}
	return nil
}

func (res *WebhookResponse) WriteResponse(w io.Writer) error {
	m := protojson.MarshalOptions{Indent: "\t"}
	b, err := m.Marshal(res)
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
