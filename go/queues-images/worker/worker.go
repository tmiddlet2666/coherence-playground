package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"github.com/oracle/coherence-go-client/v2/coherence/processors"
	"image"
	"log"
	"net/http"
	"queue-imges/common"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	var (
		ctx             = context.Background()
		cache           coherence.NamedCache[string, common.ImageThumbnail]
		jobsQueue       coherence.NamedQueue[common.ImageJob]
		err             error
		imageJob        *common.ImageJob
		thumnbNailBytes []byte
	)

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

	log.Println("Waiting for jobs...")

	for {
		imageJob, err = jobsQueue.PollHead(ctx)

		if err != nil {
			panic(err)
		}

		if err == nil && imageJob == nil {
			// nothing on the queue, sleep and try again
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		log.Println("Processing", imageJob.ImageURL)
		// we have something, process it
		thumnbNailBytes, err = createThumbnail(imageJob.ImageURL, imageJob.ThumbnailWidth)
		if err != nil {
			log.Printf("error creating thumbnail for %s: %v\n", imageJob.ImageURL, err)
			continue
		}

		// Update the cache with the completed thumbnail
		updater := processors.Update[string]("status", common.StatusCompleted).
			AndThen(processors.Update[[]byte]("thumbnail", thumnbNailBytes))

		_, err = coherence.Invoke[string, common.ImageThumbnail, any](ctx, cache, imageJob.ImageURL, updater)
		if err != nil {
			log.Printf("Unable to update thumbnail, url %s, requeueing", imageJob.ImageURL)
			err = jobsQueue.OfferTail(ctx, *imageJob)
		}

	}
}

func createThumbnail(imageURL string, thumbnailWidth int) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status downloading image: %s", resp.Status)
	}

	srcImage, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	dstImage := imaging.Resize(srcImage, thumbnailWidth, 0, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, dstImage, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
