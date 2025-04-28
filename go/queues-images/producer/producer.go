package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"queue-imges/common"
)

func main() {
	ctx := context.Background()

	config, err := common.InitializeCoherence(ctx)
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
			Status:   common.StatusPending,
			ImageURL: imageURL,
		}

		_, err = config.Cache.Put(ctx, imageURL, imageThumbnail)
		if err != nil {
			log.Printf("Failed to insert record for %s: %v\n", imageURL, err)
			continue
		}

		// Step 2: Send job into queue
		job := common.ImageJob{
			ImageURL:       imageURL,
			ThumbnailWidth: 150,
		}
		err = config.JobsQueue.OfferTail(ctx, job)
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
