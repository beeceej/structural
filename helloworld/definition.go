package helloworld

import (
	st "github.com/beeceej/structural"
)

type (
	IntroduceRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	IntroduceResponse struct {
		Message string `json:"response"`
	}
	Introduce struct {
		st.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
	}
	SayHi struct {
		st.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
	}
	GenericHello struct {
		st.Route
		RequestBody          map[string]interface{}
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
		Age                  int `arg:"age"`
	}
	HelloWorld struct {
		st.API
		RouteIntroduce    Introduce    `route:"/introduction"`
		RouteSayHi        SayHi        `route:"/say-hi"`
		RouteGenericHello GenericHello `route:"/hi"`
	}
)
