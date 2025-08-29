package web

import (
	"encoding/json"
	"net/http"

	"github.com/Berchon/Clean-Architecture/internal/entity"
	"github.com/Berchon/Clean-Architecture/internal/usecase"
	"github.com/Berchon/Clean-Architecture/pkg/events"
)

type WebGetOrderHandler struct {
	EventDispatcher events.EventDispatcherInterface
	OrderRepository entity.OrderRepositoryInterface
	OrderEvent      events.EventInterface
}

func NewWebGetOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderEvent events.EventInterface,
) *WebGetOrderHandler {
	return &WebGetOrderHandler{
		EventDispatcher: EventDispatcher,
		OrderRepository: OrderRepository,
		OrderEvent:      OrderEvent,
	}
}

func (h *WebGetOrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	var output []usecase.OrderOutputDTO
	orderUsecase := usecase.NewOrderUseCase(h.OrderRepository, h.OrderEvent, h.EventDispatcher)
	output, err := orderUsecase.GetOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
