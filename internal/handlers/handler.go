package handlers

import "neiro-api/internal/services"

type Handler struct {
	Services *services.Service
}

func NewHandler() *Handler {
	return &Handler{
		Services: services.NewService(),
	}
}
