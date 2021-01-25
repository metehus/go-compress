package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type FileCompressData struct {
	Name    string `json:"file"`
	Output  string `json:"output"`
	Quality int    `json:"quality"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/file", compressFileRoute)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Starting web server on port :6262")
	log.Fatal(http.ListenAndServe(":6262", handler))
}

func indexRoute(w http.ResponseWriter, _ *http.Request) {
	response := map[string]interface{}{
		"working": true,
	}
	json.NewEncoder(w).Encode(response)
	return
}

func compressFileRoute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	var data FileCompressData
	_ = json.NewDecoder(r.Body).Decode(&data)

	folderSegs := strings.Split(data.Output, "/")
	folderPath := strings.Join(folderSegs[:len(folderSegs)-1], "/")

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		response := map[string]interface{}{
			"success": false,
			"message": "Output path does not exist",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	file, _ := os.Open(data.Name)
	defer file.Close()
	img, _, _ := image.Decode(file)

	outputFile, _ := os.Create(data.Output)
	defer outputFile.Close()

	err := jpeg.Encode(outputFile, img, &jpeg.Options{Quality: data.Quality})

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"message": err,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	duration := time.Since(start)

	originalStat, _ := file.Stat()
	compressedStat, _ := outputFile.Stat()

	response := map[string]interface{}{
		"success":       true,
		"duration":      duration.Seconds(),
		"output":        data.Output,
		"original_size": originalStat.Size(),
		"final_size":    compressedStat.Size(),
	}

	json.NewEncoder(w).Encode(response)
}
