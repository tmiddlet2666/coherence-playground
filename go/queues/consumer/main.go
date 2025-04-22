package main

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"log"
	"queue-demo/common"
	"time"
)

func main() {
	var (
		ctx       = context.Background()
		order     *common.Order
		err       error
		received  int64
		totalTime int64
	)

	// create a new Session to the default gRPC port of 1408 using plain text
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	orderQueue, err := coherence.GetNamedQueue[common.Order](ctx, session, common.QueueName, coherence.PagedQueue)
	if err != nil {
		panic(err)
	}

	log.Println("Waiting for orders...")
	for {
		order, err = orderQueue.PollHead(ctx)

		if err != nil {
			panic(err)
		}

		if err == nil && order == nil {
			// nothing on the queue, sleep and try again
			if received > 0 {
				fmt.Printf("No more orders, orders received=%v, totalTime=%d, avg=%v\n", received, totalTime, time.Duration(float64(totalTime)/float64(received))*time.Millisecond)
			}
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		// simulate processing delay
		time.Sleep(time.Duration(10) * time.Millisecond)
		order.CompleteTime = time.Now().UnixMilli()
		processingTime := time.UnixMilli(order.CompleteTime).Sub(time.UnixMilli(order.CreateTime))
		totalTime += processingTime.Milliseconds()

		received++

		fmt.Printf("Order=%d (%s) created on %v, processing time=%v, orders received=%d\n",
			order.OrderNumber, order.Customer, time.UnixMilli(order.CreateTime), processingTime, received)
	}
}
