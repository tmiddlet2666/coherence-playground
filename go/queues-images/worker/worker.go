package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"github.com/oracle/coherence-go-client/v2/coherence/extractors"
	"github.com/oracle/coherence-go-client/v2/coherence/processors"
	"image"
	"log"
	"net/http"
	"queue-imges/common"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type stats struct {
	jobCount     int64
	processCount int64
	errorCount   int64
}

func (s stats) String() string {
	return fmt.Sprintf("jobCount: %d, processCount: %d, errorCount: %d", s.jobCount, s.processCount, s.errorCount)
}

func main() {
	var (
		ctx            = context.Background()
		err            error
		imageJob       *common.ImageJob
		thumbnailBytes []byte
		imageStatus    string
		workerStats    = stats{}
		displayedStats = false
	)

	config, err := common.InitializeCoherence(ctx)
	if err != nil {
		panic(err)
	}

	// add an index on
	err = coherence.AddIndex(ctx, config.Cache, extractors.Extract[string]("status"), true)
	if err != nil {
		panic(fmt.Sprintf("unable to add index: %v", err))
	}

	for {
		imageJob, err = config.JobsQueue.PollHead(ctx)

		if err != nil {
			panic(err)
		}

		if err == nil && imageJob == nil {
			// nothing on the queue, sleep and try again
			if !displayedStats {
				log.Println("Waiting for jobs...")
				log.Println(workerStats)
				displayedStats = true
			}
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		log.Println("Processing", imageJob.ImageURL)
		workerStats.jobCount++

		// we have something, process it

		thumbnailBytes, err = createThumbnail(imageJob.ImageURL, imageJob.ThumbnailWidth)
		imageStatus = common.StatusCompleted
		if err != nil {
			log.Printf("error creating thumbnail for %s: %v\n", imageJob.ImageURL, err)
			imageStatus = common.StatusError
			thumbnailBytes = make([]byte, 0)
			workerStats.errorCount++
		} else {
			workerStats.processCount++
		}

		displayedStats = false

		// Update the cache with the completed thumbnail
		updater := processors.Update[string]("status", imageStatus).
			AndThen(processors.Update[[]byte]("thumbnail", thumbnailBytes))

		_, err = coherence.Invoke[string, common.ImageThumbnail, any](ctx, config.Cache, imageJob.ImageURL, updater)
		if err != nil {
			log.Printf("Unable to update thumbnail, url %s, requeueing", imageJob.ImageURL)
			err = config.JobsQueue.OfferTail(ctx, *imageJob)
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

	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("not an image: content type = %s", contentType)
	}

	srcImage, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	dstImage := imaging.Resize(srcImage, thumbnailWidth, 0, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, dstImage, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
