package registry

// Resource struct
type Resource struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Scalar float64 `json:"value,omitempty"`
	//TODO customer resources
}
