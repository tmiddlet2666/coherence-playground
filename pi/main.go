package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	weatherCache coherence.NamedCache[string, JsonData]
)

type JsonData map[string]interface{}

func main() {
	session, err := coherence.NewSession(context.Background(), coherence.WithPlainText(), coherence.WithAddress("coherence:///localhost:7574"))
	if err != nil {
		log.Println("unable to connect to Coherence", err)
		return
	}
	defer session.Close()

	weatherCache, err = coherence.GetNamedCache[string, JsonData](session, "weather")
	if err != nil {
		log.Println("unable to create namedCache 'weather'", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./index.html")
	})

	http.HandleFunc("/weather/", weatherHandler)

	fmt.Println("Server running on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

type WeatherData struct {
	Source string   `json:"source"`
	Data   JsonData `json:"data"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	city := r.URL.Query().Get("city")
	if city == "" {
		httpError(w, http.StatusBadRequest, "Missing 'city' query parameter")
		return
	}
	city = strings.ToLower(city)

	ttlParam := r.URL.Query().Get("ttl")
	ttlSec := 30
	if ttlParam != "" {
		var err error
		ttlSec, err = strconv.Atoi(ttlParam)
		if err != nil {
			httpError(w, http.StatusBadRequest, "Invalid 'ttl' parameter; must be an integer")
			return
		}
	}
	ttl := time.Duration(ttlSec) * time.Second
	key := strings.ToLower(city)

	data, err := weatherCache.Get(ctx, key)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if data != nil {
		respond(w, http.StatusOK, WeatherData{
			Source: "Cache",
			Data:   *data,
		})
		return
	}

	// not in cache

	apiData, err := fetchWeather(city)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, WeatherData{
		Source: "API",
		Data:   apiData,
	})

	_, err = weatherCache.PutWithExpiry(ctx, city, apiData, ttl)

	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
	}
}

func httpError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

const WeatherURLTemplate = "https://wttr.in/%s?format=j1"

func fetchWeather(city string) (JsonData, error) {
	url := fmt.Sprintf(WeatherURLTemplate, city)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result JsonData
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return result, nil
}
