package main

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	start := time.Now()

	log.Println("migration process starting")
	defer func() {
		log.Printf("migration process finished in %s\n", time.Since(start))
	}()

	ctx := context.Background()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27019")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("moviedb").Collection("movies")

	// create unique index on the "id" field
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, idIndex)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("decompressing the csv file")

	zipFilePath := "assets/TMDB_movie_dataset_v11.csv.zip"
	outputDir := "assets"
	csvFilePath, err := unzip(zipFilePath, outputDir)
	if err != nil {
		log.Fatal(err)
	}
	// delete the decompressed csv file after processing
	defer os.Remove(csvFilePath)

	log.Println("decompressing finished")

	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read() // read the first line (headers)
	if err != nil {
		log.Fatal(err)
	}

	// map the headers to indices for easier access
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// channels for concurrency
	movieChan := make(chan domain.Movie, 100)
	var wg sync.WaitGroup

	// mutex for synchronizing access to idMap
	var mutex sync.Mutex
	// map to store ids to prevent duplicate inserts
	idMap := make(map[string]struct{})

	// worker goroutines to process the csv rows
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bulkMovies := make([]mongo.WriteModel, 0, 100) // batch size: 100
			for movie := range movieChan {
				// check if id already exists in the map
				mutex.Lock()
				if _, exists := idMap[movie.ID]; !exists {
					idMap[movie.ID] = struct{}{} // add id to map to prevent duplicates
					mutex.Unlock()

					// insert the movie
					bulkMovies = append(bulkMovies, mongo.NewInsertOneModel().SetDocument(movie))
					if len(bulkMovies) >= 100 {
						_, err := collection.BulkWrite(ctx, bulkMovies)
						if err != nil {
							log.Println("error during bulk write:", err)
						}
						bulkMovies = bulkMovies[:0]
					}
				} else {
					mutex.Unlock()
				}
			}
			if len(bulkMovies) > 0 {
				_, err := collection.BulkWrite(ctx, bulkMovies)
				if err != nil {
					log.Println("error during final bulk write:", err)
				}
			}
		}()
	}

	log.Println("started reading the csv rows and inserting data into the database")
	// read and process the csv rows
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Fatal(err)
			}

			// check if contains non latin alphabet characters e.g. cyrillic or kanji
			// accents are valid
			if hasNonLatinCharacters(record[headerMap["title"]]) {
				continue //skip
			}

			genres := strings.Split(record[headerMap["genres"]], ", ")
			keywords := strings.Split(record[headerMap["keywords"]], ", ")

			movie := domain.Movie{
				ID:            record[headerMap["id"]],
				Title:         record[headerMap["title"]],
				OriginalTitle: record[headerMap["original_title"]],
				Overview:      record[headerMap["overview"]],
				Poster:        "https://image.tmdb.org/t/p/w220_and_h330_face" + record[headerMap["poster_path"]],
				Genres:        genres,
				Keywords:      keywords,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			// send the movie to the channel
			movieChan <- movie
		}
		close(movieChan)
	}()

	wg.Wait()
}

func unzip(src string, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			outPath := filepath.Join(dest, f.Name)
			outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			if err != nil {
				return "", err
			}
			return outPath, nil
		}
	}
	return "", errors.New("no file found in zip archive")
}

func hasNonLatinCharacters(s string) bool {
	for _, char := range s {
		if !isLatinOrAccent(char) {
			return true
		}
	}
	return false
}

func isLatinOrAccent(char rune) bool {
	// basic Latin block: U+0000 to U+007F
	if char <= 0x007F {
		return true
	}
	// Latin-1 Supplement block: U+0080 to U+00FF
	if char >= 0x00A0 && char <= 0x00FF {
		return true
	}
	// Latin Extended-A block: U+0100 to U+017F
	if char >= 0x0100 && char <= 0x017F {
		return true
	}
	// Latin Extended-B block: U+0180 to U+024F
	if char >= 0x0180 && char <= 0x024F {
		return true
	}

	return false
}
