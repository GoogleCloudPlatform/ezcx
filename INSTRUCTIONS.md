# Getting started with ezcx
## What (exactly) is ezcx?
`ezcx` is a go library that facilitates the development and productionizing of Dialogflow CX Webhook fulfillment APIs.

`ezcx` provides ease-of-use proxies built on top of Google's official Dialogflow CX gRPC messages.  A proxy is a structure that acts as a substitute for some other target object.  For Example, the `ezcx.WebhookRequest` and `ezcx.WebhookResponse` structures are proxies for the cx.WebhookRequest and cx.WebhookResponse gRPC messages.  

```go
// WebhookRequest is a proxy. 
type WebhookRequest struct {
	cx.WebhookRequest // cx.WebhooKRequest is the actual gRPC message.
	ctx func() context.Context
}
```

By using the proxy pattern, `ezcx` fills the *ease-of-use gap* that is generally ubiquituous with code generated client libraries.  Ease-of-use is defined as:

1. Manual allocation, initialization, and nil-checks: The Dialogflow CX v3 and v3beta1 APIs consistently deal with nested objects. Structural allocations are a major source of developmental setbacks, are often hard to keep track of, and can lead to insiduous bugs.  Furthermore, nil checks, manual allocations, and initializations can be exceedingly verbose - so much so they can actually hinder readability.  `ezcx` aims to take care of this for the developer. 

```go
// AddTextResponses is a good example of why it ezcx makes sense.  Imagine having to re-write all this boiler plate every time you have a new virtual agent project..!

func (res *WebhookResponse) AddTextResponse(txts ...string) {
	// Check to see if objects exist else initialize.
    if res.FulfillmentResponse == nil {
		res.FulfillmentResponse = new(cx.WebhookResponse_FulfillmentResponse)
	}
	if res.FulfillmentResponse.Messages == nil {
		res.FulfillmentResponse.Messages = make([]*cx.ResponseMessage, 0)
	}
    // Create a response message object
	respMessage := &cx.ResponseMessage{}
    // Add a message
	respMessage.Message = &cx.ResponseMessage_Text_{
		Text: &cx.ResponseMessage_Text{
			Text: txts,
		},
	}
    // Append message to Fulfillment messages.
	res.FulfillmentResponse.Messages = append(res.FulfillmentResponse.Messages, respMessage)
}
```

2. Working with the standard library.  When undertaking a proper Dialogflow CX Webhook Fulfillment API project, if the developer chooses to use the Dialogflow CX Cloud Client library, they may need to undergo a rather lengthy journey that allows them to properly explore the ecosystem of extended tooling for interoperating with protobuf messages.  `ezcx` removes the need to interact directly with protobuf structures, marshalling, and unmarshalling - and reduces those interactions to  traditional Go libraries, data structures, and paradigms.  

```go
// Note the complete and utter absence of libraries like structpb, protojson, and other gRPC specific libraries.  

func CxHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	log := req.Logger()
	
    // Get session params, pull the caller-name
    params := req.GetSessionParameters()
	callerName, ok := params["caller-name"]
    if !ok {
        log.Println("Session Parameter not found: caller-name")
        return ErrParamNotFound
    }

    // Add a Text response message.
	res.AddTextResponse(fmt.Sprintf("Hi there %s, how are you?", callerName))

    // Add a param and then set the out-going SessionInfo Parameters.
	params["saidHi"] = true 
    err = res.SetSessionParameters(params)
	if err != nil {
		return err
	}
	
    return nil
}
```

3. The full solution.  While `ezcx` is modular, it's really designed to serve as the full solution.  In particular, `ezcx` provides structures for the developer to create a server that's wired up to support the entire `ezcx` ecosystem.   

## [Actually] getting started with ezcx
