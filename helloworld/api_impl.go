package helloworld

import "fmt"

func (a *HelloWorldAPI) handleRouteSayHi(rb IntroduceRequest) (IntroduceResponse, error) {
	return IntroduceResponse{Message: fmt.Sprintf("What up %s", rb.Name)}, nil
}

func (a *HelloWorldAPI) handleRouteIntroduce(rb IntroduceRequest) (IntroduceResponse, error) {
	return IntroduceResponse{Message: fmt.Sprintf("hi my name is %s, I am %d years old. my user id is %s", rb.Name, rb.Age, rb.ID)}, nil

}
func (a *HelloWorldAPI) handleRouteGenericHello(rb map[string]interface{}, age int) (IntroduceResponse, error) {
	return IntroduceResponse{Message: fmt.Sprintf("Hello, World")}, nil
}
