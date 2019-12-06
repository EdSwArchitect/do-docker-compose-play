package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	log2 "github.com/sirupsen/logrus"
)

var es *elasticsearch.Client

func query(writer http.ResponseWriter, request *http.Request) {
	index := mux.Vars(request)["index"]

	if len(index) == 0 {
		http.Error(writer, "Index not specified", http.StatusBadRequest)
		return
	}

	fmt.Printf("index %s\n", index)

	log2.Debug(fmt.Sprintf("index is: %s", index))

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	} else {
		fmt.Println("No error")
	}
}

func main() {
	fmt.Println("Sleeping 30 seconds before connecting to ESP")

	hostName := flag.String("host", "http://localhost:9200", "The ESP connection URI <http://localhost:9200>")

	flag.Parse()

	time.Sleep(30 * time.Second)

	fmt.Printf("Connecting to ESP '%s'\n", *hostName)

	cfg := elasticsearch.Config{
		Addresses: []string{
			*hostName,
		},
	}
	// es, _ = elasticsearch.NewDefaultClient()

	es, _ = elasticsearch.NewClient(cfg)

	log2.SetLevel(log2.DebugLevel)

	// log.Println(elasticsearch.Version)
	// log.Println(es.Info())

	log2.Debug(elasticsearch.Version)
	log2.Debug(es.Info())

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/query/{index}", query).Methods("GET")

	log.Fatal(http.ListenAndServe(":19090", router))
}
