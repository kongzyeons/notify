package models

type Response[T any] struct {
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}
