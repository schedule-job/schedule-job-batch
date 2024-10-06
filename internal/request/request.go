package request

import (
	"errors"
)

type RequestOption struct {
}

type RequestFormat struct {
	ID       string                 `json:"id"`
	Url      string                 `json:"url"`
	Method   string                 `json:"method"`
	Body     string                 `json:"body"`
	Handlers map[string]interface{} `json:"headers"`
}

type RequestInterface interface {
	GetID() string
	GetUrl() string
	GetMethod() string
	GetBody() string
	GetHandlers() map[string]interface{}
}

func NewRequestFormat(r RequestInterface) RequestFormat {
	var request = RequestFormat{}

	request.ID = r.GetID()
	request.Url = r.GetUrl()
	request.Method = r.GetMethod()
	request.Body = r.GetBody()
	request.Handlers = r.GetHandlers()

	return request
}

type Request struct {
	Options        RequestOption
	requests       map[string]func(interface{}, interface{}) (RequestInterface, error)
	requestOptions map[string]interface{}
}

func (r *Request) CheckSupportedRequest(name string) bool {
	return r.requests[name] != nil
}

func (r *Request) AddRequest(name string, f func(interface{}, interface{}) (RequestInterface, error), options interface{}) {
	if r.requests == nil {
		r.requests = make(map[string]func(interface{}, interface{}) (RequestInterface, error))
	}
	if r.requestOptions == nil {
		r.requestOptions = make(map[string]interface{})
	}
	r.requests[name] = f
	r.requestOptions[name] = options
}

func (r *Request) Request(name string, params interface{}) (RequestInterface, error) {
	if r.requests[name] == nil {
		return nil, errors.New("schedule not found")
	}
	return r.requests[name](params, r.requestOptions[name])
}

var Requester = Request{}
