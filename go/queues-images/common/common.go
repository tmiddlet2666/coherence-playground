package common

import (
	"context"
	"github.com/oracle/coherence-go-client/v2/coherence"
)

const (
	JobQueueName = "image-jobs"
	CacheName    = "thumbnails"

	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusError     = "error"
)

// ImageJob represents a job to create an image thumbnail.
type ImageJob struct {
	ImageURL       string `json:"imageURL"`
	ThumbnailWidth int    `json:"thumbnailWidth"`
}

type ImageThumbnail struct {
	ImageURL  string `json:"imageURL"`
	Status    string `json:"status"`
	Thumbnail []byte `json:"thumbnail"`
}

type Config struct {
	Session   *coherence.Session
	Cache     coherence.NamedCache[string, ImageThumbnail]
	JobsQueue coherence.NamedQueue[ImageJob]
}

func InitializeCoherence(ctx context.Context) (Config, error) {
	var (
		err    error
		config = Config{}
	)

	// create a new Session to the default gRPC port of 1408 using plain text
	config.Session, err = coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		return config, err
	}

	config.Cache, err = coherence.GetNamedCache[string, ImageThumbnail](config.Session, CacheName)
	if err != nil {
		return config, err
	}

	config.JobsQueue, err = coherence.GetNamedQueue[ImageJob](ctx, config.Session, JobQueueName, coherence.PagedQueue)
	if err != nil {
		return config, err
	}

	return config, nil
}
