package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/s-buhar0v/demoapp/internal/helpers"
	"github.com/s-buhar0v/demoapp/internal/metrics"
	"github.com/s-buhar0v/demoapp/internal/middleware"
	"github.com/s-buhar0v/demoapp/internal/repo"
)

func getHTTPRequestsInflightMax() float64 {
	httpRequestsInflightMaxString := os.Getenv("HTTP_REQUESTS_INFLIGHT_MAX")
	httpRequestsInflightMax := 20.0
	if httpRequestsInflightMaxString != "" {
		httpRequestsInflightMax, _ = strconv.ParseFloat(httpRequestsInflightMaxString, 32)
	}

	return httpRequestsInflightMax
}

var floatingVal = 0

func main() {
	ctx := context.Background()

	db, err := repo.New(ctx)
	if err != nil {
		log.Panic("Database initializing error. ", err)
	} else {
		log.Println("Established connection to repo postgres")
	}
	defer db.Close()

	httpRequestsInflightMax := getHTTPRequestsInflightMax()

	router := chi.NewRouter()
	router.Use(chimiddleware.Logger)
	router.Use(middleware.HTTPMetrics)
	router.Use(middleware.InflightRequests)

	metrics.HttpRequestsInflightMax.WithLabelValues().Set(httpRequestsInflightMax)

	router.Get("/code-2xx", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(helpers.Random2xx())
	})
	router.Get("/code-4xx", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(helpers.Random4xx())
	})
	router.Get("/code-5xx", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(helpers.Random5xx())
	})

	router.Get("/ms-200", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(helpers.RandomDurationMS(200))
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/ms-500", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(helpers.RandomDurationMS(500))
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/ms-1000", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(helpers.RandomDurationMS(1000))
		w.WriteHeader(http.StatusOK)
	})
	router.Get("/ms-6000", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(helpers.RandomDurationMS(6000))
		w.WriteHeader(http.StatusOK)
	})

	strings, err := helpers.ParseFile("text.txt")
	if err != nil {
		log.Printf("Failed to generate data from text.txt, error: %s", err)
	}

	router.Get("/trigger/save", func(w http.ResponseWriter, r *http.Request) {
		err = db.SaveData(ctx, strings)
		if err != nil {
			log.Println(err)
		}

		w.WriteHeader(http.StatusOK)
	})

	router.Get("/trigger/get", func(w http.ResponseWriter, r *http.Request) {
		allData, err := db.GetData(ctx, false)
		if err != nil {
			log.Println(err)
		}

		writeResponse(w, http.StatusOK, allData)
	})

	router.Get("/trigger/get/long", func(w http.ResponseWriter, r *http.Request) {
		allData, err := db.GetData(ctx, true)
		if err != nil {
			log.Println(err)
		}

		writeResponse(w, http.StatusOK, allData)
	})

	router.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func writeResponse(writer http.ResponseWriter, code int, v interface{}) {
	body, _ := json.Marshal(v)
	writer.WriteHeader(code)
	_, err := writer.Write(body)
	if err != nil {
		log.Print(err) // logging with log package because fmt package is not concurrent safe
	}
}

// randomNumber возвращает случайное число от 1 до 20
func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(20) + 1
}
