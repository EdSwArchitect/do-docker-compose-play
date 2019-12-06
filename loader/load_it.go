package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log2 "github.com/sirupsen/logrus"
)

var es *elasticsearch.Client

func publish(index string, header []string, record []string) {
	log2.Info(fmt.Sprintf("*** Index: %s. Record: %s", index, record))

	// var wg sync.WaitGroup

	// wg.Add(1)

	var body strings.Builder

	body.WriteString(`{`)

	for i, v := range header {
		body.WriteString(`"`)
		body.WriteString(v)
		body.WriteString(`": "`)
		body.WriteString(record[i])
		body.WriteString(`"`)

		if i != len(header)-1 {
			body.WriteString(`,`)
		} else {
			body.WriteString(`}`)
		}
	}

	log2.Info(body.String())

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: uuid.Must(uuid.NewRandom()).String(),
		Body:       strings.NewReader(body.String()),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// wg.Wait()

}

func sendCsv(writer http.ResponseWriter, request *http.Request) {
	index := mux.Vars(request)["index"]

	if len(index) == 0 {
		http.Error(writer, "Index not specified", http.StatusBadRequest)
		return
	}

	fmt.Printf("index %s\n", index)

	log2.Debug(fmt.Sprintf("index is: %s", index))

	var theMap interface{}

	reqBoyd, err := ioutil.ReadAll(request.Body)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)

	} else {

		// Body is expected to be JSON, so just
		// iterate

		json.Unmarshal(reqBoyd, &theMap)

		if theMap == nil {
			http.Error(writer, "Body doesn't contain JSON", http.StatusBadRequest)
			return
		}

		m := theMap.(map[string]interface{})

		// for k, v := range m {
		// 	fmt.Printf("Key: '%s' - Value: '%s'\n", k, v)
		// }

		theFile := m["file"]

		// fmt.Printf("theFile in the map: '%s'\n", theFile)

		log2.Debug(fmt.Sprintf("theFile in the map: '%s'", theFile))

		if theFile != nil {
			// fmt.Printf("The file: %s\n", theFile)

			contents, err := ioutil.ReadFile(fmt.Sprintf("%s", theFile))

			if err == nil {
				// contents is a byte array, convert tostring
				fileContents := string(contents)

				rex := csv.NewReader(strings.NewReader(fileContents))

				var header []string

				for {
					record, err := rex.Read()

					if err == io.EOF {
						break
					}

					if err != nil {
						log.Fatal(err)
					}

					if header != nil {
						publish(index, header, record)
					} else {
						header = record

						// fmt.Printf("Header: %s\n", header)

						log2.Debug(fmt.Sprintf("Header: %s", header))

					}

					// array of strings (the columns of the csv)

					// fmt.Println(record)
				}

				fmt.Fprintln(writer, "Hi")
			} else {

				s := fmt.Sprintf("Unable to open file '%s'. -- %+v", theFile, err)

				http.Error(writer, s, http.StatusNotFound)
			}

		} else {
			http.Error(writer, "File parameter missing", http.StatusNotFound)
		}
	}
}

func main() {

	time.Sleep(30 * time.Second)

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://es01:9200",
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
	router.HandleFunc("/getFile/{index}", sendCsv).Methods("POST")

	log.Fatal(http.ListenAndServe(":18080", router))

}
