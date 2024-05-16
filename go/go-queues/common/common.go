package common

import "time"

const QueueName = "orders-queue"

// Order represents a fictitious order.
type Order struct {
	OrderNumber  int       `json:"orderID"`
	Customer     string    `json:"customer"`
	OrderStatus  string    `json:"orderStatus"`
	OrderTotal   float32   `json:"orderTotal"`
	CreateTime   time.Time `json:"createTime"`
	CompleteTime time.Time `json:"completeTime"`
}
