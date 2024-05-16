package main

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/coherence"
	"github.com/oracle/coherence-go-client/coherence/aggregators"
	"github.com/oracle/coherence-go-client/coherence/extractors"
	"github.com/oracle/coherence-go-client/coherence/filters"
	"github.com/oracle/coherence-go-client/coherence/processors"
	"log"
	"math/rand"
	"time"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p Person) String() string {
	return fmt.Sprintf("Person{id=%d, name=%s, age=%d}", p.ID, p.Name, p.Age)
}

const (
	peopleCount = 100_000
)

func main() {
	ctx := context.Background()
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		log.Fatal("unable to connect to coherence", err)
	}

	defer session.Close()

	nm, err := coherence.GetNamedMap[int, Person](session, "people")
	if err != nil {
		log.Fatal("unable to create new map 'people", err)
	}

	// add an index on age
	age := extractors.Extract[int]("age")
	err = coherence.AddIndex(ctx, nm, age, true)
	if err != nil {
		log.Fatal(err)
	}

	err = nm.Clear(ctx)
	if err != nil {
		log.Fatal("unable to clear map", err)
	}

	log.Printf("Adding %d people\n", peopleCount)
	batchSize := 10_000
	buffer := make(map[int]Person)
	for i := 1; i <= peopleCount; i++ {
		p := Person{ID: i, Name: fmt.Sprintf("Name-%d", i), Age: 15 + i%50}
		buffer[p.ID] = p
		if i%batchSize == 0 {
			log.Println("Writing batch of", batchSize)
			if err = nm.PutAll(ctx, buffer); err != nil {
				panic(err)
			}
			buffer = make(map[int]Person)
		}
	}

	// write any left in the buffer
	if len(buffer) > 0 {
		if err = nm.PutAll(ctx, buffer); err != nil {
			panic(err)
		}
	}

	// do some random gets
	iter := 0
	for {
		iter++
		log.Println("")
		log.Println("Iter=", iter)

		log.Println("200 random gets")
		for i := 1; i < 200; i++ {
			_, err1 := nm.Get(ctx, rand.Intn(peopleCount)+1)
			if err1 != nil {
				log.Fatal(err1)
			}
		}

		log.Println("100 random updates via key using InvokeAllKeys")
		keys := make([]int, 100)
		for i := 0; i < 100; i++ {
			keys[i] = rand.Intn(peopleCount) + 1
		}

		_ = coherence.InvokeAllKeys[int, Person, int](ctx, nm, keys, processors.Update[int]("age", 15+rand.Intn(20)))

		count, err := coherence.AggregateFilter(ctx, nm, filters.Between(age, 20, 40), aggregators.Count())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Aggregating via age =", *count)

		log.Println("Sleeping...")
		time.Sleep(time.Duration(5) * time.Second)
	}
}
