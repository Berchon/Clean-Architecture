//go:build wireinject
// +build wireinject

// run: wire in this folder to generate the wire_gen.go

package main

import (
	"database/sql"

	"github.com/Berchon/Clean-Architecture/internal/entity"
	"github.com/Berchon/Clean-Architecture/internal/event"
	"github.com/Berchon/Clean-Architecture/internal/infra/database"
	"github.com/Berchon/Clean-Architecture/internal/infra/web"
	"github.com/Berchon/Clean-Architecture/internal/usecase"
	"github.com/Berchon/Clean-Architecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setCreatedEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setListEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderList,
	wire.Bind(new(events.EventInterface), new(*event.OrderList)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setOrderListEvent = wire.NewSet(
	event.NewOrderList,
	wire.Bind(new(events.EventInterface), new(*event.OrderList)),
)

func NewOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.OrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewOrderUseCase,
	)
	return &usecase.OrderUseCase{}
}

func NewWebCreateOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebCreateOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebCreateOrderHandler,
	)
	return &web.WebCreateOrderHandler{}
}

func NewWebGetOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebGetOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderListEvent,
		web.NewWebGetOrderHandler,
	)
	return &web.WebGetOrderHandler{}
}
