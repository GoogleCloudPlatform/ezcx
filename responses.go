package ezcx

import (
	"bytes"
	"io"

	cx "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type WebhookResponse struct {
	cx.WebhookResponse
}

func NewWebhookResponse() *WebhookResponse {
	res := new(WebhookResponse)
	res.FulfillmentResponse = new(cx.WebhookResponse_FulfillmentResponse)
	res.FulfillmentResponse.Messages = make([]*cx.ResponseMessage, 0)

	if res.SessionInfo == nil {
		res.SessionInfo = new(cx.SessionInfo)
	}

	if res.SessionInfo.Parameters == nil {
		res.SessionInfo.Parameters = make(map[string]*structpb.Value)
	}

	if res.Payload == nil {
		res.Payload = new(structpb.Struct)
	}

	return res
}

func (res *WebhookResponse) AddSessionParameters(params map[string]any) error {
	
	// Just in case.. - might be more relevant during testing.
	if res.SessionInfo.Parameters == nil {
		res.SessionInfo.Parameters = make(map[string]*structpb.Value)
	}
	for key, val := range params {
		protoVal, err := structpb.NewValue(val)
		if err != nil {
			return err
		}
		res.SessionInfo.Parameters[key] = protoVal
	}
	return nil
}

func (res *WebhookResponse) AddTextResponse(txts ...string) {
	respMessage := &cx.ResponseMessage{}
	respMessage.Message = &cx.ResponseMessage_Text_{
		Text: &cx.ResponseMessage_Text{
			Text: txts,
		},
	}
	res.FulfillmentResponse.Messages = append(res.FulfillmentResponse.Messages, respMessage)
}

func (res *WebhookResponse) AddOutputAudioTextResponse(ssml string) {
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

func (res *WebhookResponse) AddPayload(data map[string]any) error {
	if res.Payload.Fields == nil {
		res.Payload.Fields = make(map[string]*structpb.Value)
	}
	for key, val := range data {
		protoVal, err := structpb.NewValue(val)
		if err != nil {
			return err
		}
		res.Payload.Fields[key] = protoVal
	}
	return nil
}

func (res *WebhookResponse) WriteResponse(w io.Writer) error {
	m := protojson.MarshalOptions{Indent: "    "}
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

func (res *WebhookResponse) Write(w io.Writer) error {
	m := protojson.MarshalOptions{Indent: "    "}
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