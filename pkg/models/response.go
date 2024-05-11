package models

type Meta struct {
	TotalSize int `json:"totalSize"`
}

type Response[T any] struct {
	Data []T   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
