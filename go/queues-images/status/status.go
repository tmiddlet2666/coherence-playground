package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"github.com/oracle/coherence-go-client/v2/coherence/aggregators"
	"github.com/oracle/coherence-go-client/v2/coherence/extractors"
	"github.com/oracle/coherence-go-client/v2/coherence/filters"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"queue-imges/common"
	"strconv"
)

var (
	ctx       = context.Background()
	config    common.Config
	extractor = extractors.Extract[string]("status")
)

const (
	defaultHttpPort int = 8080
)

//go:embed index.html
var indexPage []byte

func main() {
	var (
		port = defaultHttpPort
		err  error
	)

	cfg, err := common.InitializeCoherence(ctx)
	if err != nil {
		panic(err)
	}

	config = cfg

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort != "" {
		port, err = strconv.Atoi(httpPort)
		if err != nil {
			port = defaultHttpPort
		}
	}

	log.Printf("HTTP server running on %d\n", port)
	log.Println("http://localhost:8080/")
	log.Println("http://localhost:8080/status")
	log.Println("http://localhost:8080/requeue")
	log.Println("http://localhost:8080/image?text=Hello")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/requeue", requeuePending)

	server := &http.Server{
		Addr: "0.0.0.0:" + strconv.Itoa(port),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// getStatus returns the processing status as json
func getStatus(w http.ResponseWriter, r *http.Request) {
	type processingStatus struct {
		TotalImages    int   `json:"totalImages"`
		PendingCount   int64 `json:"pendingCount"`
		ProcessedCount int64 `json:"processedCount"`
		ErrorCount     int64 `json:"errorCount"`
	}

	var (
		stats  = processingStatus{}
		err    error
		result int
		count  *int64
	)

	if r.Method == http.MethodGet {
		// count of images
		result, err = config.Cache.Size(ctx)
		if err == nil {
			stats.TotalImages = result
		}

		// get total pending count
		count, err = coherence.AggregateFilter(ctx, config.Cache, filters.Equal(extractor, common.StatusPending), aggregators.Count())
		if err == nil && count != nil {
			stats.PendingCount = *count
		}

		// get total processed count
		count, err = coherence.AggregateFilter(ctx, config.Cache, filters.Equal(extractor, common.StatusCompleted), aggregators.Count())
		if err == nil && count != nil {
			stats.ProcessedCount = *count
		}

		// get total error count
		count, err = coherence.AggregateFilter(ctx, config.Cache, filters.Equal(extractor, common.StatusError), aggregators.Count())
		if err == nil && count != nil {
			stats.ErrorCount = *count
		}

		_ = json.NewEncoder(w).Encode(stats)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

// requeuePending returns the processing status as json
func requeuePending(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		count := 0
		ch := config.Cache.EntrySetFilter(ctx, filters.Equal(extractor, common.StatusPending))
		for result := range ch {
			if result.Err != nil {
				panic(result.Err)
			}

			// requeue this image as it has been stuck
			job := common.ImageJob{
				ImageURL:       result.Key,
				ThumbnailWidth: 150,
			}

			if err := config.JobsQueue.OfferTail(ctx, job); err != nil {
				panic(err)
			}

			count++
		}

		_ = json.NewEncoder(w).Encode(fmt.Sprintf("re-queued: %d image URLS", count))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		text = "Default"
	}

	width := 1200
	height := 1200

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with a random light color
	bgColor := randomLightColor()
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// Draw text with random darker color
	textColor := randomDarkColor()

	// Randomize text position a little (but stay in safe zone)
	x := width/4 + rand.Intn(30) - 15 // -15 to +15 px
	y := height/2 + rand.Intn(30) - 15

	addLabel(img, x, y, text, textColor)

	// Set headers
	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, img)
}

func addLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{C: col},
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func randomLightColor() color.RGBA {
	r := uint8(rand.Intn(156) + 100) // 100-255
	g := uint8(rand.Intn(156) + 100)
	b := uint8(rand.Intn(156) + 100)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func randomDarkColor() color.RGBA {
	r := uint8(rand.Intn(100))
	g := uint8(rand.Intn(100))
	b := uint8(rand.Intn(100))
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(indexPage)
}
