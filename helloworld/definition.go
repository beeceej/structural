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
	}

	// HelloWorld is a `structural.API`, it bundles together the various routes.
	HelloWorld struct {
		st.API
		RouteIntroduce    Introduce    `route:"/introduction"`
		RouteSayHi        SayHi        `route:"/say-hi"`
		RouteGenericHello GenericHello `route:"/hi"`
	}
)
