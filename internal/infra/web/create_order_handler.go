package web

import (
	"encoding/json"
	"net/http"

	"github.com/Berchon/Clean-Architecture/internal/entity"
	"github.com/Berchon/Clean-Architecture/internal/usecase"
	"github.com/Berchon/Clean-Architecture/pkg/events"
)

type WebCreateOrderHandler struct {
	EventDispatcher events.EventDispatcherInterface
	OrderRepository entity.OrderRepositoryInterface
	OrderEvent      events.EventInterface
}

func NewWebCreateOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderEvent events.EventInterface,
) *WebCreateOrderHandler {
	return &WebCreateOrderHandler{
		EventDispatcher: EventDispatcher,
		OrderRepository: OrderRepository,
		OrderEvent:      OrderEvent,
	}
}

func (h *WebCreateOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderUsecase := usecase.NewOrderUseCase(h.OrderRepository, h.OrderEvent, h.EventDispatcher)
	output, err := orderUsecase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
