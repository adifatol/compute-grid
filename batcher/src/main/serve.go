package main

import (
	"batch"
	"encoding/json"
	"fmt"
	"net/http"
	"response"
	"flag"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var batchHandler *batch.ReadHandler
var responseHandler *response.DbHandler
var progressHandler *response.ProgressHandler

var data_file = flag.String("data_file", "GSE13507_illumina_raw.txt", "Input the name of the file to be loaded.")
var Log 	  = logging.MustGetLogger("general")

func main() {

	flag.Parse()

	batchHandler = batch.NewHandler("../data/" + *data_file)
	defer batch.Terminate(batchHandler)

	responseHandler = response.ConnectDb("result.db")
	defer response.Disconnect(responseHandler)

	progressHandler = response.NewProgress(batchHandler)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/nextBatch", nextBatch)
	router.HandleFunc("/customBatch/{loadedId}", customBatch)
	router.HandleFunc("/resultBatch", resultBatch)
	router.HandleFunc("/reset", Reset)

	Log.Critical(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Smart Computing Grid Batcher")
}

func nextBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Batch-UUID")
	w.Header().Set("Batch-UUID", "ASASASHUAHS")

	data := batch.ReadNext(batchHandler)

	Log.Info("New batch sent")
	fmt.Fprintln(w, data)
}

func customBatch(w http.ResponseWriter, r *http.Request) {
	Log.Critical("Not ready yet")
}

func resultBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Access-Control-Expose-Headers", "Batch-UUID")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Batch-UUID")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")

	var results []interface{}
	json.NewDecoder(r.Body).Decode(&results)

	response.Save(results, responseHandler)
	response.Progress(progressHandler)

	Log.Info("data received")
}

func Reset(w http.ResponseWriter, r *http.Request) {
	/* Access-Control-Allow-Origin not defined, local access only */
	batch.Reset(batchHandler)
}
