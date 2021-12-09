package route

// Describes a route in the service under test
type Route struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
