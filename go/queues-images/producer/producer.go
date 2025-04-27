package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"log"
	"os"
	"queue-imges/common"
)

var (
	ctx       = context.Background()
	cache     coherence.NamedCache[string, common.ImageThumbnail]
	jobsQueue coherence.NamedQueue[common.ImageJob]
)

func main() {
	// create a new Session to the default gRPC port of 1408 using plain text
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cache, err = coherence.GetNamedCache[string, common.ImageThumbnail](session, common.CacheName)
	if err != nil {
		panic(err)
	}

	jobsQueue, err = coherence.GetNamedQueue[common.ImageJob](ctx, session, common.JobQueueName, coherence.PagedQueue)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: producer <file-with-urls.txt>")
		os.Exit(1)
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		imageURL := scanner.Text()
		if imageURL == "" {
			continue
		}

		imageThumbnail := common.ImageThumbnail{
			Status:   common.StatusProcessing,
			ImageURL: imageURL,
		}

		log.Println("Image URL", imageURL)

		_, err = cache.Put(ctx, imageURL, imageThumbnail)
		if err != nil {
			log.Printf("Failed to insert record for %s: %v\n", imageURL, err)
			continue
		}

		// Step 2: Send job into queue
		job := common.ImageJob{
			ImageURL:       imageURL,
			ThumbnailWidth: 150,
		}
		err = jobsQueue.OfferTail(ctx, job)
		if err != nil {
			log.Printf("Failed to offer job for %s: %v\n", imageURL, err)
			continue
		}

		fmt.Printf("Submitted job for: %s\n", imageURL)
		count++
	}

	if err = scanner.Err(); err != nil {
		panic(fmt.Errorf("error reading file: %w", err))
	}

	fmt.Printf("Submitted %d jobs.\n", count)
}
