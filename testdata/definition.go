package _definition

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
		ResponseBodyEncoding st.XML
	}
	SayHi struct {
		st.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  st.JSON
		ResponseBodyEncoding st.JSON
	}
	HelloWorld struct {
		st.API
		RouteIntroduce Introduce `route:"/"`
		RouteSayHi     SayHi     `route:"/say-hi"`
	}
)
