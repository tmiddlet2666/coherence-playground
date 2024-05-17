package common

const QueueName = "orders-queue"

// Order represents a fictitious order.
type Order struct {
	OrderNumber  int     `json:"orderNumber"`
	Customer     string  `json:"customer"`
	OrderStatus  string  `json:"orderStatus"`
	OrderTotal   float32 `json:"orderTotal"`
	CreateTime   int64   `json:"createTime"`
	CompleteTime int64   `json:"completeTime"`
}
