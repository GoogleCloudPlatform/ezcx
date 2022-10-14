# Getting started with ezcx
## What (exactly) is ezcx?
`ezcx` is a go library that facilitates the development and productionizing of Dialogflow CX Webhook fulfillment APIs.

`ezcx` provides ease-of-use facades built on top of Google's official Dialogflow CX gRPC messages.  A facade is a structure that acts as a substitute for some other target object, providing a simpler and more approachable interface.  For Example, the `ezcx.WebhookRequest` and `ezcx.WebhookResponse` structures are proxies for the cx.WebhookRequest and cx.WebhookResponse gRPC messages.  

```go
// WebhookRequest is a facade. 
type WebhookRequest struct {
	cx.WebhookRequest // cx.WebhooKRequest is the actual gRPC message.
	ctx func() context.Context
}
```

By using the facade pattern, `ezcx` fills the *ease-of-use gap* that is ubiquituously absent with code generated client libraries.  Ease-of-use is defined as:

1. Manual allocation, initialization, and nil-checks: The Dialogflow CX v3 and v3beta1 APIs consistently deal with nested objects. Structural allocations are a major source of developmental setbacks, are often hard to keep track of, and can lead to insiduous bugs.  Furthermore, nil checks, manual allocations, and initializations can be exceedingly verbose - so much so they can actually hinder readability.  `ezcx` aims to take care of this for the developer. 

```go
// AddTextResponses is a good example of why ezcx makes sense.  Imagine having to re-write all this boiler plate every time you have a new virtual agent project..!

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

2. Working with the standard library.  When undertaking a proper Dialogflow CX Webhook Fulfillment API project, if the developer chooses to use the Dialogflow CX Cloud Client library, they may need to undergo a rather lengthy journey to properly explore the ecosystem of extended tooling for interoperating with protobuf messages.  `ezcx` removes the need to interact directly with protobuf structures, marshalling, and unmarshalling - and reduces those interactions to standard Go libraries, data structures, and paradigms.  

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

3. The full solution.  While `ezcx` is modular, it's really designed to serve as a one stop shop for publishing Webhook Fulfillment APIs.  In particular, `ezcx` provides structures for creating an HTTP server that's wired up to support the entire `ezcx` ecosystem.   `ezcx` provides a special type, the `ezcx.HandlerFunc`, that adapts between Webhook handler functions and HTTP handler functions:

```go
// http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ezcx.HandlerFunc type satisfies the http.Handler interface by implementing
// ServeHTTP.
type HandlerFunc func(*WebhookResponse, *WebhookRequest) error

// See serve.go for more details on the actual implementation.
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// OMITTED
}
``` 

In the spirit of providing the full solution, `ezcx` introduces it's own server, the `ezcx.Server` which is just a pre-packaged `http.Server` instance that comes pre-integrated with a variety of useful features:

1. `ezcx.Server` is signals aware and will initiate a properly contexted graceful shutdown upon receiving a SIGINT (interrupt)  or SIGTERM (terminate) signal.  Eventually, `ezcx.Server` may support hot-reconfiguration via the SIGHUP (hang up) signal.  

2. Contexts as a first party consideration.  `ezcx.Server` prefers the user to provide a user instantiated context.  This context is flowed down from the server to the HTTP Handler and from that HTTP Handler to the `ezcx.HandlerFunc`.  Within a traditional http.HandlerFunc, the context is accessible via the (*http.Request).Context method; with `ezcx` we provide that same context via the (*WebhookRequest).Context method which IS a flow down copy of (*http.Request).Context method. See examples/context-and-logging for example usage.

```go
func CxHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	ctx := req.Context() // Access the context, which is a pointer to the (*http.Request).Context method
	res.AddTextResponse("ezcx makes it easy!")
	return nil
}
```

3. The `ezcx.Server` is designed to accept functions that follow `ezcx.HandlerFunc`'s method signature: `(WebhookResponse, WebhookRequest) error` instead of http.HandlerFuncs.  The server adapts the `ezcx.HandlerFunc`s into http.Handlers for the developer. `ezcx.Server` uses a default mux - eventually, the goal would be to support custom routers.  You can provide `ezcx.HandlerFunc`s via the server's HandleCx method:

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/yaq-cc/ezcx"
)

var (
	PORT = os.Getenv("PORT")
)

func main() {
	parent := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(parent, ":"+PORT, lg)
	// Add as many Handlers as you need!
	server.HandleCx("/tell-a-joke", CxJokeHandler) 
	server.HandleCx("/send-a-text", CxTextHandler)
	server.HandleCx("/drink-coffe", CxCoffeeDrinkingHandler)
	server.ListenAndServe(parent)
}
```


## [Actually] getting started with ezcx
### Writing Webhook handlers
`ezcx` let's you focus on codifying business logic.  For instance, if you just need to send back a text response, use the (*WebhookResponse).AddTextResponse method.  

```go
func CxHelloWorldHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	res.AddTextResponse("Hello World!")
	return nil
}
```

If you want to add an outputAudioText response as well, just use AddOutputAudioTextResponse.

```go
func CxHelloWorldHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	res.AddTextResponse("Hello World!")
	res.AddOutputAudioTextResponse("<ssml>Hello World!<ssml>")
	return nil
}
```

You can, in fact, add multiple output ResponseMessages to the same WebhookResponse - but you need to be cautious and aware of some of the major caveats in doing so.

In some cases such as Text and OutputAudioText, if you do provide "multiple" response types to support a multi-channel agent, you need to make sure that fulfillments in the console are configured with those response types as well.  

To read more about why this is important, please see the Dialogflow CX API documentation for the [ResponseMessage message/object](https://cloud.google.com/dialogflow/cx/docs/reference/rest/v3/Fulfillment#ResponseMessage).  

> *Response messages are also used for output audio synthesis. The approach is as follows:*
>
>- *If at least one OutputAudioText response is present, then all OutputAudioText responses __are linearly concatenated__, and the result is used for output audio synthesis.*
>
>- *If the OutputAudioText responses are a mixture of text and SSML, then the concatenated result __is treated as SSML__; otherwise, the result is treated as either text or SSML as appropriate. The agent designer should ideally use either text or SSML consistently throughout the bot design.*
>- *Otherwise, __all Text responses are linearly concatenated__, and the result is used for output audio synthesis.*
> 
> *This approach allows for more sophisticated user experience scenarios, where the text displayed to the user may differ from what is heard.*

Sometimes you don't need to return a fulfillment and instead just add a parameter.  You can use AddSessionParameter / AddSessionParameters or SetSessionParameters to do exactly that.

AddSessionParameters will iterate over the existing session Parameters, updating or adding as necessary whereas SetSessionParameters will set the SessionParameters to the map you provide, overwriting the existing Session parameters.  

The Add / Set parameter process may return an error depending on what gets provided - always check for errors!

```go 
func CxHelloWorldHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	params := map[string]any{
		"ezcx-message": "ezcx is great",
	}
	err := res.AddSessionParameters(params)
	if err != nil {
		return err
	}
	return nil
}
```

Webhook based fulfillment is all about interacting with Dialogflow WebhookRequest provided parameters.  In the previous example we Added / Set session parameters - but what if we need to extract parameters from a given session?  That's what (*WebhookRequest).GetSessionParameters is for!

```go
func CxHelloWorldHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	params := req.GetSessionParameters()
	color, ok := params["color"]
	if !ok {
		res.AddTextResponse("I couldn't find the provided color.")
		return fmt.Errorf("missing session parameter: color")
	}
	// add a parameter
	params["color-processed"] = true
	// delete a parameter
	delete(params, "color")
	// SET the session parameters.  Necessary when deleting!
	err := res.SetSessionParameters(params)
	if err != nil {
		return err
	}
	res.AddTextResponse(fmt.Sprintf("The provided color was %s", color))
	return nil
}
```

The importance of error handling and existence checking can't be overstated.  While it does increase verbosity, it comes with the added benefit of understanding what went wrong.  



### WebhookRequest and WebhookResponse helper methods

### Integrated Contexts and Logging

### Recommended practices for using the ezcx.Server

## Unit testing

