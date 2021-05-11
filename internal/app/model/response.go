package model

type Response struct {
	Status string `json:"status"`
}

const (
	ResponseStatusSuccess = "success"
)
