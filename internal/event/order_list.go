package event

import "time"

type OrderList struct {
	Name    string
	Payload interface{}
}

func NewOrderList() *OrderList {
	return &OrderList{
		Name: "OrderList",
	}
}

func (o *OrderList) GetName() string {
	return o.Name
}

func (o *OrderList) GetPayload() interface{} {
	return o.Payload
}

func (o *OrderList) SetPayload(payload interface{}) {
	o.Payload = payload
}

func (o *OrderList) GetDateTime() time.Time {
	return time.Now()
}
