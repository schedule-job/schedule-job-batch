package request

import (
	"github.com/schedule-job/schedule-job-batch/internal/rule_based_replace"
)

type DefaultRequest struct {
	Request
	ID       string              `json:"id"`
	Url      string              `json:"url"`
	Method   string              `json:"method"`
	Body     string              `json:"body"`
	Handlers map[string][]string `json:"headers"`
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
func (r DefaultRequest) GetHandlers() map[string][]string {
	return r.Handlers
}

func NewDefaultRequest(id, url, method, body string, handlers map[string][]string) DefaultRequest {
	return DefaultRequest{
		ID:       id,
		Url:      url,
		Method:   method,
		Body:     body,
		Handlers: handlers,
	}
}

func NewDefaultRequestByInterface(data interface{}) DefaultRequest {
	m := data.(map[string]interface{})
	return NewDefaultRequest(m["id"].(string), m["url"].(string), m["method"].(string), m["body"].(string), m["headers"].(map[string][]string))
}
