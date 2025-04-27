package common

const (
	JobQueueName = "image-jobs"
	CacheName    = "thumbnails"

	StatusProcessing = "processing"
	StatusCompleted  = "completed"
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
