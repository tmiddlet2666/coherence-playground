package main

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/coherence"
	"log"
	"math/rand"
	"os"
	"queue-demo/common"
	"strconv"
	"time"
)

func main() {
	var (
		ctx              = context.Background()
		startOrderNumber int
		numOrders        int
		err              error
	)

	// check arguments
	if len(os.Args) != 3 {
		log.Println("provide starting order number and number to complete")
		return
	}

	if startOrderNumber, err = strconv.Atoi(os.Args[1]); err != nil || startOrderNumber < 0 {
		log.Println("invalid value for number of orders")
		return
	}

	if numOrders, err = strconv.Atoi(os.Args[2]); err != nil || numOrders < 0 {
		log.Println("invalid value for starting order")
		return
	}

	// create a new Session to the default gRPC port of 1408 using plain text
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	orderQueue, err := coherence.GetNamedQueue[common.Order](ctx, session, common.QueueName)
	if err != nil {
		panic(err)
	}

	defer orderQueue.Close()

	for i := 0; i < numOrders; i++ {
		newOrder := common.Order{
			OrderNumber: startOrderNumber + i,
			Customer:    fmt.Sprintf("Customer %d", i+1),
			OrderStatus: "NEW",
			OrderTotal:  rand.Float32() * 1000, //nolint
			CreateTime:  time.Now().UnixMilli(),
		}
		err = orderQueue.Offer(newOrder)

		if i%25 == 0 && i != 0 {
			log.Printf("submitted %d orders so far", i)
		}

		if err != nil {
			panic(err)
		}
	}

	log.Printf("Submitted %d orders", numOrders)
}
