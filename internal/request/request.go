package request

type RequestFormat struct {
	ID       string              `json:"id"`
	Url      string              `json:"url"`
	Method   string              `json:"method"`
	Body     string              `json:"body"`
	Handlers map[string][]string `json:"headers"`
}

type Request interface {
	GetID() string
	GetUrl() string
	GetMethod() string
	GetBody() string
	GetHandlers() map[string][]string
}

func NewRequestFormat(r Request) RequestFormat {
	var request = RequestFormat{}

	request.ID = r.GetID()
	request.Url = r.GetUrl()
	request.Method = r.GetMethod()
	request.Body = r.GetBody()
	request.Handlers = r.GetHandlers()

	return request
}
