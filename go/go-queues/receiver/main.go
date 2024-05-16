package main

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/coherence"
	"log"
	"queue-demo/common"
	"time"
)

func main() {
	var (
		ctx      = context.Background()
		order    *common.Order
		err      error
		received int64
	)

	// create a new Session to the default gRPC port of 1408 using plain text
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	blockingQueue, err := coherence.GetBlockingNamedQueue[common.Order](ctx, session, common.QueueName)
	if err != nil {
		panic(err)
	}

	defer blockingQueue.Close()

	log.Println("Waiting for orders")
	for {
		order, err = blockingQueue.Poll(time.Duration(5) * time.Second)
		if err == coherence.ErrQueueTimedOut {
			continue
		}
		if err != nil {
			panic(err)
		}

		order.CompleteTime = time.Now()
		processingTime := order.CompleteTime.Sub(order.CreateTime)

		received++

		fmt.Printf("Order=%5d created on %v, processing time=%v orders received=%d\n", order.OrderNumber, order.CreateTime, processingTime, received)
	}
}
