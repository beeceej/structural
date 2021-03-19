package helloworld

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	_age := r.URL.Query().Get("age")
	_arg417, err := strconv.Atoi(_age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response, err = a.handleRouteGenericHello(rb, _arg417); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
