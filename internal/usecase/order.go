package usecase

import (
	"github.com/Berchon/Clean-Architecture/internal/entity"
	"github.com/Berchon/Clean-Architecture/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type OrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderEvent      events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderEvent events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *OrderUseCase {
	return &OrderUseCase{
		OrderRepository: OrderRepository,
		OrderEvent:      OrderEvent,
		EventDispatcher: EventDispatcher,
	}
}

func (c *OrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderEvent.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderEvent)

	return dto, nil
}

func (c *OrderUseCase) GetOrders() ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.GetOrders()
	if err != nil {
		return []OrderOutputDTO{}, err
	}

	outputDto := []OrderOutputDTO{}
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		outputDto = append(outputDto, dto)
	}

	c.OrderEvent.SetPayload(outputDto)
	c.EventDispatcher.Dispatch(c.OrderEvent)

	return outputDto, nil
}
