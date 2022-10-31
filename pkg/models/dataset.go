package models

type ResponseError struct {
	ErrorDescr string
	Status     int
}

func (c ResponseError) Error() string {
	return c.ErrorDescr
}

type ConfigRequest struct {
	Service  string                 `json:"service"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Version  int                    `json:"version,omitempty"`
	Used     bool                   `json:"-"`
	Extended bool                   `json:"-"`
}
