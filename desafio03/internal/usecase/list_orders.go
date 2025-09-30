package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrdersInputDTO struct {
	Page int `json:"page"`
}

type ListOrdersOutputDTO struct {
	Orders []OrderOutputDTO `json:"orders"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute(input ListOrdersInputDTO) (ListOrdersOutputDTO, error) {
	orders, err := c.OrderRepository.List(input.Page)
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}
	var output ListOrdersOutputDTO
	for _, order := range orders {
		output.Orders = append(output.Orders, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}
	return output, nil
}
