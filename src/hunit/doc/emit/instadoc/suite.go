package instadoc

type Content struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Suite struct {
	Title    string     `json:"title,omitempty"`
	Detail   *Content   `json:"detail,omitempty"`
	Sections []*Section `json:"sections,omitempty"`
	Routes   []*Route   `json:"routes,omitempty"`
}

type Header struct {
	Title  string   `json:"title,omitempty"`
	Detail *Content `json:"detail,omitempty"`
}

type Section struct {
	Title  string   `json:"title,omitempty"`
	Detail *Content `json:"detail,omitempty"`
}

type Route struct {
	Title    string                 `json:"title,omitempty"`
	Detail   *Content               `json:"detail,omitempty"`
	Method   string                 `json:"method,omitempty"`
	Resource string                 `json:"resource,omitempty"`
	Attrs    map[string]interface{} `json:"attrs,omitempty"`
	Params   []*Parameter           `json:"params,omitempty"`
	Examples []*Example             `json:"examples,omitempty"`
}

type Parameter struct {
	Name   string   `json:"name,omitempty"`
	Type   string   `json:"type,omitempty"`
	Detail *Content `json:"detail,omitempty"`
}

type Example struct {
	Title    string   `json:"title,omitempty"`
	Detail   *Content `json:"detail,omitempty"`
	Request  *Listing `json:"request,omitempty"`
	Response *Listing `json:"response,omitempty"`
}

type Listing struct {
	Title  string   `json:"title,omitempty"`
	Detail *Content `json:"detail,omitempty"`
	Data   string   `json:"data,omitempty"`
}
