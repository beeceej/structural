package _definition

import "github.com/beeceej/structural"

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
		structural.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  structural.JSON
		ResponseBodyEncoding structural.XML
	}
	SayHi struct {
		structural.Route
		RequestBody          IntroduceRequest
		Response             IntroduceResponse
		RequestBodyEncoding  structural.JSON
		ResponseBodyEncoding structural.JSON
	}
	HelloWorld struct {
		structural.API
		RouteIntroduce Introduce `route:"/"`
		RouteSayHi     SayHi     `route:"/say-hi"`
	}
)
