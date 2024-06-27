# Structural

structural is a Proof of Concept project, that generates API's with automatic request body parsing into discreet types.

`structural` templates are valid Go files in themselves.

For example:

Given a `structural` template of:


```go
package helloworld

import (
	st "github.com/beeceej/structural"
)

type (

	// IntroduceRequest is a struct that represents the request body containing ID, Name, Age
	// this type may be hooked up to a `structural.Route`
	IntroduceRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// IntroduceRequest is a struct that represents the response body containing a message
	// this type may be hooked up to a `structural.Route`
	IntroduceResponse struct {
		Message string `json:"response"`
	}

	// Introduce is a `structural.Route`. When the structural tool finds this `route`,
	// it will generate a route with the following:
	// - automatic marshalling of request body into the `IntroduceRequest` type.
	// - automatic unmarshalling of response body into the `IntroduceResponse` type.
	Introduce struct {
		st.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
	}

	// SayHi is a `structural.Route`. When the structural tool finds this `route`,
	// it will generate a route with the following:
	// - automatic marshalling of request body into the `IntroduceRequest` type.
	// - automatic unmarshalling of response body into the `IntroduceResponse` type.
	SayHi struct {
		st.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
	}

	// SayHi is a `structural.Route`. This showcases that you may have a generic
	// `map[string]interface` request body rather than a discreet type.,
	GenericHello struct {
		st.Route
		RequestBody          map[string]interface{}
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
		Age                  int `arg:"age"`
	}

	// HelloWorld is a `structural.API`, it bundles together the various routes.
	HelloWorld struct {
		st.API
		RouteIntroduce    Introduce    `route:"/introduction"`
		RouteSayHi        SayHi        `route:"/say-hi"`
		RouteGenericHello GenericHello `route:"/hi"`
	}
)
```

the following code would be generated:

```go
package helloworld

import (
	"encoding/json"
	"net/http"
)

type HelloWorldAPI struct{}

func (a *HelloWorldAPI) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/introduction", a.HandleRouteIntroduce)
	mux.HandleFunc("/say-hi", a.HandleRouteSayHi)
	mux.HandleFunc("/hi", a.HandleRouteGenericHello)
}

func (a *HelloWorldAPI) HandleRouteIntroduce(w http.ResponseWriter, r *http.Request) {
	var (
		rb       IntroduceRequest
		response IntroduceResponse
		err      error
	)

	if err = json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response, err = a.handleRouteIntroduce(rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *HelloWorldAPI) HandleRouteSayHi(w http.ResponseWriter, r *http.Request) {
	var (
		rb       IntroduceRequest
		response IntroduceResponse
		err      error
	)

	if err = json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response, err = a.handleRouteSayHi(rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *HelloWorldAPI) HandleRouteGenericHello(w http.ResponseWriter, r *http.Request) {
	var (
		rb       map[string]interface{}
		response IntroduceResponse
		err      error
	)

	if err = json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response, err = a.handleRouteGenericHello(rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
```

The above code handles boiler plate and allows the user to focus on writing api handlers with discreet request body and response body types.

to see an example of the code in action, clone this repository, then run `bin/example`.

```
~/Code/beeceej/structural main*
λ ./bin/example
{"response":"hi my name is brian, I am 32 years old. my user id is abc-123"}
{"response":"map[age:32 id:1 name:brian random-key:abc]"}
```
