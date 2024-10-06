package request

import (
	"errors"

	"github.com/schedule-job/schedule-job-batch/internal/rule_based_replace"
)

type DefaultRequest struct {
	RequestInterface
	ID       string                 `json:"id"`
	Url      string                 `json:"url"`
	Method   string                 `json:"method"`
	Body     string                 `json:"body"`
	Handlers map[string]interface{} `json:"headers"`
}

func (r DefaultRequest) GetID() string {
	return r.ID
}
func (r DefaultRequest) GetUrl() string {
	return r.Url
}
func (r DefaultRequest) GetMethod() string {
	return r.Method
}
func (r DefaultRequest) GetBody() string {
	return rule_based_replace.Replacer.RuleBasedReplace(r.Body)
}
func (r DefaultRequest) GetHandlers() map[string]interface{} {
	return r.Handlers
}

func NewDefaultRequest(id, url, method, body string, handlers map[string]interface{}) *DefaultRequest {
	return &DefaultRequest{
		ID:       id,
		Url:      url,
		Method:   method,
		Body:     body,
		Handlers: handlers,
	}
}

func NewDefaultRequestByInterface(payload, _ interface{}) (RequestInterface, error) {
	m, check := payload.(map[string]interface{})
	if !check {
		return nil, errors.New("error")
	}
	return NewDefaultRequest(m["id"].(string), m["url"].(string), m["method"].(string), m["body"].(string), m["headers"].(map[string]interface{})), nil
}
